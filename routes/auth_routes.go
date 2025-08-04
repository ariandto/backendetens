package routes

import (
	"github.com/ariandto/backendetens/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/api/login", handlers.LoginWithFirebase)
	r.POST("api/logout", handlers.Logout)
	r.GET("/api/me", handlers.GetMe)
	r.DELETE("/:id", handlers.DeleteProduct)
}
