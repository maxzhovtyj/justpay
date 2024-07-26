package storage

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"justpay/internal/domain/order"
)

type Storage struct {
	Order
}

type Order interface {
	GetEvent(eventID uuid.UUID) (e order.Event, err error)
	GetEvents(orderID uuid.UUID) (events []order.Event, err error)
	CreateEvent(event order.Event) error
	Get(id uuid.UUID) (order order.Order, err error)
	GetOrders() (orders []order.Order, err error)
}

func New(conn *pgx.Conn) *Storage {
	return &Storage{
		Order: NewOrderStorage(conn),
	}
}
