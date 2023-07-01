package serviceNats

import (
	"github.com/nats-io/stan.go"
	"orders/pkg/client/repositoryNats"
)

type Orders interface {
	CreateOrderItems(msg *stan.Msg)
}

type Service struct {
	Orders
}

func NewService(repo repositoryNats.Orders) *Service {
	return &Service{
		New(repo),
	}
}

func (s *Service) InitHandler() stan.MsgHandler {
	return s.Orders.CreateOrderItems
}
