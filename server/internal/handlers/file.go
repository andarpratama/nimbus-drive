package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	userID := c.GetUint("userID")
	log.Println("userID", userID)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	// Get folder_id from form data (optional)
	folderIDStr := c.PostForm("folder_id")
	var folderID *uint
	if folderIDStr != "" {
		if id, err := strconv.ParseUint(folderIDStr, 10, 32); err == nil {
			uintID := uint(id)
			folderID = &uintID
			
			// Verify folder exists and belongs to user
			var folder models.Folder
			if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "folder not found"})
				return
			}
		}
	}

	// Remove all whitespace from the filename
	cleanFilename := ""
	for _, r := range file.Filename {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			cleanFilename += string(r)
		}
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}

	timestamp := time.Now().Unix()
	dst := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", timestamp, cleanFilename))

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	dbFile := models.File{
		Name:     cleanFilename,
		Path:     dst,
		Size:     file.Size, // Capture actual file size
		UserID:   userID,
		FolderID: folderID,
	}

	if err := database.DB.Create(&dbFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded", "file": dbFile})
}

func DownloadFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.FileAttachment(file.Path, file.Name)
}

// ServeImage handles GET /files/:id/image to serve images publicly
func ServeImage(c *gin.Context) {
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.
		Where("id = ?", fileID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Check if file is an image
	ext := strings.ToLower(filepath.Ext(file.Name))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	isImage := false
	for _, imgExt := range imageExts {
		if ext == imgExt {
			isImage = true
			break
		}
	}

	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not an image"})
		return
	}

	// Set appropriate headers for image serving
	c.Header("Content-Type", "image/"+strings.TrimPrefix(ext, "."))
	c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	c.File(file.Path)
}

func ListFiles(c *gin.Context) {
	userID := c.GetUint("userID")
	folderID := c.Query("folder_id")

	var files []models.File
	query := database.DB.Where("user_id = ?", userID).Preload("Folder")

	// If folder_id is provided, filter by folder
	if folderID != "" {
		if id, err := strconv.ParseUint(folderID, 10, 32); err == nil {
			query = query.Where("folder_id = ?", id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder_id"})
			return
		}
	} else {
		// If no folder_id, show files in root (folder_id is NULL)
		query = query.Where("folder_id IS NULL")
	}

	if err := query.Order("created_at DESC").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func DeleteFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Soft delete
	if err := database.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file moved to trash"})
}

func GetTrashedFiles(c *gin.Context) {
	userID := c.GetUint("userID")

	var trashed []models.File
	if err := database.DB.Unscoped(). // include soft-deleted
						Where("user_id = ? AND deleted_at IS NOT NULL", userID).
						Find(&trashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trash"})
		return
	}

	c.JSON(http.StatusOK, trashed)
}

func MoveFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var input struct {
		FolderID *uint `json:"folder_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the file
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// If folder_id is provided, verify it exists and belongs to the user
	if input.FolderID != nil {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.FolderID, userID).First(&folder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "folder not found"})
			return
		}
	}

	// Update file's folder
	file.FolderID = input.FolderID

	if err := database.DB.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file moved", "file": file})
}

// RestoreFile handles POST /files/:id/restore to restore a soft-deleted file
func RestoreFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found in trash"})
		return
	}

	// Restore the file by setting deleted_at to NULL
	if err := database.DB.Unscoped().Model(&file).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to restore file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file restored", "file": file})
}

// PermanentlyDeleteFile handles DELETE /files/:id/permanent to permanently delete a file
func PermanentlyDeleteFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found in trash"})
		return
	}

	// Delete the actual file from disk
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file from disk"})
		return
	}

	// Permanently delete from database
	if err := database.DB.Unscoped().Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to permanently delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file permanently deleted"})
}

func RenameFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// Clean the filename (remove whitespace)
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

	var file models.File
	if err := database.DB.
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Check if file with same name already exists in the same folder
	var existingFile models.File
	query := database.DB.Where("name = ? AND user_id = ? AND id != ?", cleanName, userID, fileID)
	if file.FolderID != nil {
		query = query.Where("folder_id = ?", file.FolderID)
	} else {
		query = query.Where("folder_id IS NULL")
	}
	
	if err := query.First(&existingFile).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "file with this name already exists"})
		return
	}

	// Get the file extension from the original path
	oldExt := filepath.Ext(file.Path)
	newExt := filepath.Ext(cleanName)
	
	// If the extension changed, use the new extension, otherwise keep the old one
	finalExt := oldExt
	if newExt != "" {
		finalExt = newExt
	}
	
	// Create new filename with proper extension
	newFilename := cleanName
	if !strings.HasSuffix(newFilename, finalExt) {
		newFilename = newFilename + finalExt
	}
	
	// Validate that the new filename is not empty after extension handling
	if newFilename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}
	
	// Create new path by replacing the filename in the old path
	oldDir := filepath.Dir(file.Path)
	newPath := filepath.Join(oldDir, newFilename)
	
	// Check if the new path already exists on disk
	if _, err := os.Stat(newPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "file with this name already exists on disk"})
		return
	}
	
	// Rename the physical file on disk
	log.Printf("Renaming file from %s to %s", file.Path, newPath)
	if err := os.Rename(file.Path, newPath); err != nil {
		log.Printf("Failed to rename physical file from %s to %s: %v", file.Path, newPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename file on disk"})
		return
	}
	
	// Update the file name, path and UpdatedAt timestamp
	file.Name = newFilename
	file.Path = newPath
	file.UpdatedAt = time.Now()

	if err := database.DB.Save(&file).Error; err != nil {
		// If database update fails, try to revert the file rename
		if revertErr := os.Rename(newPath, file.Path); revertErr != nil {
			log.Printf("Failed to revert file rename: %v", revertErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file renamed successfully", "file": file})
}
