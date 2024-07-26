package service

import (
	"github.com/google/uuid"
	"justpay/internal/domain/order"
	"sync"
)

type Subscriber struct {
	ID uuid.UUID

	events chan order.Event

	eventsBuffer   []*order.Event
	eventsBufferMX sync.RWMutex
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		ID:           uuid.New(),
		events:       make(chan order.Event, 100),
		eventsBuffer: make([]*order.Event, 5),
	}
}

func (s *Subscriber) Notify(e order.Event) {
	s.eventsBufferMX.Lock()
	defer s.eventsBufferMX.Unlock()

	for i, ev := range s.eventsBuffer {
		if i == int(e.Status) {
			s.eventsBuffer[i] = &e
			break
		}

		if ev == nil {
			s.eventsBuffer[e.Status] = &e
			return
		}
	}

	select {
	case s.events <- e:
		// ok
	default:
		// TODO
	}
}

func (s *Subscriber) ReadEvents() <-chan order.Event {
	return s.events
}
