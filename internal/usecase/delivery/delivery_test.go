package delivery

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAssignCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl) // <-- новое имя мока фабрики
	mockDeadlineCalc := NewMockDeadlineCalculator(ctrl)           // <-- новый мок калькулятора

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)

	ctx := t.Context()
	orderId := "AAA555"
	deadline := time.Date(2030, 1, 1, 10, 0, 0, 0, time.UTC) // фиксированное время, чтобы тест был стабильным

	courierWithFewestOrders := &model.Courier{Id: 10, TransportType: model.TransportTypeCar}

	mockCourierRepo.EXPECT().
		GetCourierIdWithFewestOrders(ctx).
		Return(10, nil)

	mockCourierRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(courierWithFewestOrders, nil)

	// Новая логика: фабрика возвращает калькулятор
	mockDeadlineFactory.EXPECT().
		GetDeliveryCalculator(courierWithFewestOrders.TransportType).
		Return(mockDeadlineCalc)

	// ...а калькулятор возвращает дедлайн
	mockDeadlineCalc.EXPECT().
		CalculateDeadline().
		Return(deadline)

	mockDeliveryRepo.EXPECT().
		AssignCourierWithUpdate(ctx, 10, orderId, deadline).
		Return(&model.Delivery{}, nil)

	want := &dto.AssignedDeliveryResponse{
		CourierId:      10,
		OrderId:        orderId,
		TransportType:  "car",
		Deadline:       deadline,
	}

	req := dto.AssignDeliveryRequest{OrderId: orderId}

	got, err := srv.AssignCourier(ctx, req)

	require.NoError(t, err)
	require.Equal(t, want.CourierId, got.CourierId)
	require.Equal(t, want.OrderId, got.OrderId)
	require.EqualValues(t, want.TransportType, got.TransportType)
	require.Equal(t, want.Deadline, got.Deadline)
}

func TestAssignCourier_ErrorOnGetCourierIdWithFewestOrders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	repoErr := errors.New("db error")

	mockCourierRepo.EXPECT().
		GetCourierIdWithFewestOrders(ctx).
		Return(0, repoErr)

	got, err := srv.AssignCourier(ctx, dto.AssignDeliveryRequest{OrderId: "AAA555"})

	require.ErrorIs(t, err, repoErr)
	require.Nil(t, got)
}

func TestAssignCourier_ErrorOnGetCourierById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	mockCourierRepo.EXPECT().
		GetCourierIdWithFewestOrders(ctx).
		Return(10, nil)

	repoErr := errors.New("db error")
	mockCourierRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(nil, repoErr)

	got, err := srv.AssignCourier(ctx, dto.AssignDeliveryRequest{OrderId: "AAA555"})

	require.ErrorIs(t, err, repoErr)
	require.Nil(t, got)
}

func TestAssignCourier_ErrorOnAssignCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)
	mockDeadlineCalc := NewMockDeadlineCalculator(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	deadline := time.Date(2030, 1, 1, 10, 0, 0, 0, time.UTC)

	mockCourierRepo.EXPECT().
		GetCourierIdWithFewestOrders(ctx).
		Return(10, nil)

	courier := &model.Courier{Id: 10, TransportType: model.TransportTypeCar}
	mockCourierRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(courier, nil)

	mockDeadlineFactory.EXPECT().
		GetDeliveryCalculator(courier.TransportType).
		Return(mockDeadlineCalc)

	mockDeadlineCalc.EXPECT().
		CalculateDeadline().
		Return(deadline)

	repoErr := errors.New("db error")
	mockDeliveryRepo.EXPECT().
		AssignCourierWithUpdate(ctx, courier.Id, gomock.Any(), gomock.Any()).
		Return(nil, repoErr)

	got, err := srv.AssignCourier(ctx, dto.AssignDeliveryRequest{OrderId: "AAA555"})

	require.ErrorIs(t, err, repoErr)
	require.Nil(t, got)
}

func TestUnassignCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	orderId := "AAA555"
	returned := &model.Delivery{CourierId: 10}

	mockDeliveryRepo.EXPECT().
		UnassignWithUpdate(ctx, orderId).
		Return(returned, nil)

	want := &dto.UnassignedDeliveryResponse{
		OrderId:   orderId,
		CourierId: 10,
		Status:    "unassigned",
	}

	got, err := srv.UnassignCourier(ctx, dto.UnassignDeliveryRequest{OrderId: orderId})

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestUnassignCourier_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	orderId := "AAA555"

	repoErr := errors.New("db error")
	mockDeliveryRepo.EXPECT().
		UnassignWithUpdate(ctx, orderId).
		Return(nil, repoErr)

	got, err := srv.UnassignCourier(ctx, dto.UnassignDeliveryRequest{OrderId: orderId})

	require.ErrorIs(t, err, repoErr)
	require.Nil(t, got)
}

func TestReleaseExpiredDeliveries_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	mockCourierRepo.EXPECT().
		ReleaseCouriers(ctx).
		Return(nil)

	err := srv.ReleaseExpiredDeliveries(ctx)

	require.NoError(t, err)
}

func TestReleaseExpiredDeliveries_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)
	ctx := t.Context()

	repoErr := errors.New("db error")
	mockCourierRepo.EXPECT().
		ReleaseCouriers(ctx).
		Return(repoErr)

	err := srv.ReleaseExpiredDeliveries(ctx)

	require.ErrorIs(t, err, repoErr)
}

func TestRunDeliveryChecker_StopsOnContextCancel(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockCourierRepo := NewMockCourierRepository(ctrl)
	mockDeliveryRepo := NewMockDeliveryRepository(ctrl)
	mockDeadlineFactory := NewMockdeadlineCalculatorFactory(ctrl)

	srv := NewDeliveryUseCase(mockDeliveryRepo, mockCourierRepo, mockDeadlineFactory)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := make(chan struct{}, 1)

	mockCourierRepo.EXPECT().
		ReleaseCouriers(gomock.Any()).
		DoAndReturn(func(context.Context) error {
			select {
			case called <- struct{}{}:
			default:
			}
			return nil
		}).
		MinTimes(1)

	done := make(chan struct{})
	go func() {
		srv.RunDeliveryChecker(ctx, 5*time.Millisecond)
		close(done)
	}()

	select {
	case <-called:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected ReleaseCouriers to be called at least once")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected RunDeliveryChecker to stop after context cancel")
	}
}
