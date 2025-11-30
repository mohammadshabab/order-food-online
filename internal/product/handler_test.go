package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/internal/logger"
)

func TestHandler_ListProducts(t *testing.T) {
	logger.Init("test-service", "test", 0)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockService(ctrl)
		expected := []*Product{{ID: "p1", Name: "Burger"}}

		mockSvc.EXPECT().ListProducts(gomock.Any()).Return(expected, nil)

		h := NewHandler(mockSvc)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ListProducts(c)
		if err != nil {
			t.Fatalf("handler returned err: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var resp []*Product
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp[0].ID != "p1" {
			t.Errorf("unexpected response: %+v", resp)
		}
	})

	t.Run("service_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockService(ctrl)
		mockSvc.EXPECT().ListProducts(gomock.Any()).Return(nil, errors.New("db failed"))

		h := NewHandler(mockSvc)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ListProducts(c)
		if err != nil {
			t.Fatalf("handler returned err: %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", rec.Code)
		}
	})
}

func TestHandler_GetProduct(t *testing.T) {
	logger.Init("test-service", "test", 0)
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockService(ctrl)

		validID := "3f6b5b2a-7f66-4b3f-9a1b-000000000000"
		expected := &Product{ID: validID, Name: "Burger"}

		mockSvc.EXPECT().
			GetProduct(gomock.Any(), validID).
			Return(expected, nil)

		h := NewHandler(mockSvc)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/"+validID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("productId")
		c.SetParamValues(validID)

		err := h.GetProduct(c)
		if err != nil {
			t.Fatalf("handler returned err: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("missing_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockService(ctrl)
		h := NewHandler(mockSvc)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// No param
		err := h.GetProduct(c)
		if err != nil {
			t.Fatalf("handler returned err: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid_uuid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockService(ctrl)
		h := NewHandler(mockSvc)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/x1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("productId")
		c.SetParamValues("x1")

		// Handler returns 400 before calling service
		err := h.GetProduct(c)
		if err != nil {
			t.Fatalf("handler returned err: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
