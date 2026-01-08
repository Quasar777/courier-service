package model

import "errors"

var (
	ErrInvalidJSON                 = errors.New("invalid JSON")
	ErrMissingRequiredFields       = errors.New("missing required fields")
	ErrInvalidId                   = errors.New("invalid id")
	ErrInvalidCourierName          = errors.New("invalid courier name")
	ErrInvalidCourierLastname      = errors.New("invalid courier lastname")
	ErrInvalidPhone                = errors.New("ivalid phone")
	ErrInvalidStatus               = errors.New("invalid status")
	ErrInvalidCourierTransportType = errors.New("invalid courier transport type")

	ErrCourierNotFound = errors.New("courier not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")

	ErrDeliveryNotFound    = errors.New("delivery not found")
	ErrNoAvailableCouriers = errors.New("no available couriers")
	ErrNoRelationFound     = errors.New("realtion between courier and delivery not found")
	ErrInvalidOrderID      = errors.New("invalid order id")

	ErrInternal = errors.New("internal server error")
)
