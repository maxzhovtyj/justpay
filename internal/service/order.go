package service

import (
	"github.com/google/uuid"
	"justpay/internal/domain/order"
	"justpay/internal/storage"
	"sync"
)

type OrderService struct {
	storage *storage.Storage

	orderEvents   map[string]*Events
	orderEventsMX sync.RWMutex
}

func NewOrderService(db *storage.Storage) *OrderService {
	return &OrderService{
		storage:     db,
		orderEvents: make(map[string]*Events),
	}
}

func (s *OrderService) GetOrders() ([]order.Order, error) {
	return s.storage.Order.GetOrders()
}

func (s *OrderService) NewEvent(event order.Event) error {
	e, err := s.storage.Order.GetEvent(event.ID)
	if e.ID.String() != "" && err == nil {
		return order.ErrEventAlreadyExists
	}

	o, err := s.storage.Order.Get(event.OrderID)
	if err != nil {
		return err
	}

	if o.IsFinal {
		return order.ErrFinalStatusReceived
	}

	err = s.storage.Order.CreateEvent(event)
	if err != nil {
		return err
	}

	s.NotifySubscribers(event)

	return nil
}

func (s *OrderService) Subscribe(oid uuid.UUID) (*Subscriber, error) {
	_, err := s.storage.Order.Get(oid)
	if err != nil {
		return nil, err
	}

	orderID := oid.String()

	sub := NewSubscriber()

	events, ok := s.getOrderEvents(orderID)
	if !ok {
		events = s.newOrderEvents(orderID)
		s.addOrderEvents(orderID, events)
	}

	events.addSubscriber(sub.ID.String(), sub)

	prevEvents, err := s.storage.Order.GetEvents(oid)
	if err != nil {
		return nil, err
	}

	for _, e := range prevEvents {
		sub.Notify(e)
	}

	return sub, nil
}

func (s *OrderService) Unsubscribe(orderID, userID uuid.UUID) {
	events, ok := s.getOrderEvents(orderID.String())
	if !ok {
		return
	}

	events.removeSubscriber(userID.String())
}

func (s *OrderService) NotifySubscribers(event order.Event) {
	events, ok := s.getOrderEvents(event.OrderID.String())
	if !ok {
		// means no subscribers yet
		return
	}

	events.NotifyAll(event)
}

func (s *OrderService) newOrderEvents(orderID string) *Events {
	s.orderEventsMX.Lock()
	defer s.orderEventsMX.Unlock()

	// check whether the new Events struct was created by another goroutine
	if e, ok := s.orderEvents[orderID]; ok {
		return e
	}

	return NewEvents(orderID)
}

func (s *OrderService) getOrderEvents(orderID string) (*Events, bool) {
	s.orderEventsMX.RLock()
	defer s.orderEventsMX.RUnlock()

	events, ok := s.orderEvents[orderID]
	return events, ok
}

func (s *OrderService) addOrderEvents(orderID string, e *Events) {
	s.orderEventsMX.Lock()
	defer s.orderEventsMX.Unlock()

	s.orderEvents[orderID] = e
}
