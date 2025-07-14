package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ariandto/backendazko/config"
	"github.com/ariandto/backendazko/models"

	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	name := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	path := filepath.Join("uploads", name)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	photo := models.Photo{
		Filename:  file.Filename,
		Path:      path,
		CreatedAt: time.Now(),
	}

	config.DB.Create(&photo)
	c.JSON(http.StatusOK, photo)
}

func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	config.DB.Find(&photos)
	c.JSON(http.StatusOK, photos)
}

func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo

	if err := config.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New file required"})
		return
	}

	_ = os.Remove(photo.Path)

	name := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	path := filepath.Join("uploads", name)
	c.SaveUploadedFile(file, path)

	photo.Filename = file.Filename
	photo.Path = path
	config.DB.Save(&photo)

	c.JSON(http.StatusOK, photo)
}

func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo

	if err := config.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	_ = os.Remove(photo.Path)
	config.DB.Delete(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
