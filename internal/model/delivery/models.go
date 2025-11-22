package model

import "time"

type Delivery struct {
	Id         int       `json:"id"`
	CourierId  int       `json:"courierId"`
	OrderId    int       `json:"orderId"`
	AssignedAt time.Time `json:"assignedAt"`
	Deadline   time.Time `json:"deadline"`
}

type CreateDeliveryRequest struct {
	CourierId  int       `json:"courierId"`
	OrderId    int       `json:"orderId"`
	Deadline   time.Time `json:"deadline"`
}

type UpdateDeliveryRequest struct {
	CourierId  int       `json:"courierId"`
	Deadline   time.Time `json:"deadline"`
}