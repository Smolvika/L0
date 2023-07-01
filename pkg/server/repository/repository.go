package repository

import "orders"

type Orders interface {
	GetOrderById(orderId int) (orders.OrderItems, error)
}

type Repository struct {
	Orders
}

func NewRepository(c *Cache) *Repository {
	return &Repository{
		Orders: New(c),
	}
}
