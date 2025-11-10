package usecase

import "errors"

var (
	ErrCourierNotFound       = errors.New("courier not found")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrPhoneConflict = errors.New("courier with this phone number is already exists")
)
