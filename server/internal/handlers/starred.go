package handlers
import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
func StarItem(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("StarItem: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("StarItem: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var input struct {
		ItemID   string `json:"item_id" binding:"required"`
		ItemType string `json:"item_type" binding:"required,oneof=file folder"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("StarItem: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		log.Printf("StarItem: Failed to parse itemID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	if input.ItemType == "file" {
		var file models.File
		if err := database.DB.Where("id = ? AND user_id = ?", itemID, userID).First(&file).Error; err != nil {
			log.Printf("StarItem: File not found or doesn't belong to user: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
	} else if input.ItemType == "folder" {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", itemID, userID).First(&folder).Error; err != nil {
			log.Printf("StarItem: Folder not found or doesn't belong to user: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
	}
	var existingStarred models.Starred
	var query string
	if input.ItemType == "file" {
		query = "user_id = ? AND file_id = ?"
	} else {
		query = "user_id = ? AND folder_id = ?"
	}
	if err := database.DB.Where(query, userID, itemID).First(&existingStarred).Error; err == nil {
		log.Printf("StarItem: Item already starred")
		c.JSON(http.StatusConflict, gin.H{"error": "item already starred"})
		return
	}
	starred := models.Starred{
		UserID: userID,
	}
	if input.ItemType == "file" {
		starred.FileID = &itemID
	} else {
		starred.FolderID = &itemID
	}
	if err := database.DB.Create(&starred).Error; err != nil {
		log.Printf("StarItem: Failed to create starred record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to star item"})
		return
	}
	log.Printf("StarItem: Successfully starred %s %s", input.ItemType, itemID)
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%s starred successfully", input.ItemType),
		"starred": true,
	})
}
func UnstarItem(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UnstarItem: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("UnstarItem: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var input struct {
		ItemID   string `json:"item_id" binding:"required"`
		ItemType string `json:"item_type" binding:"required,oneof=file folder"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("UnstarItem: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		log.Printf("UnstarItem: Failed to parse itemID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	var query string
	if input.ItemType == "file" {
		query = "user_id = ? AND file_id = ?"
	} else {
		query = "user_id = ? AND folder_id = ?"
	}
	result := database.DB.Where(query, userID, itemID).Delete(&models.Starred{})
	if result.Error != nil {
		log.Printf("UnstarItem: Failed to delete starred record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unstar item"})
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("UnstarItem: Item not starred")
		c.JSON(http.StatusNotFound, gin.H{"error": "item not starred"})
		return
	}
	log.Printf("UnstarItem: Successfully unstarred %s %s", input.ItemType, itemID)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s unstarred successfully", input.ItemType),
		"starred": false,
	})
}
func ToggleStarItem(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ToggleStarItem: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("ToggleStarItem: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var input struct {
		ItemID   string `json:"item_id" binding:"required"`
		ItemType string `json:"item_type" binding:"required,oneof=file folder"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("ToggleStarItem: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		log.Printf("ToggleStarItem: Failed to parse itemID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	if input.ItemType == "file" {
		var file models.File
		if err := database.DB.Where("id = ? AND user_id = ?", itemID, userID).First(&file).Error; err != nil {
			log.Printf("ToggleStarItem: File not found or doesn't belong to user: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
	} else if input.ItemType == "folder" {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", itemID, userID).First(&folder).Error; err != nil {
			log.Printf("ToggleStarItem: Folder not found or doesn't belong to user: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
	}
	var existingStarred models.Starred
	var query string
	if input.ItemType == "file" {
		query = "user_id = ? AND file_id = ?"
	} else {
		query = "user_id = ? AND folder_id = ?"
	}
	isStarred := database.DB.Where(query, userID, itemID).First(&existingStarred).Error == nil
	if isStarred {
		result := database.DB.Where(query, userID, itemID).Delete(&models.Starred{})
		if result.Error != nil {
			log.Printf("ToggleStarItem: Failed to delete starred record: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unstar item"})
			return
		}
		log.Printf("ToggleStarItem: Successfully unstarred %s %s", input.ItemType, itemID)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s unstarred successfully", input.ItemType),
			"starred": false,
		})
	} else {
		starred := models.Starred{
			UserID: userID,
		}
		if input.ItemType == "file" {
			starred.FileID = &itemID
		} else {
			starred.FolderID = &itemID
		}
		if err := database.DB.Create(&starred).Error; err != nil {
			log.Printf("ToggleStarItem: Failed to create starred record: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to star item"})
			return
		}
		log.Printf("ToggleStarItem: Successfully starred %s %s", input.ItemType, itemID)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s starred successfully", input.ItemType),
			"starred": true,
		})
	}
}
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
	log.Printf("ListStarredItems: Fetched %d starred items", len(starredItems))
	var response []gin.H
	for i, item := range starredItems {
		log.Printf("ListStarredItems: Processing item %d - FileID: %v, FolderID: %v", i, item.FileID, item.FolderID)
		if item.FileID != nil && item.File != nil {
			sizeStr := formatFileSize(item.File.Size)
			fileType := getFileTypeFromName(item.File.Name)
			response = append(response, gin.H{
				"id":         item.ID,
				"type":       fileType,
				"file_id":    item.FileID,
				"name":       item.File.Name,
				"size":       sizeStr,
				"mime_type":  item.File.MimeType,
				"created_at": item.CreatedAt,
			})
		} else if item.FolderID != nil && item.Folder != nil {
			log.Printf("ListStarredItems: Processing folder %s (%s)", item.Folder.Name, item.FolderID)
			var fileCount int64
			var subfolderCount int64
			if err := database.DB.Model(&models.File{}).Where("folder_id = ? AND deleted_at IS NULL", item.FolderID).Count(&fileCount).Error; err != nil {
				log.Printf("ListStarredItems: Error counting files in folder %s: %v", item.FolderID, err)
			}
			if err := database.DB.Model(&models.Folder{}).Where("parent_id = ? AND deleted_at IS NULL", item.FolderID).Count(&subfolderCount).Error; err != nil {
				log.Printf("ListStarredItems: Error counting subfolders in folder %s: %v", item.FolderID, err)
			}
			totalItems := fileCount + subfolderCount
			log.Printf("ListStarredItems: Folder %s has %d files and %d subfolders (total: %d)", item.Folder.Name, fileCount, subfolderCount, totalItems)
			response = append(response, gin.H{
				"id":           item.ID,
				"type":         "folder",
				"folder_id":    item.FolderID,
				"name":         item.Folder.Name,
				"size":         totalItems, 
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
func getFileTypeFromName(filename string) string {
	if filename == "" {
		return "document"
	}
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		ext = strings.ToLower(filename[idx+1:])
	}
	if ext == "" {
		return "document"
	}
	typeMap := map[string]string{
		"pdf": "pdf",
		"doc": "document",
		"docx": "document",
		"txt": "text",
		"rtf": "document",
		"odt": "document",
		"xls": "spreadsheet",
		"xlsx": "spreadsheet",
		"csv": "spreadsheet",
		"ods": "spreadsheet",
		"ppt": "presentation",
		"pptx": "presentation",
		"odp": "presentation",
		"jpg": "image",
		"jpeg": "image",
		"png": "image",
		"gif": "image",
		"bmp": "image",
		"webp": "image",
		"svg": "image",
		"ico": "image",
		"mp4": "video",
		"avi": "video",
		"mov": "video",
		"wmv": "video",
		"flv": "video",
		"webm": "video",
		"mp3": "audio",
		"wav": "audio",
		"flac": "audio",
		"aac": "audio",
		"zip": "archive",
		"rar": "archive",
		"7z": "archive",
		"tar": "archive",
		"gz": "archive",
		"js": "code",
		"ts": "code",
		"jsx": "code",
		"tsx": "code",
		"html": "code",
		"css": "code",
		"scss": "code",
		"sass": "code",
		"py": "code",
		"java": "code",
		"cpp": "code",
		"c": "code",
		"php": "code",
		"rb": "code",
		"go": "code",
		"rs": "code",
		"swift": "code",
		"kt": "code",
		"sql": "code",
		"json": "data",
		"xml": "data",
		"yaml": "data",
		"yml": "data",
		"exe": "executable",
		"msi": "executable",
		"dmg": "executable",
		"deb": "executable",
		"rpm": "executable",
	}
	if fileType, exists := typeMap[ext]; exists {
		return fileType
	}
	return "document"
} 