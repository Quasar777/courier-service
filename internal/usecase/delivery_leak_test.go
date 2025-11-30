package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/model"
	"go.uber.org/goleak"
)

// FakeCourierRepo — полностью заглушка.
// Ничего не делает, не запускает горутины, возвращает пустые данные.
type FakeCourierRepo struct{}

func (f *FakeCourierRepo) GetOneById(ctx context.Context, id int) (*model.Courier, error) {
	return &model.Courier{
		Id:            id,
		Name:          "Fake",
		Lastname:      "Courier",
		Phone:         "0000000000",
		Status:        "available",
		TransportType: "on_foot",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (f *FakeCourierRepo) GetAll(ctx context.Context) ([]model.Courier, error) {
	return []model.Courier{}, nil
}

func (f *FakeCourierRepo) Create(ctx context.Context, courier *model.CreateCourierRequest) (int, error) {
	return 1, nil
}

func (f *FakeCourierRepo) Update(ctx context.Context, courier *model.UpdateCourierRequest) error {
	return nil
}

func (f *FakeCourierRepo) Delete(ctx context.Context, id int) error {
	return nil
}

// FakeDeliveryRepo — заглушка репозитория доставки.
// Методы просто возвращают стабильные значения, без логики.
type FakeDeliveryRepo struct{}

func (f *FakeDeliveryRepo) AssignCourierWithUpdate(ctx context.Context, courierId int, orderId string, deadline time.Time) (*model.Delivery, error) {
	return &model.Delivery{
		Id:         1,
		CourierId:  courierId,
		OrderId:    orderId,
		AssignedAt: time.Now(),
		Deadline:   deadline,
	}, nil
}

func (f *FakeDeliveryRepo) UnassignWithUpdate(ctx context.Context, orderId string) (*model.Delivery, error) {
	return &model.Delivery{
		Id:         1,
		CourierId:  1,
		OrderId:    orderId,
		AssignedAt: time.Now(),
		Deadline:   time.Now(),
	}, nil
}

func (f *FakeDeliveryRepo) ReleaseCouriers(ctx context.Context) error {
	// Ничего не делаем — в тестах нам важен только сам вызов
	return nil
}

func (f *FakeDeliveryRepo) GetCourierIdWithFewestOrders(ctx context.Context) (int, error) {
	// Возвращаем фиксированный id, чтобы юзкейс мог работать
	return 1, nil
}


func TestRunDeliveryChecker_NoLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u := &DeliveryUseCase{
		DeliveryRepo:     &FakeDeliveryRepo{},
		CourierRepo:      &FakeCourierRepo{},
	}

	// запускаем тикер
	go u.RunDeliveryChecker(ctx, 10*time.Millisecond)

	// даём ему поработать
	time.Sleep(50 * time.Millisecond)

	// останавливаем контекст
	cancel()

	// даём горутине завершиться
	time.Sleep(20 * time.Millisecond)
}