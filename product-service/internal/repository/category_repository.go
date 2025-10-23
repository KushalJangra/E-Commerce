package repository

import (
	"errors"

	"product-service/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	return r.DB.Create(category).Error
}

func (r *CategoryRepository) GetCategoryById(id int) (*model.Category, error) {
	var category model.Category
	result := r.DB.Where("id = ?", id).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(id int, updatedCategory *model.Category) error {
	category, err := r.GetCategoryById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}
	category.Name = updatedCategory.Name
	category.Description = updatedCategory.Description
	category.Products = updatedCategory.Products
	return r.DB.Save(category).Error
}

func (r *CategoryRepository) DeleteCategory(id int) error {
	category, err := r.GetCategoryById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}
	return r.DB.Delete(category).Error
}
