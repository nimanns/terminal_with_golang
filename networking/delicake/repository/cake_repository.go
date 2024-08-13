package repository

import (
	"delicake/models"
	"gorm.io/gorm"
)

type CakeRepository struct {
	DB *gorm.DB
}

func (r *CakeRepository) Create(cake *models.Cake) error {
	return r.DB.Create(cake).Error
}

func (r *CakeRepository) GetAll() ([]models.Cake, error) {
	var cakes []models.Cake
	err := r.DB.Find(&cakes).Error
	return cakes, err
}
