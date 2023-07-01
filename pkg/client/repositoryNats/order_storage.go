package repositoryNats

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"orders"
	"orders/pkg/server/repository"
)

type Storage struct {
	c  *repository.Cache
	db *sqlx.DB
}

func New(cache *repository.Cache, db *sqlx.DB) *Storage {
	return &Storage{
		c:  cache,
		db: db,
	}
}

func (r *Storage) CreateOrderItems(order *orders.OrderItems) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createOrderQuery := fmt.Sprintf(`INSERT INTO %s (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, repository.OrdersTable)
	row := tx.QueryRow(createOrderQuery, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)
	if err = row.Scan(&order.Order.Id); err != nil {
		tx.Rollback()
		return err
	}

	createDeliveryQuery := fmt.Sprintf(`INSERT INTO %s (order_id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, repository.DeliveryTable)
	_, err = tx.Exec(createDeliveryQuery, order.Order.Id, order.Name, order.Phone, order.Zip, order.City, order.Address, order.Region, order.Region)
	if err != nil {
		tx.Rollback()
		return err
	}

	createPaymentQuery := fmt.Sprintf(`INSERT INTO %s (order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11)`, repository.PaymentTable)
	_, err = tx.Exec(createPaymentQuery, order.Order.Id, order.Transaction, order.RequestId, order.Currency, order.Provider, order.Amount, order.PaymentDt, order.Bank, order.DeliveryCost, order.GoodsTotal, order.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		createItemQuery := fmt.Sprintf(`INSERT INTO %s (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, repository.ItemsTable)
		_, err = tx.Exec(createItemQuery, order.Order.Id, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	r.c.AddCache(order.Order.Id, *order)

	return nil
}
