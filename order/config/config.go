package config

import "os"

var (
	DatabaseURL        string
	AccountServiceURL  string
	ProductServiceURL  string
	BootstrapServers   string
)

const (
	GrpcPort int = 8080
)

func init() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	AccountServiceURL = os.Getenv("ACCOUNT_SERVICE_URL")
	ProductServiceURL = os.Getenv("PRODUCT_SERVICE_URL")
	BootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
}
