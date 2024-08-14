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

func (s *CakeService) GetAllCakes() ([]models.Cake, error) {
	return s.repo.GetAll()
}

// service/order_service.go
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
