package repository

import (
	"errors"
	"orders"
)

type OrderCache struct {
	*Cache
}

func New(cache *Cache) *OrderCache {
	return &OrderCache{cache}
}

func (c *OrderCache) GetOrderById(OrderId int) (orders.OrderItems, error) {
	value, ok := c.GetCache(OrderId)
	if !ok {
		return value, errors.New("there is no such Id")
	}

	return value, nil
}
