package order

import "time"

type Order struct {
	OrderID   string
	Status    string
	CreatedAt time.Time
}

type OrderID struct {
	OrderID string
}