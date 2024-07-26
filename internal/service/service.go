package service

import (
	"github.com/google/uuid"
	"justpay/internal/domain/order"
	"justpay/internal/storage"
)

type Service struct {
	Order Order
}

type Order interface {
	NewEvent(event order.Event) error
	Subscribe(oid uuid.UUID) (*Subscriber, error)
	Unsubscribe(orderID, userID uuid.UUID)
	NotifySubscribers(event order.Event)
	GetOrders() ([]order.Order, error)
}

func New(db *storage.Storage) *Service {
	return &Service{
		Order: NewOrderService(db),
	}
}
