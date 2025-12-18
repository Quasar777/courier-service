package repository_test

import (
	"context"
	"testing"

	"github.com/Quasar777/courier-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type DeliveryRepositoryTestSuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	repo       *repository.DeliveryRepository
	couriersID []int
}

func (s *DeliveryRepositoryTestSuite) SetupSuite() {
	err := godotenv.Load("../../.env.test")
	s.Require().NoError(err)

	pool, err := mustInitPool(context.Background())
	s.Require().NoError(err)

	s.pool = pool
	s.repo = repository.NewDeliveryRepository(s.pool)
}

func (s *DeliveryRepositoryTestSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *DeliveryRepositoryTestSuite) SetupTest() {
	s.couriersID = nil // важно: не копим между тестами
	ctx := context.Background()

	ids := make([]int, 0, 3)

	var id1 int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "Andrew", "Downsky", "+79990000001", "available", "on_foot").Scan(&id1)
	s.Require().NoError(err)
	ids = append(ids, id1)

	var id2 int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "Sergey", "Sergeev", "+79990000002", "available", "scooter").Scan(&id2)
	s.Require().NoError(err)
	ids = append(ids, id2)

	var id3 int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "Ivan", "Ivanov", "+79990000003", "available", "car").Scan(&id3)
	s.Require().NoError(err)
	ids = append(ids, id3)

	s.couriersID = ids
}

func (s *DeliveryRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `TRUNCATE delivery, couriers RESTART IDENTITY CASCADE`)
	s.Require().NoError(err)

	s.couriersID = nil
}

func TestDeliveryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryRepositoryTestSuite))
}
