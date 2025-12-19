package model

import "errors"

var (
	ErrMissingRequiredFields       = errors.New("missing required fields")
	ErrInvalidId                   = errors.New("invalid id")
	ErrInvalidPhone                = errors.New("ivalid phone")
	ErrInvalidStatus               = errors.New("invalid status")
	ErrInvalidCourierTransportType = errors.New("invalid courier transport type")

	ErrCourierNotFound = errors.New("courier not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")

	ErrDeliveryNotFound    = errors.New("delivery not found")
	ErrNoAvailableCouriers = errors.New("no available couriers")
	ErrNoRelationFound     = errors.New("realtion between courier and delivery not found")
)
