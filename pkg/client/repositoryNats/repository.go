package repositoryNats

import (
	"github.com/jmoiron/sqlx"
	"orders"
	"orders/pkg/server/repository"
)

type Orders interface {
	CreateOrderItems(Order *orders.OrderItems) error
}
type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB, c *repository.Cache) *Repository {
	return &Repository{
		Orders: New(c, db),
	}
}
