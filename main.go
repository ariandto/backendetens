package main

import (
	"log"
	"os"
	"time"

	"github.com/ariandto/backendetens/config"
	"github.com/ariandto/backendetens/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found")
	}

	// Initialize Firebase Admin SDK
	config.InitFirebase()

	// Connect to Database
	config.ConnectDB()

	// Initialize Gin router
	router := gin.Default()

	// Serve static files (e.g., uploaded images)
	router.Static("/uploads", "./uploads")

	// Custom CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"https://etensports.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register all product-related routes
	routes.RegisterProductRoutes(router)

	routes.RegisterAuthRoutes(router)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5700"
	}
	log.Printf("🚀 Server running at http://localhost:%s", port)
	router.Run(":" + port)
}
