package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StarFile stars a file for the authenticated user
func StarFile(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("StarFile: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("StarFile: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		log.Printf("StarFile: Failed to parse fileID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	// Check if file exists and belongs to user
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		log.Printf("StarFile: File not found or doesn't belong to user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Check if already starred
	var existingStarred models.Starred
	if err := database.DB.Where("user_id = ? AND file_id = ?", userID, fileID).First(&existingStarred).Error; err == nil {
		log.Printf("StarFile: File already starred")
		c.JSON(http.StatusConflict, gin.H{"error": "file already starred"})
		return
	}

	// Create starred record
	starred := models.Starred{
		UserID: userID,
		FileID: &fileID,
	}

	if err := database.DB.Create(&starred).Error; err != nil {
		log.Printf("StarFile: Failed to create starred record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to star file"})
		return
	}

	log.Printf("StarFile: Successfully starred file %s", fileID)
	c.JSON(http.StatusCreated, gin.H{"message": "file starred successfully"})
}

// UnstarFile removes star from a file
func UnstarFile(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UnstarFile: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("UnstarFile: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		log.Printf("UnstarFile: Failed to parse fileID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	// Delete starred record
	result := database.DB.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&models.Starred{})
	if result.Error != nil {
		log.Printf("UnstarFile: Failed to delete starred record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unstar file"})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("UnstarFile: File not starred")
		c.JSON(http.StatusNotFound, gin.H{"error": "file not starred"})
		return
	}

	log.Printf("UnstarFile: Successfully unstarred file %s", fileID)
	c.JSON(http.StatusOK, gin.H{"message": "file unstarred successfully"})
}

// StarFolder stars a folder for the authenticated user
func StarFolder(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("StarFolder: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("StarFolder: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	folderIDStr := c.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		log.Printf("StarFolder: Failed to parse folderID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	// Check if folder exists and belongs to user
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		log.Printf("StarFolder: Folder not found or doesn't belong to user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	// Check if already starred
	var existingStarred models.Starred
	if err := database.DB.Where("user_id = ? AND folder_id = ?", userID, folderID).First(&existingStarred).Error; err == nil {
		log.Printf("StarFolder: Folder already starred")
		c.JSON(http.StatusConflict, gin.H{"error": "folder already starred"})
		return
	}

	// Create starred record
	starred := models.Starred{
		UserID:   userID,
		FolderID: &folderID,
	}

	if err := database.DB.Create(&starred).Error; err != nil {
		log.Printf("StarFolder: Failed to create starred record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to star folder"})
		return
	}

	log.Printf("StarFolder: Successfully starred folder %s", folderID)
	c.JSON(http.StatusCreated, gin.H{"message": "folder starred successfully"})
}

// UnstarFolder removes star from a folder
func UnstarFolder(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UnstarFolder: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("UnstarFolder: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	folderIDStr := c.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		log.Printf("UnstarFolder: Failed to parse folderID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	// Delete starred record
	result := database.DB.Where("user_id = ? AND folder_id = ?", userID, folderID).Delete(&models.Starred{})
	if result.Error != nil {
		log.Printf("UnstarFolder: Failed to delete starred record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unstar folder"})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("UnstarFolder: Folder not starred")
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not starred"})
		return
	}

	log.Printf("UnstarFolder: Successfully unstarred folder %s", folderID)
	c.JSON(http.StatusOK, gin.H{"message": "folder unstarred successfully"})
}

// ListStarredItems returns all starred items (files and folders) for the authenticated user
func ListStarredItems(c *gin.Context) {
	log.Printf("ListStarredItems: Starting function")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ListStarredItems: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("ListStarredItems: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var starredItems []models.Starred
	if err := database.DB.Preload("File").Preload("Folder").Where("user_id = ?", userID).Find(&starredItems).Error; err != nil {
		log.Printf("ListStarredItems: Failed to fetch starred items: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch starred items"})
		return
	}
	
	log.Printf("ListStarredItems: Raw starred items from DB: %+v", starredItems)
	
	// Test the first item
	if len(starredItems) > 0 {
		firstItem := starredItems[0]
		log.Printf("ListStarredItems: First item - FileID: %v, FolderID: %v, File: %v, Folder: %v", 
			firstItem.FileID, firstItem.FolderID, firstItem.File != nil, firstItem.Folder != nil)
		if firstItem.Folder != nil {
			log.Printf("ListStarredItems: Folder name: %s", firstItem.Folder.Name)
		}
	}
	
	log.Printf("ListStarredItems: Fetched %d starred items", len(starredItems))

	// Format response
	var response []gin.H
	for i, item := range starredItems {
		log.Printf("ListStarredItems: Processing item %d - FileID: %v, FolderID: %v, File: %v, Folder: %v", i, item.FileID, item.FolderID, item.File != nil, item.Folder != nil)
		log.Printf("ListStarredItems: Item details - ID: %s, Type: %s", item.ID, item.Type)
		if item.FileID != nil && item.File != nil {
			// Format file size for display
			sizeStr := formatFileSize(item.File.Size)
			
			response = append(response, gin.H{
				"id":        item.ID,
				"type":      "file",
				"file_id":   item.FileID,
				"name":      item.File.Name,
				"size":      sizeStr,
				"mime_type": item.File.MimeType,
				"created_at": item.CreatedAt,
			})
		} else if item.FolderID != nil && item.Folder != nil {
			log.Printf("ListStarredItems: Entering folder processing block")
			// Get folder item count for display
			var fileCount int64
			var subfolderCount int64
			
			// Debug logging
			log.Printf("ListStarredItems: Processing folder %s (%s)", item.Folder.Name, item.FolderID)
			
			// Count files in folder
			if err := database.DB.Model(&models.File{}).Where("folder_id = ? AND deleted_at IS NULL", item.FolderID).Count(&fileCount).Error; err != nil {
				log.Printf("ListStarredItems: Error counting files in folder %s: %v", item.FolderID, err)
			}
			
			// Count subfolders in folder
			if err := database.DB.Model(&models.Folder{}).Where("parent_id = ? AND deleted_at IS NULL", item.FolderID).Count(&subfolderCount).Error; err != nil {
				log.Printf("ListStarredItems: Error counting subfolders in folder %s: %v", item.FolderID, err)
			}
			
			totalItems := fileCount + subfolderCount
			log.Printf("ListStarredItems: Folder %s has %d files and %d subfolders (total: %d)", item.Folder.Name, fileCount, subfolderCount, totalItems)
			
			var size string
			if totalItems > 0 {
				size = fmt.Sprintf("%d item%s", totalItems, totalItems != 1 ? "s" : "")
			} else {
				size = "Empty"
			}
			
			log.Printf("ListStarredItems: Setting folder size to: %s", size)
			
			response = append(response, gin.H{
				"id":           item.ID,
				"type":         "folder",
				"folder_id":    item.FolderID,
				"name":         item.Folder.Name,
				"size":         totalItems,
				"total_items":  totalItems,
				"created_at":   item.CreatedAt,
			})
		} else {
			log.Printf("ListStarredItems: Item doesn't match file or folder conditions - FileID: %v, FolderID: %v", item.FileID, item.FolderID)
		}
	}

	log.Printf("ListStarredItems: Found %d starred items", len(response))
	c.JSON(http.StatusOK, gin.H{
		"starred_items": response,
		"total_items":   len(response),
	})
}

// CheckStarredStatus checks if a file or folder is starred by the user
func CheckStarredStatus(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("CheckStarredStatus: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("CheckStarredStatus: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		log.Printf("CheckStarredStatus: Failed to parse itemID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	itemType := c.Query("type") // "file" or "folder"
	if itemType != "file" && itemType != "folder" {
		log.Printf("CheckStarredStatus: Invalid item type: %s", itemType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item type"})
		return
	}

	var starred models.Starred
	var query string
	if itemType == "file" {
		query = "user_id = ? AND file_id = ?"
	} else {
		query = "user_id = ? AND folder_id = ?"
	}

	isStarred := database.DB.Where(query, userID, itemID).First(&starred).Error == nil

	log.Printf("CheckStarredStatus: Item %s (type: %s) starred: %t", itemID, itemType, isStarred)
	c.JSON(http.StatusOK, gin.H{"is_starred": isStarred})
}

// formatFileSize formats bytes into human readable format
func formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	
	if exp >= len(sizes) {
		exp = len(sizes) - 1
	}
	
	value := float64(bytes) / math.Pow(unit, float64(exp))
	return fmt.Sprintf("%.1f %s", value, sizes[exp])
} 