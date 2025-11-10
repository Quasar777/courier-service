package repository

import "errors"

var (
	ErrCourierNotFound = errors.New("courier with a such id is not found")
	ErrPhoneConflict   = errors.New("courier with this phone number is already exists")
)
