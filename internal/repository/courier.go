package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) *CourierRepository {
	return &CourierRepository{pool: pool}
}

func (r *CourierRepository) GetOneById(ctx context.Context, id int) (*model.Courier, error) {
	var courier model.Courier
	err := r.pool.QueryRow(ctx,`
		SELECT id, name, lastname, phone, status, transport_type
		FROM couriers 
		WHERE id = $1
		`, id).Scan(
		&courier.Id, 
		&courier.Name, 
		&courier.Lastname, 
		&courier.Phone, 
		&courier.Status,
		&courier.TransportType,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrCourierNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &courier, nil
}

func (r *CourierRepository) GetAll(ctx context.Context) ([]model.Courier, error) {
	rows, err := r.pool.Query(ctx,`
		SELECT id, name, lastname, phone, status, transport_type
		FROM couriers
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var couriers []model.Courier
	for rows.Next() {
		var courier model.Courier
		err := rows.Scan(
			&courier.Id, 
			&courier.Name, 
			&courier.Lastname, 
			&courier.Phone, 
			&courier.Status,
			&courier.TransportType,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading data: %w", err)
		}
		couriers = append(couriers, courier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if couriers == nil {
		couriers = []model.Courier{}
	}

	return couriers, nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *model.CreateCourierRequest) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, courier.Name, courier.Lastname, courier.Phone, courier.Status, courier.TransportType).Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			fmt.Println(err)
			return 0, model.ErrPhoneConflict
		}
		fmt.Println(err)
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

func (r *CourierRepository) Update(ctx context.Context, courier *model.UpdateCourierRequest) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE couriers
		SET name = $1, lastname = $2, phone = $3, status = $4, transport_type = $5
		WHERE id = $6
	`,  courier.Name, 
	    courier.Lastname, 
	    courier.Phone, 
	    courier.Status, 
		courier.TransportType,
		courier.Id,)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return model.ErrPhoneConflict
		}
		return fmt.Errorf("database error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrCourierNotFound
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE couriers
		SET updated_at = NOW()
		WHERE id = $1
	`, courier.Id)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func (r *CourierRepository) Delete(ctx context.Context, id int) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM couriers
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrCourierNotFound
	}

	return nil
}