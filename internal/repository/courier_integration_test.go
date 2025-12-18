package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
	"github.com/Quasar777/courier-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type CourierRepositoryTestSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	repo *repository.CourierRepository
	couriersId []int
}

func (s *CourierRepositoryTestSuite) SetupSuite() {
	err := godotenv.Load("../../.env.test")
	s.Require().NoError(err)

	pool, err := mustInitPool(context.Background())
	s.Require().NoError(err)

	s.pool = pool

	s.repo = repository.NewCourierRepository(s.pool)
}

func (s *CourierRepositoryTestSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *CourierRepositoryTestSuite) SetupTest() {
	ctx := context.Background()

	var id1 int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, "Andrew",
		"Downsky",
		"+79990000001",
		"available",
		"on_foot",
	).Scan(&id1)
	s.Require().NoError(err)

	var id2 int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, "Sergey",
		"Sergeev",
		"+79990000002",
		"available",
		"scooter",
	).Scan(&id2)
	s.Require().NoError(err)

	var id3 int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, "Ivan",
		"Ivanov",
		"+79990000003",
		"available",
		"car",
	).Scan(&id3)
	s.Require().NoError(err)

	
	s.couriersId = append(s.couriersId, id1, id2, id3)
}

func (s *CourierRepositoryTestSuite) TestGetById() {
	ctx := context.Background()

	courier, err := s.repo.GetOneById(ctx, s.couriersId[0])

	s.Require().NoError(err)
	s.Equal(s.couriersId[0], courier.Id)
	s.Equal("Andrew", courier.Name)
	s.Equal("Downsky", courier.Lastname)
	s.Equal("+79990000001", courier.Phone)
	s.Equal("available", courier.Status)
	s.Equal("on_foot", courier.TransportType)
}

func (s *CourierRepositoryTestSuite) TestGetById_NotFound() {
	ctx := context.Background()

	notExistingID := s.couriersId[0] + 10_000_000

	courier, err := s.repo.GetOneById(ctx, notExistingID)

	s.Require().Error(err)
	s.Nil(courier)
	s.ErrorIs(err, model.ErrCourierNotFound)
}

func (s *CourierRepositoryTestSuite) TestGetById_InvalidId() {
	ctx := context.Background()

	courier, err := s.repo.GetOneById(ctx, 0)

	s.Require().Error(err)
	s.Nil(courier)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Empty() {
	ctx := context.Background()

	for _, id := range s.couriersId {
		_, err := s.pool.Exec(ctx, `
			DELETE FROM couriers WHERE id = $1	
		`, id)
		s.Require().NoError(err)
	}

	couriers, err := s.repo.GetAll(ctx)

	s.Require().NoError(err)
	s.NotNil(couriers)
	s.Len(couriers, 0)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Multiple() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
	s.Require().NoError(err)

	var id1, id2 int

	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "A", "One", "111", "available", "on_foot").Scan(&id1)
	s.Require().NoError(err)

	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "B", "Two", "222", "busy", "bike").Scan(&id2)
	s.Require().NoError(err)

	couriers, err := s.repo.GetAll(ctx)

	s.Require().NoError(err)
	s.Len(couriers, 2)

	s.Equal(id1, couriers[0].Id)
	s.Equal(id2, couriers[1].Id)

	s.Equal("A", couriers[0].Name)
	s.Equal("B", couriers[1].Name)

}

func (s *CourierRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
	s.Require().NoError(err)

	courier := &model.CreateCourierRequest{
		Name:          "Andrew",
		Lastname:      "Downsky",
		Phone:         "+79990009988",
		Status:        "available",
		TransportType: "on_foot",
	}

	id, err := s.repo.Create(ctx, courier)

	s.Require().NoError(err)
	s.Greater(id, 0)

	var dbCourier model.Courier
	err = s.pool.QueryRow(ctx, `
		SELECT id, name, lastname, phone, status, transport_type
		FROM couriers
		WHERE id = $1
	`, id).Scan(&dbCourier.Id, &dbCourier.Name, &dbCourier.Lastname, &dbCourier.Phone, &dbCourier.Status, &dbCourier.TransportType)

	s.Require().NoError(err)
	s.Equal(courier.Name, dbCourier.Name)
	s.Equal(courier.Lastname, dbCourier.Lastname)
	s.Equal(courier.Phone, dbCourier.Phone)
	s.Equal(courier.Status, dbCourier.Status)
	s.Equal(courier.TransportType, dbCourier.TransportType)
}

func (s *CourierRepositoryTestSuite) TestCreate_PhoneConflict() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
	s.Require().NoError(err)

	courier1 := &model.CreateCourierRequest{
		Name:          "Andrew",
		Lastname:      "Downsky",
		Phone:         "+79990009988",
		Status:        "available",
		TransportType: "on_foot",
	}

	_, err = s.repo.Create(ctx, courier1)
	s.Require().NoError(err)

	courier2 := &model.CreateCourierRequest{
		Name:          "John",
		Lastname:      "Doe",
		Phone:         "+79990009988",
		Status:        "busy",
		TransportType: "bike",
	}

	id2, err := s.repo.Create(ctx, courier2)

	s.Require().Error(err)
	s.Equal(model.ErrPhoneConflict, err)

	s.Equal(id2, 0)
}

