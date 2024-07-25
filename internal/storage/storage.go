package storage

type Storage struct {
	Order
}

type Order interface {
}

func New() *Storage {
	return &Storage{}
}
