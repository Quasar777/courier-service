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
	return u.repository.Create(ctx, &req)
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, req model.UpdateCourierRequest) error {
	return u.repository.Update(ctx, &req)
}

func (u *CourierUseCase) DeleteCourier(ctx context.Context, id int) error {
	return u.repository.Delete(ctx, id)
}