package courier

import (
	"context"
	"unicode"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
)

type CourierUseCase struct {
	repository CourierRepository
}

func NewCourierUseCase(r CourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: r}
}

func (u *CourierUseCase) GetCourier(ctx context.Context, id int) (*model.Courier, error) {
	if id <= 0 {
		return nil, model.ErrInvalidId
	}

	courier, err := u.repository.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func (u *CourierUseCase) GetCouriers(ctx context.Context) ([]model.Courier, error) {
	return u.repository.GetAll(ctx)
}

func (u *CourierUseCase) CreateCourier(ctx context.Context, req dto.CreateCourierRequest) (int, error) {
	if req.Name == "" || req.Lastname == "" || req.Phone == "" {
        return 0, model.ErrMissingRequiredFields
    }

	if !IsNameValid(req.Name) {
		return 0, model.ErrInvalidCourierName
	}

	if !IsNameValid(req.Lastname) {
		return 0, model.ErrInvalidCourierLastname
	}
	
	if !isPhoneValid(req.Phone) {
		return 0, model.ErrInvalidPhone
	}

	if !isStatusValid(req.Status) {
		return 0, model.ErrInvalidStatus
	}
	if req.Status == "" {
		req.Status = "available"
	}

	if !isTransportyTypeValid(req.TransportType) {
		return 0, model.ErrInvalidCourierTransportType
	}
	if req.TransportType == "" {
		req.TransportType = "on_foot"
	}
	
	return u.repository.Create(ctx, &req)
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, req dto.UpdateCourierRequest) error {
	if req.Id <= 0 {
		return model.ErrMissingRequiredFields
	}

	existingCourier, err := u.repository.GetOneById(ctx, req.Id)
	if err != nil {
		return err
	}

	if req.Name == "" {
		req.Name = existingCourier.Name
	} else {
		if !IsNameValid(req.Name) {
			return model.ErrInvalidCourierName
		}
	}

	if req.Lastname == "" {
		req.Lastname = existingCourier.Lastname
	} else {
		if !IsNameValid(req.Lastname) {
			return model.ErrInvalidCourierLastname
		}
	}

	if req.Phone == "" {
		req.Phone = existingCourier.Phone
	} else {
		if !isPhoneValid(req.Phone) {
			return model.ErrInvalidPhone
		}
	}

	if req.Status == "" {
		req.Status = string(existingCourier.Status)
	} else {
		if !isStatusValid(req.Status) {
			return model.ErrInvalidStatus
		}
	}

	if req.TransportType == "" {
		req.TransportType = string(existingCourier.TransportType)
	} else {
		if !isTransportyTypeValid(req.TransportType) {
			return model.ErrInvalidCourierTransportType
		}
	}
 
	return u.repository.Update(ctx, &req)
}

func (u *CourierUseCase) DeleteCourier(ctx context.Context, id int) error {
	return u.repository.Delete(ctx, id)
}

// Имя должно содержать только буквы латиницы или кириллицы.
// Длина имени допустима от 2 до 20 символов
func IsNameValid(s string) bool {
	length := 0

	for _, r := range s {
		length++

		if !isLatinLetter(r) && !isCyrillicLetter(r) {
			return false
		}
	}

	return length >= 2 && length <= 20
}

func isLatinLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isCyrillicLetter(r rune) bool {
	return (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') || r == 'Ё' || r == 'ё'
}

func isStatusValid(s string) bool {
	switch s {
	case "available", "busy", "paused", "":
		return true
	default:
		return false
	}
}

func isTransportyTypeValid(s string) bool {
	switch s {
	case "car", "scooter", "on_foot", "":
		return true
	default:
		return false
	}
}

func isPhoneValid(s string) bool {
	if len(s) != 12 {
		return false
	}
	if s[0] != '+' || s[1] != '7' || s[2] != '9' {
		return false
	}
	for _, r := range s[3:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}