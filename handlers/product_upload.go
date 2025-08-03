package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ariandto/backendetens/config"
	"github.com/ariandto/backendetens/models"
	"github.com/gin-gonic/gin"
)

func UploadProduct(c *gin.Context) {
	name := c.PostForm("name")
	size := c.PostForm("size")
	stock := c.PostForm("stock")
	price := c.PostForm("price")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar diperlukan"})
		return
	}

	// Simpan file ke folder uploads/
	filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	dst := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// Simpan data ke database (panggil prosedur atau langsung pakai GORM)
	product := models.Product{
		Name:     name,
		Size:     size,
		Stock:    atoi(stock),
		Price:    atof(price),
		ImageURL: "/" + dst, // misal: /uploads/jersey-a.jpg
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil ditambahkan", "product": product})
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
