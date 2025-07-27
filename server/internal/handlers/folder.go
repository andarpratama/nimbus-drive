package handlers

import (
	"net/http"
	"time"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateFolder handles POST /folders to create a new folder
func CreateFolder(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If parent_id is provided, verify it exists and belongs to the user
	if input.ParentID != nil {
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.ParentID, userID).First(&parentFolder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent folder not found"})
			return
		}
	}

	// Check for duplicate folder name in the same parent folder
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ?", userID, input.Name)
	
	if input.ParentID != nil {
		// Check in specific parent folder
		query = query.Where("parent_id = ?", *input.ParentID)
	} else {
		// Check in root folder (parent_id is NULL)
		query = query.Where("parent_id IS NULL")
	}
	
	if err := query.First(&existingFolder).Error; err == nil {
		// Folder with same name already exists in this location
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in this location"})
		return
	}

	folder := models.Folder{
		Name:     input.Name,
		ParentID: input.ParentID,
		UserID:   userID,
	}

	if err := database.DB.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "folder created", "folder": folder})
}

// GetFolders handles GET /folders to fetch all folders for a user
func GetFolders(c *gin.Context) {
	userID := c.GetUint("userID")

	var folders []models.Folder
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("Parent").
		Preload("Files").
		Order("created_at DESC").
		Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch folders"})
		return
	}

	// Add item counts to each folder
	type FolderWithCount struct {
		models.Folder
		FileCount    int64 `json:"file_count"`
		SubfolderCount int64 `json:"subfolder_count"`
		TotalItems   int64 `json:"total_items"`
	}

	var foldersWithCounts []FolderWithCount
	for _, folder := range folders {
		// Count files in this folder
		var fileCount int64
		database.DB.Model(&models.File{}).Where("folder_id = ? AND user_id = ?", folder.ID, userID).Count(&fileCount)
		
		// Count subfolders in this folder
		var subfolderCount int64
		database.DB.Model(&models.Folder{}).Where("parent_id = ? AND user_id = ?", folder.ID, userID).Count(&subfolderCount)
		
		folderWithCount := FolderWithCount{
			Folder:       folder,
			FileCount:    fileCount,
			SubfolderCount: subfolderCount,
			TotalItems:   fileCount + subfolderCount,
		}
		foldersWithCounts = append(foldersWithCounts, folderWithCount)
	}

	c.JSON(http.StatusOK, gin.H{"folders": foldersWithCounts})
}

// GetFolderByID handles GET /folders/:id to fetch a single folder
func GetFolderByID(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Param("id")

	var folder models.Folder
	if err := database.DB.
		Where("id = ? AND user_id = ?", folderID, userID).
		Preload("Parent").
		Preload("Files").
		First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

// UpdateFolder handles PUT /folders/:id to update a folder
func UpdateFolder(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Param("id")

	var input struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if folder exists and belongs to user
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	// If parent_id is provided, verify it exists and belongs to the user
	if input.ParentID != nil {
		// Prevent circular reference - folder cannot be its own parent
		if *input.ParentID == folder.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "folder cannot be its own parent"})
			return
		}

		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.ParentID, userID).First(&parentFolder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent folder not found"})
			return
		}
	}

	// Check for duplicate folder name in the same parent folder (excluding current folder)
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ? AND id != ?", userID, input.Name, folderID)
	
	if input.ParentID != nil {
		// Check in specific parent folder
		query = query.Where("parent_id = ?", *input.ParentID)
	} else {
		// Check in root folder (parent_id is NULL)
		query = query.Where("parent_id IS NULL")
	}
	
	if err := query.First(&existingFolder).Error; err == nil {
		// Folder with same name already exists in this location
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in this location"})
		return
	}

	// Update folder
	folder.Name = input.Name
	folder.ParentID = input.ParentID

	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder updated", "folder": folder})
}

// DeleteFolder handles DELETE /folders/:id to delete a folder
func DeleteFolder(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Param("id")

	// Check if folder exists and belongs to user
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	// Check if folder has subfolders
	var subfolderCount int64
	if err := database.DB.Model(&models.Folder{}).Where("parent_id = ?", folderID).Count(&subfolderCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subfolders"})
		return
	}

	if subfolderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete folder with subfolders"})
		return
	}

	// Check if folder has files
	var fileCount int64
	if err := database.DB.Model(&models.File{}).Where("folder_id = ?", folderID).Count(&fileCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check files"})
		return
	}

	if fileCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete folder with files"})
		return
	}

	// Soft delete the folder
	if err := database.DB.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}

// GetFolderTree handles GET /folders/tree to get folder hierarchy
func GetFolderTree(c *gin.Context) {
	userID := c.GetUint("userID")

	var folders []models.Folder
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("Parent").
		Preload("Files").
		Order("name ASC").
		Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch folder tree"})
		return
	}

	// Build folder tree structure
	folderMap := make(map[uint]*models.Folder)
	var rootFolders []*models.Folder

	for i := range folders {
		folderMap[folders[i].ID] = &folders[i]
	}

	for i := range folders {
		if folders[i].ParentID == nil {
			rootFolders = append(rootFolders, &folders[i])
		}
	}

	c.JSON(http.StatusOK, gin.H{"folders": rootFolders})
}

// GetFolderContents handles GET /folders/:id/contents to get folder contents
func GetFolderContents(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Param("id")

	// Verify folder exists and belongs to user
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	// Get subfolders
	var subfolders []models.Folder
	if err := database.DB.Where("parent_id = ? AND user_id = ?", folderID, userID).Find(&subfolders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subfolders"})
		return
	}

	// Get files in this folder
	var files []models.File
	if err := database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder":     folder,
		"subfolders": subfolders,
		"files":      files,
	})
}

// RenameFolder handles PATCH /folders/:id/rename to rename a folder
func RenameFolder(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// Clean the folder name (remove whitespace)
	cleanName := ""
	for _, r := range input.Name {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			cleanName += string(r)
		}
	}

	if cleanName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
		return
	}

	// Check if folder exists and belongs to user
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	// Check for duplicate folder name in the same parent folder
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ? AND id != ?", userID, cleanName, folderID)
	
	if folder.ParentID != nil {
		// Check in specific parent folder
		query = query.Where("parent_id = ?", folder.ParentID)
	} else {
		// Check in root folder (parent_id is NULL)
		query = query.Where("parent_id IS NULL")
	}
	
	if err := query.First(&existingFolder).Error; err == nil {
		// Folder with same name already exists in this location
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in this location"})
		return
	}

	// Update the folder name and UpdatedAt timestamp
	folder.Name = cleanName
	folder.UpdatedAt = time.Now()

	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder renamed successfully", "folder": folder})
} 