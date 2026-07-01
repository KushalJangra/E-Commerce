package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kushaljangra/e-commerce/pkg/middleware"
)

func RegisterRoutes(engine *gin.Engine, server *Server) {
	engine.Use(middleware.GinContextToContextMiddleware())
	engine.Use(middleware.AuthorizeJWT())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group("/api")
	{
		api.POST("/auth/register", server.Register)
		api.POST("/auth/login", server.Login)

		api.GET("/accounts", server.ListAccounts)
		api.GET("/accounts/:id", server.GetAccount)
		api.GET("/accounts/:id/orders", server.GetAccountOrders)

		api.GET("/products", server.ListProducts)
		api.GET("/products/:id", server.GetProduct)
		api.POST("/products", server.CreateProduct)
		api.PUT("/products/:id", server.UpdateProduct)
		api.DELETE("/products/:id", server.DeleteProduct)

		api.POST("/orders", server.CreateOrder)

		api.POST("/payments/customer-portal", server.CreateCustomerPortalSession)
		api.POST("/payments/checkout", server.CreateCheckoutSession)
	}
}
