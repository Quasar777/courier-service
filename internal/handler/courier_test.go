package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Quasar777/courier-service/internal/handler"
	srvmocks "github.com/Quasar777/courier-service/internal/handler/mocks"
	"github.com/Quasar777/courier-service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCourier_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/courier/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	expected := &model.Courier{Id: 5}
	mockUseCase.EXPECT().
		GetCourier(gomock.Any(), 5).
		Return(expected, nil)
	
	c.Get(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got model.Courier
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, *expected, got)
}

func TestCourier_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)
	
	req := httptest.NewRequest(http.MethodGet, "/courier/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)
	
	mockUseCase.EXPECT().
		GetCourier(gomock.Any(), 5).
		Return(nil, model.ErrCourierNotFound)

	c.Get(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.JSONEq(t, `{"error": "Courier not found"}`, rr.Body.String())
}

func TestCourier_Get_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)
	
	req := httptest.NewRequest(http.MethodGet, "/courier/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)
	
	dbErr := errors.New("db error")
	mockUseCase.EXPECT().
		GetCourier(gomock.Any(), 5).
		Return(nil, dbErr)

	c.Get(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error": "Database error"}`, rr.Body.String())
}

func TestCourier_Get_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)
	
	req := httptest.NewRequest(http.MethodGet, "/courier/sdfsdf", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "sdfsdf")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	c.Get(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error": "Invalid id"}`, rr.Body.String())
}

func TestCourier_GetMany_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)	
	req := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	rr := httptest.NewRecorder()

	expected := []model.Courier{
		{Id: 5},
		{Id: 6},
	}

	mockUseCase.EXPECT().
		GetCouriers(gomock.Any()).
		Return(expected, nil)

	c.GetMany(rr, req)

	var got []model.Courier
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, expected, got)
}

func TestCourier_GetMany_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t) 
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl) 
	c := handler.NewCourierController(mockUseCase) 
	req := httptest.NewRequest(http.MethodGet, "/couriers", nil) 
	rr := httptest.NewRecorder()

	srvErr := errors.New("error in service")
	mockUseCase.EXPECT().
		GetCouriers(gomock.Any()).
		Return(nil, srvErr)
	
	c.GetMany(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error": "Database error"}`, rr.Body.String())
}

func TestCourier_Create_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)	

	reqBody := `{
		"name": "testName",
		"lastname": "testLastName",
		"phone": "111",
		"status": "available",
		"transportType": "car"
	}`
	req := httptest.NewRequest(http.MethodPost, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	expectReq := model.CreateCourierRequest{ 
		Name: "testName", 
		Lastname: "testLastName", 
		Phone: "111", 
		Status: "available", 
		TransportType: "car", 
	}

	mockUseCase.EXPECT(). 
		CreateCourier(gomock.Any(), expectReq). 
		Return(10, nil) 
	
	c.Create(rr, req)
	
	want := `{
		"id": 10,
		"message": "Courier created succesfully"
	}`
	
	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.JSONEq(t, want, rr.Body.String())
}

func TestCourier_Create_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)	

	reqBody := `{sdfdsfsdf}sdfsdfsdf`
	req := httptest.NewRequest(http.MethodPost, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()
	
	
	c.Create(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error": "Invalid JSON"}`, rr.Body.String())
}

func TestCourier_Create_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUseCase := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUseCase)	

	reqBody := `{
		"name": "testName",
		"lastname": "testLastName",
		"phone": "111",
		"status": "available",
		"transportType": "car"
	}`
	req := httptest.NewRequest(http.MethodPost, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()
	
	srvErr := errors.New("Service error")
	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(0, srvErr)

	c.Create(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error": "Database error"}`, rr.Body.String())
}

func TestCourier_Update_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{
		"id": 10,
		"name": "NewName",
		"lastname": "NewLast",
		"phone": "111",
		"status": "available",
		"transportType": "car"
	}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	expectReq := model.UpdateCourierRequest{
		Id:            10,
		Name:          "NewName",
		Lastname:      "NewLast",
		Phone:         "111",
		Status:        "available",
		TransportType: "car",
	}

	mockUC.EXPECT().
		UpdateCourier(gomock.Any(), expectReq).
		Return(nil)

	c.Update(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.JSONEq(t, `{"message":"Courier updated successfully"}`, rr.Body.String())
}

func TestCourier_Update_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader("{bad json"))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().UpdateCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Invalid JSON"}`, rr.Body.String())
}

func TestCourier_Update_IdRequired(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{"name":"X"}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	mockUC.EXPECT().UpdateCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Id is required"}`, rr.Body.String())
}

func TestCourier_Update_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{"id":10,"name":"X"}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	expectReq := model.UpdateCourierRequest{Id: 10, Name: "X"}

	mockUC.EXPECT().
		UpdateCourier(gomock.Any(), expectReq).
		Return(model.ErrCourierNotFound)

	c.Update(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.JSONEq(t, `{"error":"Courier not found"}`, rr.Body.String())
}

func TestCourier_Update_MissingFields(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{"id":10}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	expectReq := model.UpdateCourierRequest{Id: 10}

	mockUC.EXPECT().
		UpdateCourier(gomock.Any(), expectReq).
		Return(model.ErrMissingRequiredFields)

	c.Update(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Missing required fields"}`, rr.Body.String())
}

func TestCourier_Update_PhoneConflict(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{"id":10,"phone":"111"}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	expectReq := model.UpdateCourierRequest{Id: 10, Phone: "111"}

	mockUC.EXPECT().
		UpdateCourier(gomock.Any(), expectReq).
		Return(model.ErrPhoneConflict)

	c.Update(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
	require.JSONEq(t, `{"error":"Courier with this phone is already exists"}`, rr.Body.String())
}

func TestCourier_Update_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	reqBody := `{"id":10,"name":"X"}`
	req := httptest.NewRequest(http.MethodPut, "/couriers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	expectReq := model.UpdateCourierRequest{Id: 10, Name: "X"}
	dbErr := errors.New("db error")

	mockUC.EXPECT().
		UpdateCourier(gomock.Any(), expectReq).
		Return(dbErr)

	c.Update(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error":"Database error"}`, rr.Body.String())
}

func TestCourier_Delete_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	req := httptest.NewRequest(http.MethodDelete, "/couriers/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	mockUC.EXPECT().
		DeleteCourier(gomock.Any(), 5).
		Return(nil)

	c.Delete(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.JSONEq(t, `{"message":"Courier deleted successfully"}`, rr.Body.String())
}

func TestCourier_Delete_InvalidID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	req := httptest.NewRequest(http.MethodDelete, "/couriers/abc", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	mockUC.EXPECT().DeleteCourier(gomock.Any(), gomock.Any()).Times(0)

	c.Delete(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"Invalid id"}`, rr.Body.String())
}

func TestCourier_Delete_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	req := httptest.NewRequest(http.MethodDelete, "/couriers/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	mockUC.EXPECT().
		DeleteCourier(gomock.Any(), 5).
		Return(errors.New("courier not found"))

	c.Delete(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.JSONEq(t, `{"error":"Courier not found"}`, rr.Body.String())
}

func TestCourier_Delete_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUC := srvmocks.NewMockCourierUseCase(ctrl)
	c := handler.NewCourierController(mockUC)

	req := httptest.NewRequest(http.MethodDelete, "/couriers/5", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	mockUC.EXPECT().
		DeleteCourier(gomock.Any(), 5).
		Return(errors.New("db error"))

	c.Delete(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"error":"Database error"}`, rr.Body.String())
}