package product

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestSetup(t *testing.T) {
	t.Run("should register GET /product and GET /product/:productId routes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mock repo
		mockRepo := NewMockRepository(ctrl)

		// create echo instance
		e := echo.New()

		// call the Setup function to register routes
		Setup(e, mockRepo)

		// verify routes
		routes := e.Routes()
		foundList := false
		foundGet := false
		for _, r := range routes {
			if r.Method == http.MethodGet && r.Path == "/product" {
				foundList = true
			}
			if r.Method == http.MethodGet && r.Path == "/product/:productId" {
				foundGet = true
			}
		}

		if !foundList {
			t.Errorf("expected GET /product to be registered but it was not")
		}
		if !foundGet {
			t.Errorf("expected GET /product/:productId to be registered but it was not")
		}
	})
}
