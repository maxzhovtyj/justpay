package service

import (
	"justpay/internal/service/order"
	"justpay/internal/storage"
)

type Service struct {
	Order Order
}

type Order interface {
	Create()
	NewEvent()
	Subscribe(orderID string) (*order.Subscriber, error)
	Unsubscribe()
}

func New(db *storage.Storage) *Service {
	return &Service{
		Order: order.New(db),
	}
}
