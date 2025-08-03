package handlers

import (
	"net/http"
	"os"

	"github.com/ariandto/backendetens/config"
	"github.com/ariandto/backendetens/models"
	"github.com/gin-gonic/gin"
)

// ✅ Fungsi pembantu untuk verifikasi admin tanpa middleware
func isAdmin(c *gin.Context) bool {
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Cookie tidak ditemukan"})
		return false
	}

	authClient, err := config.FirebaseApp.Auth(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase Auth tidak tersedia"})
		return false
	}

	token, err := authClient.VerifySessionCookie(c, sessionCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session tidak valid"})
		return false
	}

	email, _ := token.Claims["email"].(string)
	if email != os.Getenv("ADMIN_EMAIL") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengakses fitur ini"})
		return false
	}

	return true
}

// Ambil semua produk via stored procedure
func GetProducts(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var products []models.Product

	result := config.DB.Raw("CALL sp_get_products()").Scan(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Tambah produk baru
func CreateProduct(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Exec("CALL sp_create_product(?, ?, ?, ?, ?)",
		p.Name, p.Size, p.Stock, p.Price, p.ImageURL)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan"})
}

// Perbarui produk
func UpdateProduct(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	id := c.Param("id")

	// Ambil field dari FormData
	name := c.PostForm("name")
	size := c.PostForm("size")
	stock := c.PostForm("stock")
	price := c.PostForm("price")
	imageURL := c.PostForm("image_url") // jika tidak upload baru

	// Cek apakah user upload gambar baru
	file, err := c.FormFile("image")
	if err == nil {
		// ✅ Simpan file ke folder uploads
		dst := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload file"})
			return
		}
		imageURL = "/" + dst // update imageURL jika upload berhasil
	}

	// Eksekusi stored procedure
	result := config.DB.Exec("CALL sp_update_product(?, ?, ?, ?, ?, ?)",
		id, name, size, stock, price, imageURL)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diperbarui"})
}

// Hapus produk
func DeleteProduct(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	id := c.Param("id")

	result := config.DB.Exec("CALL sp_delete_product(?)", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
