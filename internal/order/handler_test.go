package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/mohammadshabab/order-food-online/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateOrder(t *testing.T) {
	// Initialize logger to avoid nil panics
	logger.Init("test-service", "test", 0)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc)
	e := echo.New()

	t.Run("invalid JSON bind", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte("{invalid-json")))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, ErrOrderInvalid.Code, rec.Code)
	})

	t.Run("validation fails", func(t *testing.T) {
		body := OrderReq{Items: &[]OrderItem{}}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, 422, rec.Code) // order must have at least one item
	})

	t.Run("service returns error", func(t *testing.T) {
		body := OrderReq{Items: &[]OrderItem{{ProductID: "p1", Quantity: 1}}}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockSvc.EXPECT().CreateOrder(gomock.Any(), &body).Return(nil, errors.New("db error"))

		err := h.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		body := OrderReq{Items: &[]OrderItem{{ProductID: "p1", Quantity: 2}}}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Use Order instead of OrderResp
		expectedOrder := &Order{ID: "order1", Items: *body.Items}
		mockSvc.EXPECT().CreateOrder(gomock.Any(), &body).Return(expectedOrder, nil)

		err := h.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp Order
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder.ID, resp.ID)
		assert.Equal(t, expectedOrder.Items, resp.Items)
	})
}
