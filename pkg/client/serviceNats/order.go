package serviceNats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"orders"
	"orders/pkg/client/repositoryNats"
)

type OrderService struct {
	repo repositoryNats.Orders
}

func New(repo repositoryNats.Orders) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrderItems(msg *stan.Msg) {
	order := new(orders.OrderItems)
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		return
	}
	if err := s.repo.CreateOrderItems(order); err != nil {
		log.Fatalf(err.Error())
	}
}
