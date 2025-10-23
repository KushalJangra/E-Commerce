package model

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `json:"description,omitempty"`
	Products    []Product `gorm:"many2many:product_categories" json:"products,omitempty"`
}
