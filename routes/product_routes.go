package routes

import (
	"github.com/ariandto/backendetens/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	api := r.Group("/api/products")
	{
		api.GET("", handlers.GetProducts)
		api.POST("", handlers.CreateProduct) // gunakan CreateProduct, bukan upload
		api.PUT("/:id", handlers.UpdateProduct)
		api.DELETE("/:id", handlers.DeleteProduct)
	}
}