package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{pool: pool}
}

func (r *DeliveryRepository) Assign(ctx context.Context, req model.AssignedDeliveryResponse) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.CourierId, req.OrderId, time.Now(), req.Deadline).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}


func (r *DeliveryRepository) UnAssign(ctx context.Context, orderId *model.UnAssignDeliveryRequest) (*model.UnAssignedDeliveryResponse, error) {
	return &model.UnAssignedDeliveryResponse{}, nil
}