package model

import "time"

const (
	StatusAvailable CourierStatus = "available"
	StatusBusy      CourierStatus = "busy"
	StatusPaused    CourierStatus = "paused"

	TransportTypeOnFoot  CourierTransportType = "on_foot"
	TransportTypeScooter CourierTransportType = "scooter"
	TransportTypeCar     CourierTransportType = "car"
)

var TransportTypeFromDto = map[string]CourierTransportType{
	"on_foot": TransportTypeOnFoot,
	"scooter": TransportTypeScooter,
	"car":     TransportTypeCar,
}

var DtoFromTransportType = map[CourierTransportType]string{
	TransportTypeOnFoot:  "on_foot",
	TransportTypeScooter: "scooter",
	TransportTypeCar:     "car",
}

type Courier struct {
	Id            int
	Name          string
	Lastname      string
	Phone         string
	Status        CourierStatus
	TransportType CourierTransportType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CourierStatus string
type CourierTransportType string