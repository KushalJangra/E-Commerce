package router

import (
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	productRepository := &repository.ProductRepository{DB: db}
	productService := &service.ProductService{Repo: *productRepository}
	productHandler := &handler.ProductHandler{ProductService: productService}

	r := gin.Default()
	// r.Use(middleware.Logging())
	// r.Use(middleware.ErrorHandler())

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetProducts)
	r.GET("/product/:id", productHandler.GetProductById)
	r.PATCH("/product/:id", productHandler.UpdateProductById)
	r.DELETE("/product/:id", productHandler.DeleteProductById)
	return r
}
