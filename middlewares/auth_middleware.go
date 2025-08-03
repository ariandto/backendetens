package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/ariandto/backendetens/config"
	"github.com/gin-gonic/gin"
)

func CheckAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session cookie tidak ditemukan"})
			c.Abort()
			return
		}

		authClient, err := config.FirebaseApp.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inisialisasi Firebase Auth"})
			c.Abort()
			return
		}

		// Verifikasi session token
		token, err := authClient.VerifySessionCookieAndCheckRevoked(context.Background(), cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session tidak valid atau sudah dicabut"})
			c.Abort()
			return
		}

		// Periksa apakah email adalah admin
		email, ok := token.Claims["email"].(string)
		if !ok || email != os.Getenv("ADMIN_EMAIL") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses terbatas hanya untuk admin"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckAuthFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session")
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Cookie tidak ditemukan"})
			return
		}

		client, err := config.FirebaseApp.Auth(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Firebase tidak siap"})
			return
		}

		// Verifikasi token dari cookie
		token, err := client.VerifyIDToken(c, cookie.Value)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		// Simpan email dan UID di context
		c.Set("uid", token.UID)
		c.Set("email", token.Claims["email"])
		c.Set("name", token.Claims["name"])
		c.Set("admin", token.Claims["email"] == os.Getenv("ADMIN_EMAIL"))

		c.Next()
	}
}
