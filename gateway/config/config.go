package config

import "os"

var (
	AccountURL string
	ProductURL string
	OrderURL   string
	PaymentURL string
)

const HTTPPort = 8080

func init() {
	AccountURL = os.Getenv("ACCOUNT_SERVICE_URL")
	ProductURL = os.Getenv("PRODUCT_SERVICE_URL")
	OrderURL = os.Getenv("ORDER_SERVICE_URL")
	PaymentURL = os.Getenv("PAYMENT_SERVICE_URL")
}
