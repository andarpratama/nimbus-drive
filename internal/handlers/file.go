package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	userID := c.GetUint("user_id") // assuming JWT middleware sets this

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
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
		Name:   cleanFilename,
		Path:   dst,
		UserID: userID,
	}

	if err := database.DB.Create(&dbFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded", "file": dbFile})
}

func DownloadFile(c *gin.Context) {
	userID := c.GetUint("user_id")
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

func ListFiles(c *gin.Context) {
	userID := c.GetUint("user_id")

	var files []models.File
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
