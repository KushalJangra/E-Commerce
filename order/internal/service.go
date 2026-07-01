package internal

import (
	"context"
	"time"

	"github.com/kushaljangra/e-commerce/order/models"
)

type Service interface {
	PostOrder(ctx context.Context, accountID uint64, totalPrice float64, products []*models.OrderedProduct) (*models.Order, error)
	GetOrdersForAccount(ctx context.Context, accountID uint64) ([]*models.Order, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, status string) error
}

type orderService struct {
	repository Repository
}

func NewOrderService(repository Repository) Service {
	return &orderService{repository}
}

func (service orderService) PostOrder(ctx context.Context, accountID uint64, totalPrice float64, products []*models.OrderedProduct) (*models.Order, error) {
	order := models.Order{
		AccountID:  accountID,
		TotalPrice: totalPrice,
		Products:   products,
		CreatedAt:  time.Now().UTC(),
	}
	err := service.repository.PutOrder(ctx, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (service orderService) GetOrdersForAccount(ctx context.Context, accountID uint64) ([]*models.Order, error) {
	return service.repository.GetOrdersForAccount(ctx, accountID)
}

func (service orderService) UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, paymnetStatus string) error {
	return service.repository.UpdateOrderPaymentStatus(ctx, orderId, paymnetStatus)
}
