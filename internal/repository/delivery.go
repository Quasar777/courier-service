package repository

import (
	"context"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{pool: pool}
}

func (r *DeliveryRepository) Assign(ctx context.Context, orderId model.AssignDeliveryRequest) (model.AssignedDeliveryResponse, error) {
	return model.AssignedDeliveryResponse{}, nil
}


func (r *DeliveryRepository) UnAssign(ctx context.Context, orderId model.UnAssignDeliveryRequest) (model.UnAssignedDeliveryResponse, error) {
	return model.UnAssignedDeliveryResponse{}, nil
}