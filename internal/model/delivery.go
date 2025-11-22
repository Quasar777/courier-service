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
	CourierId int       `json:"courierId"`
	OrderId   int       `json:"orderId"`
	Deadline  time.Time `json:"deadline"`
}

type UpdateDeliveryRequest struct {
	CourierId int       `json:"courierId"`
	Deadline  time.Time `json:"deadline"`
}

type AssignedDeliveryResponse struct {
	CourierId     int       `json:"courierId"`
	OrderId       int       `json:"orderId"`
	TransportType string    `json:"transportType"`
	Deadline      time.Time `json:"deadline"`
}

type UnAssignedDeliveryResponse struct {
	OrderId   int    `json:"orderId"`
	CourierId int    `json:"courierId"`
	Status    string `json:"status"`
}
