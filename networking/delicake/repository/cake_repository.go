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

func (r *CakeRepository) GetAll(page, pageSize int) ([]models.Cake, int64, error) {
	var cakes []models.Cake
	var total int64

	offset := (page - 1) * pageSize

	err := r.DB.Model(&models.Cake{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.Offset(offset).Limit(pageSize).Find(&cakes).Error
	return cakes, total, err
}
