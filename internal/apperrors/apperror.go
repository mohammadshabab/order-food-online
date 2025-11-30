package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Level   Level  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d: %s | cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code int, msg string, level Level) *AppError {
	return &AppError{Code: code, Message: msg, Level: level}
}

func Wrap(code int, msg string, level Level, err error) *AppError {
	if err == nil {
		return New(code, msg, level)
	}
	if ae, ok := err.(*AppError); ok {
		return ae
	}
	return &AppError{Code: code, Message: msg, Err: err, Level: level}
}

// Helpers
func Debug(msg string, err error) *AppError {
	return Wrap(http.StatusOK, msg, LevelDebug, err)
}
func BadRequest(msg string, err error) *AppError {
	return Wrap(http.StatusBadRequest, msg, LevelWarn, err)
}
func Unauthorized(msg string, err error) *AppError {
	return Wrap(http.StatusUnauthorized, msg, LevelWarn, err)
}
func NotFound(msg string, err error) *AppError {
	return Wrap(http.StatusNotFound, msg, LevelWarn, err)
}
func Internal(msg string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, msg, LevelError, err)
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	type out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	return json.Marshal(out{Code: e.Code, Message: e.Message})
}
