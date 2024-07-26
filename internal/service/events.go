package service

import (
	"justpay/internal/domain/order"
	"sync"
)

type Events struct {
	OrderID string

	subscribers   map[string]*Subscriber
	subscribersMX sync.RWMutex
}

func NewEvents(orderID string) *Events {
	return &Events{
		OrderID:     orderID,
		subscribers: make(map[string]*Subscriber),
	}
}

func (e *Events) NotifyAll(event order.Event) {
	e.subscribersMX.RLock()
	defer e.subscribersMX.RUnlock()

	for _, sub := range e.subscribers {
		sub.Notify(event)
	}
}

func (e *Events) addSubscriber(userID string, sub *Subscriber) {
	e.subscribersMX.Lock()
	defer e.subscribersMX.Unlock()

	e.subscribers[userID] = sub
}

func (e *Events) removeSubscriber(userID string) {
	e.subscribersMX.Lock()
	defer e.subscribersMX.Unlock()

	delete(e.subscribers, userID)
}
