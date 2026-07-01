package internal

import (
	account "github.com/kushaljangra/e-commerce/account/client"
	order "github.com/kushaljangra/e-commerce/order/client"
	payment "github.com/kushaljangra/e-commerce/payment/client"
	product "github.com/kushaljangra/e-commerce/product/client"
)

type Server struct {
	Account *account.Client
	Product *product.Client
	Order   *order.Client
	Payment *payment.Client
}

func NewServer(accountURL, productURL, orderURL, paymentURL string) (*Server, error) {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return nil, err
	}

	productClient, err := product.NewClient(productURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		productClient.Close()
		return nil, err
	}

	paymentClient, err := payment.NewClient(paymentURL)
	if err != nil {
		accountClient.Close()
		productClient.Close()
		orderClient.Close()
		return nil, err
	}

	return &Server{
		Account: accountClient,
		Product: productClient,
		Order:   orderClient,
		Payment: paymentClient,
	}, nil
}

func (s *Server) Close() {
	s.Account.Close()
	s.Product.Close()
	s.Order.Close()
	s.Payment.Close()
}
