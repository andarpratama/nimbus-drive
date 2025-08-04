package handlers
import (
	"net/http"
	"time"
	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
func CreateFolder(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var input struct {
		Name     string      `json:"name" binding:"required"`
		ParentID *uuid.UUID `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.ParentID != nil {
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.ParentID, userID).First(&parentFolder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent folder not found"})
			return
		}
	}
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ?", userID, input.Name)
	if input.ParentID != nil {
		query = query.Where("parent_id = ?", *input.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.First(&existingFolder).Error; err == nil {
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
func GetFolders(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
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
	type FolderWithCount struct {
		models.Folder
		FileCount    int64 `json:"file_count"`
		SubfolderCount int64 `json:"subfolder_count"`
		TotalItems   int64 `json:"total_items"`
	}
	var foldersWithCounts []FolderWithCount
	for _, folder := range folders {
		var fileCount int64
		database.DB.Model(&models.File{}).Where("folder_id = ? AND user_id = ?", folder.ID, userID).Count(&fileCount)
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
func GetFolderByID(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
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
func UpdateFolder(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Param("id")
	var input struct {
		Name     string      `json:"name" binding:"required"`
		ParentID *uuid.UUID `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	if input.ParentID != nil {
		if input.ParentID.String() == folder.ID.String() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "folder cannot be its own parent"})
			return
		}
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.ParentID, userID).First(&parentFolder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent folder not found"})
			return
		}
	}
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ? AND id != ?", userID, input.Name, folderID)
	if input.ParentID != nil {
		query = query.Where("parent_id = ?", *input.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.First(&existingFolder).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in this location"})
		return
	}
	folder.Name = input.Name
	folder.ParentID = input.ParentID
	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update folder"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "folder updated", "folder": folder})
}
func DeleteFolder(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Param("id")
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	var subfolderCount int64
	if err := database.DB.Model(&models.Folder{}).Where("parent_id = ?", folderID).Count(&subfolderCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subfolders"})
		return
	}
	if subfolderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete folder with subfolders"})
		return
	}
	var fileCount int64
	if err := database.DB.Model(&models.File{}).Where("folder_id = ?", folderID).Count(&fileCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check files"})
		return
	}
	if fileCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete folder with files"})
		return
	}
	if err := database.DB.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete folder"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}
func GetFolderTree(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
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
	folderMap := make(map[uuid.UUID]*models.Folder)
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
func GetFolderContents(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Param("id")
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	var subfolders []models.Folder
	if err := database.DB.Where("parent_id = ? AND user_id = ?", folderID, userID).Find(&subfolders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subfolders"})
		return
	}
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
func RenameFolder(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Param("id")
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
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
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ? AND id != ?", userID, cleanName, folderID)
	if folder.ParentID != nil {
		query = query.Where("parent_id = ?", folder.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.First(&existingFolder).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in this location"})
		return
	}
	folder.Name = cleanName
	folder.UpdatedAt = time.Now()
	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename folder"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "folder renamed successfully", "folder": folder})
} 
func MoveFolder(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Param("id")
	var input struct {
		NewParentID *uuid.UUID `json:"new_parent_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	if input.NewParentID != nil && input.NewParentID.String() == folder.ID.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder cannot be moved to itself"})
		return
	}
	if input.NewParentID != nil {
		var newParentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.NewParentID, userID).First(&newParentFolder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "destination folder not found"})
			return
		}
		if isDescendant(folder.ID, *input.NewParentID, userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot move folder into its own subfolder"})
			return
		}
	}
	var existingFolder models.Folder
	query := database.DB.Where("user_id = ? AND name = ? AND id != ?", userID, folder.Name, folderID)
	if input.NewParentID != nil {
		query = query.Where("parent_id = ?", *input.NewParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.First(&existingFolder).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "a folder with this name already exists in the destination"})
		return
	}
	folder.ParentID = input.NewParentID
	folder.UpdatedAt = time.Now()
	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move folder"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "folder moved successfully", "folder": folder})
}
func isDescendant(ancestorID, targetID uuid.UUID, userID uuid.UUID) bool {
	var folders []models.Folder
	if err := database.DB.Where("user_id = ?", userID).Find(&folders).Error; err != nil {
		return false
	}
	childrenMap := make(map[uuid.UUID][]uuid.UUID)
	for _, folder := range folders {
		if folder.ParentID != nil {
			childrenMap[*folder.ParentID] = append(childrenMap[*folder.ParentID], folder.ID)
		}
	}
	visited := make(map[uuid.UUID]bool)
	queue := []uuid.UUID{ancestorID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		if current == targetID {
			return true
		}
		for _, child := range childrenMap[current] {
			if !visited[child] {
				queue = append(queue, child)
			}
		}
	}
	return false
} 