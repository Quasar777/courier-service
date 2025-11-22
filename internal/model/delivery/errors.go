package model

import "errors"

var (
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrMissingRequiredFields = errors.New("missing required fields")
)
