package model

import "time"

type Courier struct {
	Id            int
	Name          string
	Lastname      string
	Phone         string
	Status        string
	TransportType string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateCourierRequest struct {
	Name          string
	Lastname      string
	Phone         string
	Status        string
	TransportType string
}

type UpdateCourierRequest struct {
	Id            int
	Name          string
	Lastname      string
	Phone         string
	Status        string
	TransportType string
}
