package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/jackc/pgx/v5"
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
		SET status = 'busy'
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
        return nil, fmt.Errorf("database error: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }

    return d, nil
}


func (r *DeliveryRepository) UnassignWithUpdate(ctx context.Context, orderId string) (*model.Delivery, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer tx.Rollback(ctx)

	d := &model.Delivery{
		OrderId: orderId,
	}

	err = tx.QueryRow(ctx, `
		SELECT courier_id
		FROM delivery 
		WHERE order_id = $1
	`, orderId).Scan(&d.CourierId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoRelationFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE couriers 
		SET status = 'available' 
		WHERE id = $1
	`, d.CourierId)
	
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM delivery
		WHERE order_id = $1
	`, orderId)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }

	return d, nil
}