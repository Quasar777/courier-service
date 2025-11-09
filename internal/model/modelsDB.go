package model

import "time"

type CourierDB struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
