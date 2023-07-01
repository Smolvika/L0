package service

import (
	"orders"
	"orders/pkg/server/repository"
)

type Orders interface {
	GetOrderById(orderId int) (orders.OrderItems, error)
}

type Service struct {
	Orders
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Orders: New(repo),
	}
}
