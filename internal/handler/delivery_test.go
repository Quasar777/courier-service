package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Quasar777/courier-service/internal/handler"
	"github.com/Quasar777/courier-service/internal/handler/dto"
	srvmocks "github.com/Quasar777/courier-service/internal/handler/mocks"
	"github.com/Quasar777/courier-service/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestDelivery_Assign_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	reqBody := `{"orderId":"AAA555"}`
	req := httptest.NewRequest(http.MethodPost, "/deliveries/assign", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	deadline := time.Now().Add(5 * time.Minute)

	mockUC.EXPECT().
		AssignCourier(gomock.Any(), dto.AssignDeliveryRequest{OrderId: "AAA555"}).
		Return(&dto.AssignedDeliveryResponse{
			CourierId:      10,
			OrderId:        "AAA555",
			TransportType:  "car",
			Deadline:       deadline,
		}, nil)

	c.Assign(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got dto.AssignedDeliveryResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, 10, got.CourierId)
	require.Equal(t, "AAA555", got.OrderId)
	require.Equal(t, "car", got.TransportType)
	require.True(t, got.Deadline.Equal(deadline))
}

func TestDelivery_Assign_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/assign", strings.NewReader("{bad json"))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().AssignCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Assign(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Invalid JSON"}`, rr.Body.String())
}

func TestDelivery_Assign_MissingOrderId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/assign", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().AssignCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Assign(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Missing orderId"}`, rr.Body.String())
}


func TestDelivery_Assign_NoAvailableCouriers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/assign", strings.NewReader(`{"orderId":"AAA555"}`))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().
		AssignCourier(gomock.Any(), dto.AssignDeliveryRequest{OrderId: "AAA555"}).
		Return(nil, model.ErrNoAvailableCouriers)

	c.Assign(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
	require.JSONEq(t, `{"error":"No available couriers"}`, rr.Body.String())
}

func TestDelivery_Assign_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/assign", strings.NewReader(`{"orderId":"AAA555"}`))
	rr := httptest.NewRecorder()

	dbErr := errors.New("db error")
	mockUC.EXPECT().
		AssignCourier(gomock.Any(), dto.AssignDeliveryRequest{OrderId: "AAA555"}).
		Return(nil, dbErr)

	c.Assign(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error":"Database error"}`, rr.Body.String())
}

func TestDelivery_Unassign_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	reqBody := `{"orderId":"AAA555"}`
	req := httptest.NewRequest(http.MethodPost, "/deliveries/unassign", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockUC.EXPECT().
		UnassignCourier(gomock.Any(), dto.UnassignDeliveryRequest{OrderId: "AAA555"}).
		Return(&dto.UnassignedDeliveryResponse{
			OrderId:   "AAA555",
			CourierId: 10,
			Status:    "unassigned",
		}, nil)

	c.Unassign(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got dto.UnassignedDeliveryResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, dto.UnassignedDeliveryResponse{
		OrderId:   "AAA555",
		CourierId: 10,
		Status:    "unassigned",
	}, got)
}

func TestDelivery_Unassign_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/unassign", strings.NewReader("{bad json"))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().UnassignCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Unassign(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Invalid JSON"}`, rr.Body.String())
}

func TestDelivery_Unassign_MissingOrderId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/unassign", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().UnassignCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Unassign(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Missing orderId"}`, rr.Body.String())
}


func TestDelivery_Unassign_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/unassign", strings.NewReader(`{"orderId":"AAA555"}`))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().
		UnassignCourier(gomock.Any(), dto.UnassignDeliveryRequest{OrderId: "AAA555"}).
		Return(nil, model.ErrNoRelationFound)

	c.Unassign(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.JSONEq(t, `{"error":"Delivery not found"}`, rr.Body.String())
}

func TestDelivery_Unassign_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockDeliveryUseCase(ctrl)
	c := handler.NewDeliveryController(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/deliveries/unassign", strings.NewReader(`{"orderId":"AAA555"}`))
	rr := httptest.NewRecorder()

	dbErr := errors.New("db error")
	mockUC.EXPECT().
		UnassignCourier(gomock.Any(), dto.UnassignDeliveryRequest{OrderId: "AAA555"}).
		Return(nil, dbErr)

	c.Unassign(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error":"Database error"}`, rr.Body.String())
}