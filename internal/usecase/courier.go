package usecase

import (
	"context"
	"errors"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/Quasar777/courier-service/internal/repository"
)

type CourierUseCase struct {
	repository repository.CourierRepository
}

func NewCourierUseCase(r repository.CourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: r}
}

func (u *CourierUseCase) GetCourier(ctx context.Context, id int) (*model.Courier, error) {
	courierDB, err := u.repository.GetOneById(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCourierNotFound) {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}

	courier := &model.Courier{
		Id: courierDB.Id,
		Name: courierDB.Name,
		Lastname: courierDB.Lastname,
		Status: courierDB.Status,
	}

	return courier, nil
}

func (u *CourierUseCase) GetCouriers(ctx context.Context) ([]model.Courier, error) {
	couriersDB, err := u.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var couriers []model.Courier
	for _, courierDB := range couriersDB {
		courier := model.Courier{
			Id: courierDB.Id,
			Name: courierDB.Name,
			Lastname: courierDB.Lastname,
			Phone: courierDB.Phone,
			Status: courierDB.Status,
		}
		couriers = append(couriers, courier)
	}

	if couriers == nil {
		couriers = []model.Courier{}
	}

	return couriers, nil
}

func (u *CourierUseCase) CreateCourier(ctx context.Context, req model.CreateCourierRequest) (int, error) {
	if req.Name == "" || req.Lastname == "" || req.Phone == "" {
		return 0, ErrMissingRequiredFields
	}

	courierDB := &model.CourierDB{
		Name: req.Name,
		Lastname: req.Lastname,
		Phone: req.Phone,
		Status: req.Status,
	}
	id, err := u.repository.Create(ctx, courierDB)
	if err != nil {
		if errors.Is(err, repository.ErrPhoneConflict) {
			return 0, ErrPhoneConflict
		}
		return 0, err
	}

	return id, nil
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, req model.UpdateCourierRequest) error {
	if req.Name == "" || req.Lastname == "" ||
	   req.Phone == "" || req.Status == "" {
		return ErrMissingRequiredFields
	}

	err := u.repository.Update(ctx, &req)	
	if err != nil {
		if errors.Is(err, repository.ErrPhoneConflict) {
			return ErrPhoneConflict
		}
		if errors.Is(err, repository.ErrCourierNotFound) {
			return  ErrCourierNotFound
		}
		return err
	}

	return nil
}

func (u *CourierUseCase) DeleteCourier(ctx context.Context, id int) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCourierNotFound) {
			return ErrCourierNotFound
		}
		return err
	}

	return nil
}