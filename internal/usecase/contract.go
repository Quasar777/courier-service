package usecase

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int) (*model.CourierDB, error)
	GetAll(ctx context.Context) ([]model.CourierDB, error)
	Create(ctx context.Context, courier *model.CreateCourierRequest) (int, error)
	Update(ctx context.Context, courier *model.UpdateCourierRequest) error
	Delete(ctx context.Context, id int) error
}