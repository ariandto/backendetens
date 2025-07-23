package main

import (
	"log"
	"os"
	"time"

	"github.com/ariandto/backendazko/config"
	"github.com/ariandto/backendazko/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4600"
	}

	os.MkdirAll("uploads", os.ModePerm)
	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173","https://lacakazko.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	// Tambahkan prefix group /api
	api := r.Group("/api")
	{
		api.POST("/photos", handlers.UploadPhoto)
		api.GET("/photos", handlers.GetPhotos)
		api.GET("/users", handlers.GetAllUsers)
		api.PUT("/photos/:id", handlers.UpdatePhoto)
		api.DELETE("/photos/:id", handlers.DeletePhoto)
		api.POST("/users", handlers.CreateUserWithPhoto) // ✅ tambah endpoint user
	}

	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}