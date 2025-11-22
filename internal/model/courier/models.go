package model

import "time"

type Courier struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Lastname      string    `json:"lastname"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	TransportType string    `json:"transportType"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type CreateCourierRequest struct {
	Name          string `json:"name"`
	Lastname      string `json:"lastname"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transportType"`
}

type UpdateCourierRequest struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Lastname      string `json:"lastname"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transportType"`
}
