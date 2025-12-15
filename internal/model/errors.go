package model

import "errors"

var (
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidId = errors.New("invalid id")

	ErrCourierNotFound = errors.New("courier not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")

	ErrDeliveryNotFound    = errors.New("delivery not found")
	ErrNoAvailableCouriers = errors.New("no available couriers")
	ErrNoRelationFound     = errors.New("realtion between courier and delivery not found")
)
