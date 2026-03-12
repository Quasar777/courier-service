package delivery

import (
	"context"

	uc "github.com/Quasar777/courier-service/internal/usecase/delivery"
)

type DeliveryUseCase interface {
	AssignCourier(ctx context.Context, req uc.AssignInput) (*uc.AssignOutput, error)
	UnassignCourier(ctx context.Context, req uc.UnassignInput) (*uc.UnassignOutput, error) 
}