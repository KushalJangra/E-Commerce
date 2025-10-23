package service

import (
	"errors"
	"product-service/internal/model"
	"product-service/internal/repository"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func (s *ProductService) CreateProduct(name string, description string, price float64, stock int) error {
	product := model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		// Categories:  categories,
	}
	return s.Repo.CreateProduct(&product)
}

func (s *ProductService) GetProducts(searchQuery string) ([]model.Product, error) {
	products, err := s.Repo.GetProducts(searchQuery)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductById(id string) (*model.Product, error) {
	product, err := s.Repo.GetProductById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProductById(id string, updateData map[string]interface{}) (*model.Product, error) {

	if len(updateData) == 0 {
		return nil, errors.New("no fields to update")
	}

	product, err := s.Repo.UpdateProductById(id, updateData)

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProductById(id string) error {

	err := s.Repo.DeleteProductById(id)

	if err != nil {
		return err
	}

	return nil
}
