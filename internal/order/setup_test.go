package order

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestSetup(t *testing.T) {
	t.Run("should register POST /order route", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mock repo
		mockRepo := NewMockRepository(ctrl)

		// echo instance
		e := echo.New()

		// pass nil for promo validator
		Setup(e, mockRepo, nil)

		// verify route
		routes := e.Routes()
		found := false
		for _, r := range routes {
			if r.Method == http.MethodPost && r.Path == "/order" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected POST /order to be registered but it was not")
		}
	})
}
