package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kushaljangra/e-commerce/gateway/config"
	"github.com/kushaljangra/e-commerce/gateway/internal"
)

func main() {
	server, err := internal.NewServer(
		config.AccountURL,
		config.ProductURL,
		config.OrderURL,
		config.PaymentURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	engine := gin.Default()
	internal.RegisterRoutes(engine, server)

	log.Fatal(engine.Run(fmt.Sprintf(":%d", config.HTTPPort)))
}
