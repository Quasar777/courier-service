package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
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

func (s *DeliveryRepositoryTestSuite) TestAssign_Success() {
	ctx := context.Background()

	courierID := s.couriersID[0]
	orderID := "order-1"
	deadline := time.Now().Add(30 * time.Minute)

	d, err := s.repo.AssignCourierWithUpdate(ctx, courierID, orderID, deadline)

	s.Require().NoError(err)
	s.Require().NotNil(d)
	s.Greater(d.Id, 0)
	s.Equal(courierID, d.CourierId)
	s.Equal(orderID, d.OrderId)

	var (
		dbID       int
		dbCourier  int
		dbOrderID  string
		dbAssigned time.Time
		dbDeadline time.Time
	)
	err = s.pool.QueryRow(ctx, `
		SELECT id, courier_id, order_id, assigned_at, deadline
		FROM delivery
		WHERE id = $1
	`, d.Id).Scan(&dbID, &dbCourier, &dbOrderID, &dbAssigned, &dbDeadline)
	s.Require().NoError(err)

	s.Equal(d.Id, dbID)
	s.Equal(courierID, dbCourier)
	s.Equal(orderID, dbOrderID)
	s.Equal(d.Deadline.Format("2006-01-02 15:04:05"), dbDeadline.Format("2006-01-02 15:04:05"))

	var status string
	err = s.pool.QueryRow(ctx, `
		SELECT status
		FROM couriers
		WHERE id = $1
	`, courierID).Scan(&status)
	s.Require().NoError(err)
	s.Equal("busy", status)
}

func (s *DeliveryRepositoryTestSuite) TestAssign_CourierNotFound() {
	ctx := context.Background()

	nonExistingCourier := 999999999
	orderID := "order-no-courier"
	deadline := time.Now().Add(30 * time.Minute)

	d, err := s.repo.AssignCourierWithUpdate(ctx, nonExistingCourier, orderID, deadline)

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrCourierNotFound)
	s.Nil(d)
}

