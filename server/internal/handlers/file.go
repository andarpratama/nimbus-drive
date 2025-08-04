package handlers
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)
func UploadFile(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UploadFile: Panic recovered: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	log.Println("UploadFile: Starting file upload")
	userIDStr := c.GetString("userID")
	log.Printf("UploadFile: userID from context: %s", userIDStr)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("UploadFile: Failed to parse userID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	log.Printf("UploadFile: Parsed userID: %s", userID)
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("UploadFile: Failed to get form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}
	log.Printf("UploadFile: Received file: %s, size: %d", file.Filename, file.Size)
	folderIDStr := c.PostForm("folder_id")
	log.Printf("UploadFile: folder_id from form: %s", folderIDStr)
	var folderID *uuid.UUID
	if folderIDStr != "" {
		if id, err := uuid.Parse(folderIDStr); err == nil {
			folderID = &id
			log.Printf("UploadFile: Parsed folderID: %s", folderID)
			var folder models.Folder
			if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
				log.Printf("UploadFile: Folder not found or doesn't belong to user: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "folder not found"})
				return
			}
			log.Printf("UploadFile: Folder verified: %s", folder.Name)
		} else {
			log.Printf("UploadFile: Failed to parse folderID: %v", err)
		}
	} else {
		log.Printf("UploadFile: No folder_id provided, uploading to root")
	}
	cleanFilename := ""
	for _, r := range file.Filename {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			cleanFilename += string(r)
		}
	}
	log.Printf("UploadFile: Clean filename: %s", cleanFilename)
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if !filepath.IsAbs(uploadDir) {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Printf("UploadFile: Failed to get current directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to determine upload directory"})
			return
		}
		uploadDir = filepath.Join(currentDir, uploadDir)
	}
	log.Printf("UploadFile: Upload directory path: %s", uploadDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		if os.IsExist(err) {
			log.Printf("UploadFile: Upload directory already exists: %s", uploadDir)
		} else {
			log.Printf("UploadFile: Failed to create upload directory %s: %v", uploadDir, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
			return
		}
	} else {
		log.Printf("UploadFile: Upload directory created successfully: %s", uploadDir)
	}
	if info, err := os.Stat(uploadDir); err != nil {
		log.Printf("UploadFile: Failed to stat upload directory %s: %v", uploadDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload directory not accessible"})
		return
	} else if !info.IsDir() {
		log.Printf("UploadFile: Upload path is not a directory: %s", uploadDir)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload path is not a directory"})
		return
	}
	log.Printf("UploadFile: Upload directory verified successfully")
	timestamp := time.Now().Unix()
	dst := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", timestamp, cleanFilename))
	log.Printf("UploadFile: Destination path: %s", dst)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		log.Printf("UploadFile: Failed to save uploaded file to %s: %v", dst, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	log.Printf("UploadFile: File saved successfully to %s", dst)
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		ext := strings.ToLower(filepath.Ext(cleanFilename))
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".pdf":
			mimeType = "application/pdf"
		case ".txt":
			mimeType = "text/plain"
		case ".doc", ".docx":
			mimeType = "application/msword"
		case ".xls", ".xlsx":
			mimeType = "application/vnd.ms-excel"
		case ".ppt", ".pptx":
			mimeType = "application/vnd.ms-powerpoint"
		default:
			mimeType = "application/octet-stream"
		}
	}
	dbFile := models.File{
		Name:     cleanFilename,
		Path:     dst,
		Size:     file.Size,
		MimeType: mimeType,
		UserID:   userID,
		FolderID: folderID,
	}
	log.Printf("UploadFile: Creating database record for file: %s", dbFile.Name)
	if err := database.DB.Create(&dbFile).Error; err != nil {
		log.Printf("UploadFile: Failed to save file metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}
	log.Printf("UploadFile: Database record created successfully, file ID: %s", dbFile.ID)
	c.JSON(http.StatusOK, gin.H{"message": "file uploaded", "file": dbFile})
}
func DownloadFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
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
func ServeImage(c *gin.Context) {
	fileID := c.Param("id")
	var file models.File
	if err := database.DB.
		Where("id = ?", fileID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
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
	c.Header("Content-Type", "image/"+strings.TrimPrefix(ext, "."))
	c.Header("Cache-Control", "public, max-age=31536000") 
	c.File(file.Path)
}
func ListFiles(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	folderID := c.Query("folder_id")
	var files []models.File
	query := database.DB.Where("user_id = ?", userID).Preload("Folder")
	if folderID != "" {
		if id, err := uuid.Parse(folderID); err == nil {
			query = query.Where("folder_id = ?", id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder_id"})
			return
		}
	} else {
		query = query.Where("folder_id IS NULL")
	}
	if err := query.Order("created_at DESC").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch files"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}
func DeleteFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	fileID := c.Param("id")
	var file models.File
	if err := database.DB.
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if err := database.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file moved to trash"})
}
func GetTrashedFiles(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var trashed []models.File
	if err := database.DB.Unscoped(). 
						Where("user_id = ? AND deleted_at IS NOT NULL", userID).
						Find(&trashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trash"})
		return
	}
	c.JSON(http.StatusOK, trashed)
}
func MoveFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	fileID := c.Param("id")
	var input struct {
		FolderID *uuid.UUID `json:"folder_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if input.FolderID != nil {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.FolderID, userID).First(&folder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "folder not found"})
			return
		}
	}
	file.FolderID = input.FolderID
	if err := database.DB.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file moved", "file": file})
}
func RestoreFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	fileID := c.Param("id")
	var file models.File
	if err := database.DB.Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found in trash"})
		return
	}
	if err := database.DB.Unscoped().Model(&file).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to restore file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file restored", "file": file})
}
func PermanentlyDeleteFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	fileID := c.Param("id")
	var file models.File
	if err := database.DB.Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found in trash"})
		return
	}
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file from disk"})
		return
	}
	if err := database.DB.Unscoped().Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to permanently delete file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file permanently deleted"})
}
func RenameFile(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	fileID := c.Param("id")
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
	var file models.File
	if err := database.DB.
		Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
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
	oldExt := filepath.Ext(file.Path)
	newExt := filepath.Ext(cleanName)
	finalExt := oldExt
	if newExt != "" {
		finalExt = newExt
	}
	newFilename := cleanName
	if !strings.HasSuffix(newFilename, finalExt) {
		newFilename = newFilename + finalExt
	}
	if newFilename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}
	oldDir := filepath.Dir(file.Path)
	newPath := filepath.Join(oldDir, newFilename)
	if _, err := os.Stat(newPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "file with this name already exists on disk"})
		return
	}
	log.Printf("Renaming file from %s to %s", file.Path, newPath)
	if err := os.Rename(file.Path, newPath); err != nil {
		log.Printf("Failed to rename physical file from %s to %s: %v", file.Path, newPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename file on disk"})
		return
	}
	file.Name = newFilename
	file.Path = newPath
	file.UpdatedAt = time.Now()
	if err := database.DB.Save(&file).Error; err != nil {
		if revertErr := os.Rename(newPath, file.Path); revertErr != nil {
			log.Printf("Failed to revert file rename: %v", revertErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file renamed successfully", "file": file})
}

func GetRecentFiles(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "20")
	limit := 20
	if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
		limit = parsedLimit
	}

	var files []models.File
	query := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("updated_at DESC").
		Limit(limit)

	if err := query.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recent files"})
		return
	}

	// Prepare response with file details
	var response []gin.H
	for _, file := range files {
		fileInfo := gin.H{
			"id":         file.ID,
			"name":       file.Name,
			"size":       file.Size,
			"mime_type":  file.MimeType,
			"created_at": file.CreatedAt,
			"updated_at": file.UpdatedAt,
			"folder_id":  file.FolderID,
		}
		response = append(response, fileInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"total": len(response),
		"limit": limit,
	})
}
