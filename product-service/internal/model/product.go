package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `json:"description,omitempty"`
	Price       float64    `gorm:"not null" json:"price"`
	Stock       int        `gorm:"not null" json:"stock"`
	Categories  []Category `gorm:"many2many:product_categories" json:"categories,omitempty"`
}

// BeforeCreate hook to set UUID
func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}
