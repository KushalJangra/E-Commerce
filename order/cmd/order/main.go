package main

import (
	"log"
	"time"

	"github.com/kushaljangra/e-commerce/order/config"
	"github.com/kushaljangra/e-commerce/order/internal"
	"github.com/tinrab/retry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var repository internal.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Println(err)
		}
		repository, err = internal.NewPostgresRepository(db)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	log.Printf("Listening on port %d...", config.GrpcPort)
	service := internal.NewOrderService(repository)
	log.Fatal(internal.ListenGRPC(service, config.AccountServiceURL, config.ProductServiceURL, config.GrpcPort))
}
