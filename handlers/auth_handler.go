package handlers

import (
	"net/http"
	"os"

	"github.com/ariandto/backendetens/config"
	"github.com/gin-gonic/gin"
)

func LoginWithFirebase(c *gin.Context) {
	var req struct {
		IDToken string `json:"idToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.IDToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID token dibutuhkan"})
		return
	}

	client, err := config.FirebaseApp.Auth(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inisialisasi Firebase Auth"})
		return
	}

	// Verifikasi token Firebase
	token, err := client.VerifyIDToken(c, req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	// Set session cookie (HTTP-only)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    req.IDToken,
		HttpOnly: true,
		Secure:   false, // set true jika menggunakan HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 7 hari
	})

	// Kirim respons sukses (tanpa token)
	c.JSON(http.StatusOK, gin.H{
		"uid":   token.UID,
		"email": email,
		"name":  name,
		"admin": email == os.Getenv("ADMIN_EMAIL"), // misal:
	})
}

func GetMe(c *gin.Context) {
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Belum login"})
		return
	}

	authClient, err := config.FirebaseApp.Auth(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inisialisasi Firebase"})
		return
	}

	token, err := authClient.VerifySessionCookie(c, sessionCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session tidak valid"})
		return
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	c.JSON(http.StatusOK, gin.H{
		"uid":   token.UID,
		"email": email,
		"name":  name,
		"admin": email == os.Getenv("ADMIN_EMAIL"),
	})
}
