package repository

import (
	"delicake/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Preload("Cake").Find(&orders).Error
	return orders, err
}
