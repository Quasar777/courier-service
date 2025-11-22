package model

import "errors"

var (
	ErrMissingRequiredFields = errors.New("missing required fields")

	ErrCourierNotFound = errors.New("courier not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")
	
	ErrDeliveryNotFound = errors.New("delivery not found")
)
