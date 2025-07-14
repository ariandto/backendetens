package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ariandto/backendazko/config"
	"github.com/ariandto/backendazko/models"
	"github.com/gin-gonic/gin"
)

func CreateUserWithPhoto(c *gin.Context) {
	var user models.User

	// Ambil data dari form
	user.NIK = c.PostForm("nik")
	user.Name = c.PostForm("name")
	user.Department = c.PostForm("department")
	user.Shift = c.PostForm("shift")
	user.Phone = c.PostForm("phone")

	// === DEBUGGING: Cetak input dari form ===
	fmt.Println("DEBUG - Diterima dari form:")
	fmt.Println("NIK:", user.NIK)
	fmt.Println("Name:", user.Name)
	fmt.Println("Department:", user.Department)
	fmt.Println("Shift:", user.Shift)
	fmt.Println("Phone:", user.Phone)

	// Ambil file foto
	file, err := c.FormFile("photo")
	if err != nil {
		fmt.Println("DEBUG - Foto tidak ditemukan:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto tidak ditemukan"})
		return
	}

	filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	path := filepath.Join("uploads", filename)

	// Simpan file
	if err := c.SaveUploadedFile(file, path); err != nil {
		fmt.Println("DEBUG - Gagal simpan foto:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan foto"})
		return
	}
	user.Photo = path

	// Simpan ke database
	if err := config.DB.Create(&user).Error; err != nil {
		fmt.Println("DEBUG - Gagal menyimpan user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user"})
		return
	}

	fmt.Println("DEBUG - User berhasil disimpan:", user)

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat user"})
		return
	}
	c.JSON(http.StatusOK, users)
}
