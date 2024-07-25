package order

import "justpay/internal/domain/order"

type Subscriber struct {
	Events chan order.Event
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Events: make(chan order.Event),
	}
}