func (s *CourierRepositoryTestSuite) TestUpdate_Success() {
    ctx := context.Background()

    _, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
    s.Require().NoError(err)

    var id int
    err = s.pool.QueryRow(ctx, `
        INSERT INTO couriers (name, lastname, phone, status, transport_type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, "Old", "Name", "phone-old", "available", "on_foot").Scan(&id)
    s.Require().NoError(err)

    req := &model.UpdateCourierRequest{
        Id:            id,
        Name:          "New",
        Lastname:      "Surname",
        Phone:         "phone-new",
        Status:        "busy",
        TransportType: "bike",
    }

    err = s.repo.Update(ctx, req)
    s.Require().NoError(err)


    var got model.Courier
    err = s.pool.QueryRow(ctx, `
        SELECT id, name, lastname, phone, status, transport_type
        FROM couriers
        WHERE id = $1
    `, id).Scan(&got.Id, &got.Name, &got.Lastname, &got.Phone, &got.Status, &got.TransportType)
    s.Require().NoError(err)

    s.Equal(req.Id, got.Id)
    s.Equal(req.Name, got.Name)
    s.Equal(req.Lastname, got.Lastname)
    s.Equal(req.Phone, got.Phone)
    s.Equal(req.Status, got.Status)
    s.Equal(req.TransportType, got.TransportType)
}

func (s *CourierRepositoryTestSuite) TestUpdate_NotFound() {
    ctx := context.Background()

    req := &model.UpdateCourierRequest{
        Id:            999999999,
        Name:          "X",
        Lastname:      "Y",
        Phone:         "phone-x",
        Status:        "available",
        TransportType: "on_foot",
    }

    err := s.repo.Update(ctx, req)

    s.Require().Error(err)
    s.ErrorIs(err, model.ErrCourierNotFound)
}

func (s *CourierRepositoryTestSuite) TestUpdate_PhoneConflict() {
    ctx := context.Background()

    _, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
    s.Require().NoError(err)

    var id1, id2 int
    err = s.pool.QueryRow(ctx, `
        INSERT INTO couriers (name, lastname, phone, status, transport_type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, "A", "One", "phone-111", "available", "on_foot").Scan(&id1)
    s.Require().NoError(err)

    err = s.pool.QueryRow(ctx, `
        INSERT INTO couriers (name, lastname, phone, status, transport_type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, "B", "Two", "phone-222", "available", "bike").Scan(&id2)
    s.Require().NoError(err)

    req := &model.UpdateCourierRequest{
        Id:            id2,
        Name:          "B",
        Lastname:      "Two",
        Phone:         "phone-111",
        Status:        "busy",
        TransportType: "bike",
    }

    err = s.repo.Update(ctx, req)
    s.Require().Error(err)
    s.ErrorIs(err, model.ErrPhoneConflict)
}

func (s *CourierRepositoryTestSuite) TestDelete_Success() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `DELETE FROM couriers`)
	s.Require().NoError(err)

	var id int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO couriers (name, lastname, phone, status, transport_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "Delete", "Me", "delete-111", "available", "on_foot").Scan(&id)
	s.Require().NoError(err)

	err = s.repo.Delete(ctx, id)
	s.Require().NoError(err)

	var count int
	err = s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM couriers WHERE id = $1
	`, id).Scan(&count)
	s.Require().NoError(err)
	s.Equal(0, count)
}

func (s *CourierRepositoryTestSuite) TestDelete_NotFound() {
	ctx := context.Background()

	err := s.repo.Delete(ctx, 999999999) 

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrCourierNotFound)
}


func TestCourierRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}

func (s *CourierRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()

	// for _, id := range s.couriersId {
	// 	_, err := s.pool.Exec(ctx, `
	// 		DELETE FROM couriers WHERE id = $1	
	// 	`, id)
	// 	s.Require().NoError(err)
	// }

	_, err := s.pool.Exec(ctx, `TRUNCATE couriers`)
	s.Require().NoError(err)

	s.couriersId = nil
	
}

func getConnectionString() (string, error) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		return "", err
	}

	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		return "", err
	}

	return connString, nil
}

func mustInitPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString, err := getConnectionString()
	if err != nil {
		return nil, err
	}

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pingAttemptsLimit := 3
	var pingErr error

	for i := range pingAttemptsLimit {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = pool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Printf("db ping attempt %d failed: %v", i, pingErr)
		if i < pingAttemptsLimit {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if pingErr != nil {
		return nil, pingErr
	}

	log.Println("Database test connection pool established")
	return pool, nil
}
