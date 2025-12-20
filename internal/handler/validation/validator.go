package validation

import (
	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
)

type Validator struct{}

func ValidateRequest[T any](req *T) error {
	switch v := any(req).(type) {
	case *dto.CreateCourierRequest:
		return validateCreate(v)
	case *dto.UpdateCourierRequest:
		return validateUpdate(v)
	default:
		return nil
	}
}

func validateCreate(req *dto.CreateCourierRequest) error {
	if req.Name == "" || req.Lastname == "" || req.Phone == "" {
        return model.ErrMissingRequiredFields
    }
	if err := validateStatus(req.Status); err != nil {
		return err
	}
	if err := validateTransport(req.TransportType); err != nil {
		return err
	}
	if len(req.Phone) < 10 {
		return model.ErrInvalidPhone
	}
	return nil
}

func validateUpdate(req *dto.UpdateCourierRequest) error {
	if req.Id <= 0 {
		return model.ErrMissingRequiredFields
	}
	if req.Name == "" && req.Phone == "" && req.Status == "" && req.TransportType == "" {
		return model.ErrMissingRequiredFields
	}
	if req.Status != "" {
		if err := validateStatus(req.Status); err != nil {
			return err
		}
	}
	if req.TransportType != "" {
		if err := validateTransport(req.TransportType); err != nil {
			return err
		}
	}
	if req.Phone != "" && len(req.Phone) < 10 {
		return model.ErrInvalidPhone
	}
	return nil
}

func validateStatus(status string) error {
	switch status {
	case "available", "busy", "paused", "":
		return nil
	default:
		return model.ErrInvalidStatus
	}
}

func validateTransport(t string) error {
	switch t {
	case "car", "scooter", "on_foot", "":
		return nil
	default:
		return model.ErrInvalidCourierTransportType
	}
}