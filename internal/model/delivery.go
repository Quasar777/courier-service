package model

import "time"

type Delivery struct {
	Id         int       `json:"id"`
	CourierId  int       `json:"courierId"`
	OrderId    string    `json:"orderId"`
	AssignedAt time.Time `json:"assignedAt"`
	Deadline   time.Time `json:"deadline"`
}

type CreateDeliveryRequest struct {
	CourierId int       `json:"courierId"`
	OrderId   string    `json:"orderId"`
	Deadline  time.Time `json:"deadline"`
}

type UpdateDeliveryRequest struct {
	CourierId int       `json:"courierId"`
	Deadline  time.Time `json:"deadline"`
}

type AssignDeliveryRequest struct {
	OrderId string `json:"orderId"`
}

type UnassignDeliveryRequest struct {
	OrderId string `json:"orderId"`
}

type AssignedDeliveryResponse struct {
	CourierId     int       `json:"courierId"`
	OrderId       string    `json:"orderId"`
	TransportType string    `json:"transportType"`
	Deadline      time.Time `json:"deadline"`
}

type UnassignedDeliveryResponse struct {
	OrderId   string `json:"orderId"`
	CourierId int    `json:"courierId"`
	Status    string `json:"status"`
}
