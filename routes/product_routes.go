package routes

import (
	"github.com/ariandto/backendetens/handlers"
	"github.com/ariandto/backendetens/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	// middleware auth dari cookie
	api := r.Group("/api/products", middlewares.CheckAuthFromCookie())
	{
		api.GET("", handlers.GetProducts)
		api.POST("/upload", middlewares.CheckAdminOnly(), handlers.UploadProduct)
		api.PUT("/:id", middlewares.CheckAdminOnly(), handlers.UpdateProduct)
		api.DELETE("/:id", middlewares.CheckAdminOnly(), handlers.DeleteProduct)
	}
}
