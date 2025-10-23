package model

type ProductCategory struct {
	CategoryID uint `gorm:"primaryKey"`
	ProductID  uint `gorm:"primaryKey"`
}
