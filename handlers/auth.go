package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/ariandto/backendazko/config"
	"github.com/ariandto/backendazko/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func GoogleLogin(c *gin.Context) {
	var body struct {
		IDToken string `json:"id_token"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), body.IDToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	photo := payload.Claims["picture"].(string) // ✅ ambil foto Google

	var user models.User
	config.DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		user = models.User{
			Name:  name,
			Photo: photo,
		}
		config.DB.Create(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}
