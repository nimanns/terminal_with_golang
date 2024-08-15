package service

import (
	"delicake/models"
	"delicake/repository"
)

type CakeService struct {
	repo *repository.CakeRepository
}

func NewCakeService(repo *repository.CakeRepository) *CakeService {
	return &CakeService{repo: repo}
}

func (s *CakeService) CreateCake(cake *models.Cake) error {
	return s.repo.Create(cake)
}

func (s *CakeService) GetAllCakes(page, pageSize int) ([]models.Cake, int64, error) {
	return s.repo.GetAll(page, pageSize)
}

