package courier

import (
	"context"

	uc "github.com/Quasar777/courier-service/internal/usecase/courier"
	"github.com/Quasar777/courier-service/internal/model"
)

type CourierUseCase interface {
	GetCourier(ctx context.Context, id int) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, req uc.CreateCourierInput) (int, error)
	UpdateCourier(ctx context.Context, req uc.UpdateCourierInput) error
	DeleteCourier(ctx context.Context, id int) error
}