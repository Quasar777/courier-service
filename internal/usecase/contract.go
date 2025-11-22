package usecase

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int) (*model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, courier *model.CreateCourierRequest) (int, error)
	Update(ctx context.Context, courier *model.UpdateCourierRequest) error
	Delete(ctx context.Context, id int) error
}