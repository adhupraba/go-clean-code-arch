package order

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adhupraba/ecom/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res := s.db.QueryRow(
		context.Background(),
		`INSERT INTO orders ("userId", total, status, address) VALUES ($1, $2, $3, $4) RETURNING id`,
		order.UserID, order.Total, order.Status, order.Address,
	)

	var id int

	err := res.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec(
		context.Background(),
		`INSERT INTO "orderItems" ("orderId", "productId", quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price,
	)

	return err
}
