package service

import (
	"orders"
	"orders/pkg/server/repository"
)

type OrderService struct {
	repo repository.Orders
}

func New(repo repository.Orders) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s OrderService) GetOrderById(OrderId int) (orders.OrderItems, error) {
	return s.repo.GetOrderById(OrderId)
}