func (s *DeliveryRepositoryTestSuite) TestUnassign_Success() {
	ctx := context.Background()

	courierID := s.couriersID[0]
	orderID := "order-unassign-1"
	deadline := time.Now().Add(30 * time.Minute)

	_, err := s.repo.AssignCourierWithUpdate(ctx, courierID, orderID, deadline)
	s.Require().NoError(err)

	var before string
	err = s.pool.QueryRow(ctx, `SELECT status FROM couriers WHERE id = $1`, courierID).Scan(&before)
	s.Require().NoError(err)
	s.Equal("busy", before)

	d, err := s.repo.UnassignWithUpdate(ctx, orderID)

	s.Require().NoError(err)
	s.Require().NotNil(d)
	s.Equal(orderID, d.OrderId)
	s.Equal(courierID, d.CourierId)

	var cnt int
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM delivery WHERE order_id = $1`, orderID).Scan(&cnt)
	s.Require().NoError(err)
	s.Equal(0, cnt)

	var after string
	err = s.pool.QueryRow(ctx, `SELECT status FROM couriers WHERE id = $1`, courierID).Scan(&after)
	s.Require().NoError(err)
	s.Equal("available", after)
}

func (s *DeliveryRepositoryTestSuite) TestUnassign_NoRelationFound() {
	ctx := context.Background()

	d, err := s.repo.UnassignWithUpdate(ctx, "order-does-not-exist")

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrNoRelationFound)
	s.Nil(d)
}

func (s *DeliveryRepositoryTestSuite) TestReleaseCouriers_ReleasesOnlyExpiredBusy() {
	ctx := context.Background()

	// Подготовка:
	// courier1 busy + deadline в прошлом => должен стать available
	// courier2 busy + deadline в будущем => должен остаться busy
	// courier3 available + deadline в прошлом => должен остаться available

	c1 := s.couriersID[0]
	c2 := s.couriersID[1]
	c3 := s.couriersID[2]

	_, err := s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy' WHERE id = $1`, c1)
	s.Require().NoError(err)
	_, err = s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy' WHERE id = $1`, c2)
	s.Require().NoError(err)

	past := time.Now().Add(-2 * time.Hour)
	future := time.Now().Add(2 * time.Hour)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), $3)
	`, c1, "order-expired", past)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), $3)
	`, c2, "order-active", future)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), $3)
	`, c3, "order-expired-available", past)
	s.Require().NoError(err)

	
	err = s.repo.ReleaseCouriers(ctx)
	s.Require().NoError(err)

	
	status := func(id int) string {
		var st string
		e := s.pool.QueryRow(ctx, `SELECT status FROM couriers WHERE id = $1`, id).Scan(&st)
		s.Require().NoError(e)
		return st
	}

	s.Equal("available", status(c1))
	s.Equal("busy", status(c2))
	s.Equal("available", status(c3))
}

func (s *DeliveryRepositoryTestSuite) TestReleaseCouriers_NoExpired_NoChanges() {
	ctx := context.Background()

	c1 := s.couriersID[0]
	c2 := s.couriersID[1]

	_, err := s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy' WHERE id = $1`, c1)
	s.Require().NoError(err)
	_, err = s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy' WHERE id = $1`, c2)
	s.Require().NoError(err)

	future := time.Now().Add(2 * time.Hour)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), $3)
	`, c1, "order-f1", future)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), $3)
	`, c2, "order-f2", future)
	s.Require().NoError(err)

	err = s.repo.ReleaseCouriers(ctx)
	s.Require().NoError(err)

	var st1, st2 string
	s.Require().NoError(s.pool.QueryRow(ctx, `SELECT status FROM couriers WHERE id = $1`, c1).Scan(&st1))
	s.Require().NoError(s.pool.QueryRow(ctx, `SELECT status FROM couriers WHERE id = $1`, c2).Scan(&st2))
	s.Equal("busy", st1)
	s.Equal("busy", st2)
}

func (s *DeliveryRepositoryTestSuite) TestGetCourierIdWithFewestOrders_Success() {
	ctx := context.Background()

	c1 := s.couriersID[0]
	c2 := s.couriersID[1]
	c3 := s.couriersID[2]

	_, err := s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), NOW() + interval '1 hour'),
		       ($1, $3, NOW(), NOW() + interval '1 hour')
	`, c1, "order-c1-1", "order-c1-2")
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), NOW() + interval '1 hour')
	`, c2, "order-c2-1")
	s.Require().NoError(err)


	id, err := s.repo.GetCourierIdWithFewestOrders(ctx)


	s.Require().NoError(err)
	s.Equal(c3, id)
}

func (s *DeliveryRepositoryTestSuite) TestGetCourierIdWithFewestOrders_IgnoresBusy() {
	ctx := context.Background()

	c1 := s.couriersID[0]
	c2 := s.couriersID[1]
	c3 := s.couriersID[2]

	_, err := s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy' WHERE id = $1`, c3)
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), NOW() + interval '1 hour'),
		       ($1, $3, NOW(), NOW() + interval '1 hour')
	`, c1, "order1", "order2")
	s.Require().NoError(err)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW(), NOW() + interval '1 hour')
	`, c2, "order3")
	s.Require().NoError(err)

	id, err := s.repo.GetCourierIdWithFewestOrders(ctx)

	s.Require().NoError(err)
	s.Equal(c2, id)
}

func (s *DeliveryRepositoryTestSuite) TestGetCourierIdWithFewestOrders_NoAvailable() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `UPDATE couriers SET status = 'busy'`)
	s.Require().NoError(err)

	id, err := s.repo.GetCourierIdWithFewestOrders(ctx)

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrNoAvailableCouriers)
	s.Equal(0, id)
}


func (s *DeliveryRepositoryTestSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *DeliveryRepositoryTestSuite) SetupTest() {
	s.couriersID = nil
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
