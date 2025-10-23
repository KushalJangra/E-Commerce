package repository

import (
	"errors"

	"product-service/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r *ProductRepository) CreateProduct(product *model.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) GetProductById(id string) (*model.Product, error) {
	var product model.Product
	result := r.DB.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) GetProducts(searchQuery string) ([]model.Product, error) {
	var products []model.Product

	query := r.DB
	if searchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) UpdateProductById(id string, updateData map[string]interface{}) (*model.Product, error) {
	product, err := r.GetProductById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
	}

	if err := r.DB.Model(&product).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) DeleteProductById(id string) error {
	product, err := r.GetProductById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}
	return r.DB.Delete(product).Error
}
