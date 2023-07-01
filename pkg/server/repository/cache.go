package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"orders"
	"strconv"
	"sync"
)

type Cache struct {
	m  map[int]orders.OrderItems
	mx sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		m:  make(map[int]orders.OrderItems),
		mx: sync.RWMutex{},
	}
}

func (c *Cache) AddCache(key int, value orders.OrderItems) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *Cache) GetCache(key int) (orders.OrderItems, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func CreateCache(db *sqlx.DB) (*Cache, error) {

	cache := NewCache()

	query := fmt.Sprintf(`SELECT * FROM %s ord INNER JOIN %s del ON ord.id = del.order_id INNER JOIN %s pay ON pay.order_id = ord.id`, OrdersTable, DeliveryTable, PaymentTable)
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		order := orders.OrderInfo{}

		err = rows.StructScan(&order)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf(`SELECT * FROM %s WHERE order_id = $1`, ItemsTable)
		rowsI, err := db.Queryx(query, strconv.Itoa(order.Id))
		if err != nil {
			return nil, err
		}
		var items []orders.Item
		for rowsI.Next() {
			item := orders.Item{}
			err := rowsI.StructScan(&item)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}

		cache.AddCache(order.Id, orders.OrderItems{
			OrderInfo: order,
			Items:     items,
		})
	}

	return cache, nil
}
