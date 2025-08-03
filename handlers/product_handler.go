package handlers

import (
	"net/http"

	"github.com/ariandto/backendetens/config"
	"github.com/ariandto/backendetens/models"
	"github.com/gin-gonic/gin"
)

// Ambil semua produk via stored procedure
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Panggil stored procedure
	result := config.DB.Raw("CALL sp_get_products()").Scan(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Tambah produk baru
func CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Langsung gunakan stored procedure
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
	id := c.Param("id")
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Exec("CALL sp_update_product(?, ?, ?, ?, ?, ?)",
		id, p.Name, p.Size, p.Stock, p.Price, p.ImageURL)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diperbarui"})
}

// Hapus produk
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Exec("CALL sp_delete_product(?)", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
