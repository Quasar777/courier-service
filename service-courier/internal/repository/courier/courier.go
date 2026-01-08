package courier

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Quasar777/courier-service/internal/handler/dto"
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
		return nil, model.ErrInternal
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
			return nil, model.ErrInternal
		}
		couriers = append(couriers, courier)
	}

	if err = rows.Err(); err != nil {
		return nil, model.ErrInternal
	}

	if couriers == nil {
		couriers = []model.Courier{}
	}

	return couriers, nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *dto.CreateCourierRequest) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, courier.Name, courier.Lastname, courier.Phone, courier.Status, courier.TransportType).Scan(&id)

	if err != nil {
		fmt.Printf("err type=%T err=%q\n", err, err.Error())
		if strings.Contains(err.Error(), "duplicate key value") {
			return 0, model.ErrPhoneConflict
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

func (r *CourierRepository) Update(ctx context.Context, courier *dto.UpdateCourierRequest) error {
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

func (r *CourierRepository) ReleaseCouriers(ctx context.Context) error {
	 _, err := r.pool.Exec(ctx, `
		UPDATE couriers 
		SET status = 'available'
		WHERE id IN (
			SELECT courier_id 
			FROM delivery
			WHERE deadline < NOW()
		) AND status = 'busy'
	`)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func (r *CourierRepository) GetCourierIdWithFewestOrders(ctx context.Context) (int, error) {
	var id int

	err := r.pool.QueryRow(ctx,`
		SELECT c.id
		FROM couriers c
		LEFT JOIN delivery d ON d.courier_id = c.id
		WHERE c.status = 'available'
		GROUP BY c.id
		ORDER BY COUNT(d.id) ASC
		LIMIT 1;
		`).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrNoAvailableCouriers
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}