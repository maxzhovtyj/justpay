package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"justpay/internal/domain/order"
	"time"
)

type OrderStorage struct {
	conn *pgx.Conn
}

func NewOrderStorage(db *pgx.Conn) *OrderStorage {
	return &OrderStorage{conn: db}
}

func (s *OrderStorage) GetEvent(eventID uuid.UUID) (e order.Event, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	q := fmt.Sprintf("SELECT id, order_id, user_id, status, is_final, created_at, updated_at FROM events WHERE id = $1")

	row := s.conn.QueryRow(ctx, q, eventID.String())
	if err != nil {
		return e, err
	}

	if err = row.Scan(
		&e.ID,
		&e.OrderID,
		&e.UserID,
		&e.Status,
		&e.IsFinal,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return e, err
	}

	return
}

func (s *OrderStorage) GetEvents(orderID uuid.UUID) (events []order.Event, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	q := fmt.Sprintf("SELECT id, order_id, user_id, status, is_final, created_at, updated_at FROM events WHERE order_id = $1")

	rows, err := s.conn.Query(ctx, q, orderID.String())
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var e order.Event

		if err = rows.Scan(
			&e.ID,
			&e.OrderID,
			&e.UserID,
			&e.Status,
			&e.IsFinal,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return
}

func (s *OrderStorage) CreateEvent(event order.Event) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	q := fmt.Sprintf("INSERT INTO events (id, order_id, user_id, status, is_final, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")

	_, err := s.conn.Exec(ctx, q, event.ID, event.OrderID, event.UserID, event.Status, event.IsFinal, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderStorage) Get(id uuid.UUID) (order order.Order, err error) {
	q := fmt.Sprintf("SELECT id, user_id, status, created_at, updated_at FROM orders WHERE id = $1 LIMIT 1")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	row := s.conn.QueryRow(ctx, q, id.String())

	err = row.Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	return
}

func (s *OrderStorage) GetOrders() (orders []order.Order, err error) {
	q := fmt.Sprintf("SELECT id, user_id, status, created_at, updated_at FROM orders WHERE id = $1 LIMIT 1")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	rows, err := s.conn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o order.Order

		err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)

		orders = append(orders, o)
	}

	return
}
