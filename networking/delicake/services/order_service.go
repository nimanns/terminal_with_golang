package service

import (
	"delicake/models"
	"delicake/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return s.repo.Create(order)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.repo.GetAll()
}
