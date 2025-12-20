package courier_test

import (
	"errors"
	"testing"

	"github.com/Quasar777/courier-service/internal/handler/dto"
	"github.com/Quasar777/courier-service/internal/model"
	"github.com/Quasar777/courier-service/internal/usecase/courier"
	"github.com/Quasar777/courier-service/internal/usecase/common/mocks"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCourier_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	ctx := t.Context()
	expected := &model.Courier{Id: 5}
	mockRepo.EXPECT().
		GetOneById(ctx, 5).
		Return(expected, nil)

	got, err := s.GetCourier(ctx, 5)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetCourier_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	repoErr := errors.New("Database error")
	ctx := t.Context()
	mockRepo.EXPECT().
		GetOneById(ctx, 5).
		Return(nil, repoErr)
	
	got, err := s.GetCourier(ctx, 5)

	require.ErrorIs(t, err, repoErr)
	require.Nil(t, got)
}

func TestGetCourier_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	ctx := t.Context()
	mockRepo.EXPECT().GetOneById(ctx, 10).Return(nil, model.ErrCourierNotFound)

	got, err := s.GetCourier(ctx, 10)

	require.ErrorIs(t, err, model.ErrCourierNotFound)
	require.Nil(t, got)
}

func TestGetCourier_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	got, err := s.GetCourier(t.Context(), -7)

	require.ErrorIs(t, err, model.ErrInvalidId)
	require.Nil(t, got)
}

func TestCreateCourier_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	request := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
	}
	expectedReq := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
		Status: "available",
		TransportType: "on_foot",
	}
	mockRepo.EXPECT().
		Create(gomock.Any(), &expectedReq).
		Return(1, nil)

	got, err := s.CreateCourier(t.Context(), request)

	require.NoError(t, err)
	require.Equal(t, 1, got)
}

func TestCreateCourier_SuccessWithPassedStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	request := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
		Status: "paused",
	}
	expectedReq := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
		Status: "paused",
		TransportType: "on_foot",
	}
	mockRepo.EXPECT().
		Create(gomock.Any(), &expectedReq).
		Return(1, nil)

	got, err := s.CreateCourier(t.Context(), request)

	require.NoError(t, err)
	require.Equal(t, 1, got)
}

func TestCreateCourier_SuccessWithPassedTransportType(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	request := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
		TransportType: "car",
	}
	expectedReq := dto.CreateCourierRequest{
		Name: "Andrew",
		Lastname: "Ravshanov",
		Phone: "+79998887766",
		Status: "available",
		TransportType: "car",
	}
	mockRepo.EXPECT().
		Create(gomock.Any(), &expectedReq).
		Return(1, nil)

	got, err := s.CreateCourier(t.Context(), request)

	require.NoError(t, err)
	require.Equal(t, 1, got)
}

func TestUpdateCourier_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)
	ctx := t.Context()
	req := dto.UpdateCourierRequest{
		Id:       10,
		Name:     "Andrew",
		Lastname: "",
		Phone:    "",
		Status:   "",
		TransportType: "",
	}
	existing := &model.Courier{
		Id:            10,
		Name:          "OldName",
		Lastname:      "Smith",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}
	mockRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(existing, nil)
	expectedReq := dto.UpdateCourierRequest{
		Id:            10,
		Name:          "Andrew",
		Lastname:      "Smith",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}
	mockRepo.EXPECT().
		Update(ctx, &expectedReq).
		Return(nil)

	err := s.UpdateCourier(ctx, req)

	require.NoError(t, err)
}

func TestUpdateCourier_GetOneByIdError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	req := dto.UpdateCourierRequest{
		Id:   10,
		Name: "NewName",
	}

	repoErr := errors.New("db error")
	mockRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(nil, repoErr)

	err := s.UpdateCourier(ctx, req)

	require.ErrorIs(t, err, repoErr)
}

func TestUpdateCourier_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	req := dto.UpdateCourierRequest{
		Id:   10,
		Name: "NewName",
	}

	mockRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(nil, model.ErrCourierNotFound)

	err := s.UpdateCourier(ctx, req)

	require.ErrorIs(t, err, model.ErrCourierNotFound)
}

func TestUpdateCourier_UpdateError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	req := dto.UpdateCourierRequest{
		Id:   10,
		Name: "Andrew",
	}

	existing := &model.Courier{
		Id:            10,
		Name:          "OldName",
		Lastname:      "Smith",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}

	mockRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(existing, nil)

	expectedReq := dto.UpdateCourierRequest{
		Id:            10,
		Name:          "Andrew",
		Lastname:      "Smith",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}

	repoErr := errors.New("update failed")
	mockRepo.EXPECT().
		Update(ctx, &expectedReq).
		Return(repoErr)

	err := s.UpdateCourier(ctx, req)

	require.ErrorIs(t, err, repoErr)
}

func TestUpdateCourier_AllFieldsProvided(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	req := dto.UpdateCourierRequest{
		Id:            10,
		Name:          "Andrew",
		Lastname:      "Ravshanov",
		Phone:         "+7999",
		Status:        "available",
		TransportType: "bike",
	}

	mockRepo.EXPECT().
		GetOneById(ctx, 10).
		Return(&model.Courier{Id: 10}, nil)

	mockRepo.EXPECT().
		Update(ctx, &req).
		Return(nil)

	err := s.UpdateCourier(ctx, req)

	require.NoError(t, err)
}

func TestDeleteCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	courierId := 10

	mockRepo.EXPECT().
		Delete(ctx, courierId).
		Return(nil)

	err := s.DeleteCourier(ctx, courierId)

	require.NoError(t, err)
}

func TestDeleteCourier_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCourierRepository(ctrl)
	s := courier.NewCourierUseCase(mockRepo)

	ctx := t.Context()
	courierId := 10

	repoErr := errors.New("delete failed")
	mockRepo.EXPECT().
		Delete(ctx, courierId).
		Return(repoErr)

	err := s.DeleteCourier(ctx, courierId)

	require.ErrorIs(t, err, repoErr)
}