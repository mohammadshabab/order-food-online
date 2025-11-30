package apperrors

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppError_Error(t *testing.T) {
	e := AppError{
		Code:    400,
		Message: "bad request",
		Err:     errors.New("invalid input"),
	}

	require.Equal(t, "400: bad request | cause: invalid input", e.Error())
}

func TestAppError_Error_NoCause(t *testing.T) {
	e := AppError{
		Code:    404,
		Message: "not found",
	}

	require.Equal(t, "404: not found", e.Error())
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("root error")
	e := &AppError{Err: cause}

	unwrappedErr := errors.Unwrap(e)
	require.Equal(t, cause, unwrappedErr)
}

func TestNew(t *testing.T) {
	e := New(401, "unauthorized", LevelWarn)

	require.Equal(t, 401, e.Code)
	require.Equal(t, "unauthorized", e.Message)
	require.Equal(t, LevelWarn, e.Level)
	require.Nil(t, e.Err)
}

func TestWrap_WithNilError_ReturnsNew(t *testing.T) {
	e := Wrap(500, "internal", LevelError, nil)

	require.Equal(t, 500, e.Code)
	require.Equal(t, "internal", e.Message)
	require.Nil(t, e.Err)
}

func TestWrap_WithAppError_ReturnsSameInstance(t *testing.T) {
	orig := AppError{Code: 400, Message: "existing"}

	e := Wrap(500, "ignored", LevelError, &orig)

	require.Equal(t, &orig, e)

}

func TestWrap_WithNormalError_Wraps(t *testing.T) {
	cause := errors.New("db error")

	e := Wrap(500, "internal", LevelError, cause)

	require.Equal(t, 500, e.Code)
	require.Equal(t, "internal", e.Message)
	require.Equal(t, cause, e.Err)
}

func TestHelpers(t *testing.T) {
	br := BadRequest("bad", nil)
	require.Equal(t, http.StatusBadRequest, br.Code)
	require.Equal(t, LevelWarn, br.Level)

	unauth := Unauthorized("auth", nil)
	require.Equal(t, http.StatusUnauthorized, unauth.Code)

	nf := NotFound("missing", nil)
	require.Equal(t, http.StatusNotFound, nf.Code)

	intErr := Internal("fail", nil)
	require.Equal(t, http.StatusInternalServerError, intErr.Code)
	require.Equal(t, LevelError, intErr.Level)

	debug := Debug("trace", nil)
	require.Equal(t, http.StatusOK, debug.Code)
	require.Equal(t, LevelDebug, debug.Level)
}

func TestMarshalJSON(t *testing.T) {
	e := AppError{
		Code:    404,
		Message: "not found",
	}

	data, err := e.MarshalJSON()
	require.NoError(t, err)

	var parsed struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	require.Equal(t, 404, parsed.Code)
	require.Equal(t, "not found", parsed.Message)
}
