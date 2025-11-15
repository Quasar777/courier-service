package model

import "errors"

var (
	ErrCourierNotFound = errors.New("courier not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")
	ErrMissingRequiredFields = errors.New("missing required fields")
)
