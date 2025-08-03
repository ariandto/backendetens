package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ariandto/backendetens/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found. Using environment variables directly.")
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}

	// Auto migrate hanya untuk Product
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("❌ Gagal migrasi tabel Product:", err)
	}

	DB = db
	log.Println("✅ Terkoneksi ke MySQL & migrasi tabel Product selesai")
	return db
}
