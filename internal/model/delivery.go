package model

import "time"

type Delivery struct {
	Id         int
	CourierId  int
	OrderId    string
	AssignedAt time.Time
	Deadline   time.Time
}