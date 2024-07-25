package order

import (
	"justpay/internal/domain/order"
	"justpay/internal/storage"
	"sync"
)

type Service struct {
	storage *storage.Storage

	orderEvents   map[string]*Events
	orderEventsMX sync.RWMutex
}

type Events struct {
	OrderID string

	subscribers   map[string]*Subscriber
	subscribersMX sync.RWMutex
}

func (s *Service) Create(order order.Order) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) NewEvent() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Subscribe(orderID string) (*Subscriber, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Unsubscribe() {
	//TODO implement me
	panic("implement me")
}

func New(db *storage.Storage) *Service {
	return &Service{
		storage:     db,
		subscribers: make(map[string]*Subscriber),
	}
}

func (s *Service) NotifySubscribers(event order.Event) {
	event
}

func (s *Service) addSubscriber(orderID string, sub *Subscriber) {
	s.subscribersMX.Lock()
	defer s.subscribersMX.Unlock()

	s.subscribers[orderID] = append(s.subscribers[orderID], sub)
}
