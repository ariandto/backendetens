package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ariandto/backendetens/config"
	"github.com/ariandto/backendetens/handlers"
	"github.com/ariandto/backendetens/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4600"
	}

	// Buat folder uploads jika belum ada
	os.MkdirAll("uploads", os.ModePerm)

	// Koneksi ke database
	db := config.ConnectDB() // Ensure ConnectDB returns *gorm.DB

	// Auto migrate Visitor
	if err := db.AutoMigrate(&models.Visitor{}); err != nil {
		log.Fatalf("AutoMigrate Visitor failed: %v", err)
	}

	// Load timezone Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

	// Inisialisasi router Gin
	r := gin.Default()

	// 📌 Endpoint bebas CORS untuk /api/visit (manual, sebelum middleware global)
	r.POST("/api/visit", func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = c.Request.Header.Get("X-Forwarded-For")
		}
		if ip == "" {
			log.Println("⚠️ Tidak bisa mendapatkan IP pengunjung")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak bisa mendeteksi IP"})
			return
		}

		now := time.Now().In(loc)
		today := now.Format("2006-01-02")

		var visitor models.Visitor
		result := db.Where("ip = ? AND DATE(visited_date) = ?", ip, today).First(&visitor)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				newVisitor := models.Visitor{
					IP:          ip,
					VisitedDate: now,
				}
				if err := db.Create(&newVisitor).Error; err != nil {
					if strings.Contains(err.Error(), "Duplicate entry") {
						log.Printf("⚠️ Duplikat terdeteksi dari IP: %s (%s)", ip, today)
						c.JSON(http.StatusOK, gin.H{"message": "Already visited today", "ip": ip})
						return
					}
					log.Printf("❌ Gagal mencatat kunjungan: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
					return
				}
				log.Printf("✅ IP baru tercatat: %s (%s)", ip, today)
				c.JSON(http.StatusOK, gin.H{"message": "Visitor recorded", "ip": ip})
			} else {
				log.Printf("❌ Error database saat cek pengunjung: %v", result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
		} else {
			log.Printf("ℹ️ Pengunjung sudah tercatat hari ini: %s", ip)
			c.JSON(http.StatusOK, gin.H{"message": "Already visited today", "ip": ip})
		}
	})

	// 🌐 Middleware CORS untuk endpoint lainnya
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://lacakazko.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static folder untuk foto
	r.Static("/uploads", "./uploads")

	// Group prefix /api
	api := r.Group("/api")
	{
		// Endpoint visitor count (pakai CORS)
		api.GET("/visitors/count", func(c *gin.Context) {
			now := time.Now().In(loc)
			today := now.Format("2006-01-02")

			var count int64
			db.Model(&models.Visitor{}).Where("DATE(visited_date) = ?", today).Count(&count)
			c.JSON(200, gin.H{"total_unique_visitors_today": count})
		})

		api.GET("/visitors/view", func(c *gin.Context) {
			var visitors []models.Visitor

			// Ambil semua data pengunjung dari DB, urut dari terbaru
			if err := db.Order("visited_date DESC").Find(&visitors).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Gagal mengambil data pengunjung",
				})
				return
			}

			// Kembalikan sebagai JSON
			c.JSON(http.StatusOK, gin.H{
				"total":    len(visitors),
				"visitors": visitors,
			})
		})

		// Upload & user routes
		api.POST("/photos", handlers.UploadPhoto)
		api.GET("/photos", handlers.GetPhotos)
		api.PUT("/photos/:id", handlers.UpdatePhoto)
		api.DELETE("/photos/:id", handlers.DeletePhoto)
		api.POST("/users", handlers.CreateUserWithPhoto)
		api.GET("/users", handlers.GetAllUsers)
	}

	log.Println("🚀 Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
