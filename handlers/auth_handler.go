package handlers

import (
	"net/http"
	"os"

	"time"

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

	// Verifikasi ID token
	token, err := client.VerifyIDToken(c, req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Buat session cookie dari ID token
	expiresIn := time.Hour * 24 * 5 // 5 hari
	sessionCookie, err := client.SessionCookie(c, req.IDToken, expiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat session cookie"})
		return
	}

	// Set cookie HTTP-Only
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    sessionCookie, // 🔁 session cookie, bukan ID token lagi
		MaxAge:   int(expiresIn.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // set true jika HTTPS
		Path:     "/",
	})

	c.JSON(http.StatusOK, gin.H{
		"uid":   token.UID,
		"email": token.Claims["email"],
		"name":  token.Claims["name"],
		"admin": token.Claims["email"] == os.Getenv("ADMIN_EMAIL"),
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

func Logout(c *gin.Context) {
	// Hapus cookie 'session' dengan overwrite kosong
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    "",
		MaxAge:   -1, // menghapus cookie
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // set true jika HTTPS
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}
