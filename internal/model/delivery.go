package model

import "time"

type Delivery struct {
	Id         int
	CourierId  int
	OrderId    string
	AssignedAt time.Time
	Deadline   time.Time
}

type CreateDeliveryRequest struct {
	CourierId int
	OrderId   string
	Deadline  time.Time
}

type UpdateDeliveryRequest struct {
	CourierId int
	Deadline  time.Time
}

type AssignDeliveryRequest struct {
	OrderId string
}

type UnassignDeliveryRequest struct {
	OrderId string
}

type AssignedDeliveryResponse struct {
	CourierId     int
	OrderId       string
	TransportType string
	Deadline      time.Time
}

type UnassignedDeliveryResponse struct {
	OrderId   string
	CourierId int
	Status    string
}