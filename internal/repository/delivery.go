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

func (r *DeliveryRepository) AssignCourierWithUpdate(ctx context.Context, courierId int, orderId string, deadline time.Time) (*model.Delivery, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE couriers
		SET status = 'paused'
		WHERE id = $1
	`, courierId)
	
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

    d := &model.Delivery{
        CourierId:  courierId,
        OrderId:    orderId,
        AssignedAt: time.Now(),
        Deadline:   deadline,
    }

	err = tx.QueryRow(ctx, `
        INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, d.CourierId, d.OrderId, d.AssignedAt, d.Deadline).Scan(&d.Id)
    if err != nil {
        return nil, fmt.Errorf("insert delivery: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("commit tx: %w", err)
    }

    return d, nil
}


func (r *DeliveryRepository) UnAssign(ctx context.Context, orderId *model.UnAssignDeliveryRequest) (*model.UnAssignedDeliveryResponse, error) {
	return &model.UnAssignedDeliveryResponse{}, nil
}