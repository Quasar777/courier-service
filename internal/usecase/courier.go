package usecase

// UseCase / Service / Manager - business logic layer 

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierUseCase struct {
	repository CourierRepository
}

func NewCourierUseCase(r CourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: r}
}

func (u *CourierUseCase) GetCourier(ctx context.Context, id int) (*model.Courier, error) {
	return u.repository.GetOneById(ctx, id)
}

func (u *CourierUseCase) GetCouriers(ctx context.Context) ([]model.Courier, error) {
	return u.repository.GetAll(ctx)
}

func (u *CourierUseCase) CreateCourier(ctx context.Context, req model.CreateCourierRequest) (int, error) {
	if req.Name == "" || req.Lastname == "" || req.Phone == "" {
        return 0, model.ErrMissingRequiredFields
    }

	if req.Status == "" {
		req.Status = "available"
	}

	if req.TransportType == "" {
		req.TransportType = "on_foot"
	}
	
	return u.repository.Create(ctx, &req)
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, req model.UpdateCourierRequest) error {
	if req.Name == "" && req.Lastname == "" && req.Phone == "" && req.Status == "" && req.TransportType == "" {
        return model.ErrMissingRequiredFields
    }

	existingCourier, err := u.repository.GetOneById(ctx, req.Id)
	if err != nil {
		return err
	}

	if req.Name == "" {
		req.Name = existingCourier.Name
	}
	if req.Lastname == "" {
		req.Lastname = existingCourier.Lastname
	}
	if req.Phone == "" {
		req.Phone = existingCourier.Phone
	}
	if req.Status == "" {
		req.Status = existingCourier.Status
	}
	if req.TransportType == "" {
		req.TransportType = existingCourier.TransportType
	}
 
	return u.repository.Update(ctx, &req)
}

func (u *CourierUseCase) DeleteCourier(ctx context.Context, id int) error {
	return u.repository.Delete(ctx, id)
}