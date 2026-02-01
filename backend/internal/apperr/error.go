package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"-"`
	Err     error  `json:"-"`
	Op      string `json:"-"`
}

func (e *AppError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("[%s] %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

// Конструктор
func New(code int, errType, msg, op string, err error) *AppError {
	return &AppError{
		Code:    code,
		Type:    errType,
		Message: msg,
		Op:      op,
		Err:     err,
	}
}

// Типовые ошибки
func NotFound(msg, op string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", msg, op, nil)
}

func Internal(err error, op string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL", "Internal server error", op, err)
}

func BadRequest(err error, msg, op string) *AppError {
	displayMsg := msg
	if displayMsg == "" && err != nil {
		displayMsg = err.Error()
	}
	return New(http.StatusBadRequest, "BAD_REQUEST", displayMsg, op, err)
}

func BadRequestWithoutError(msg, op string) *AppError {
	return BadRequest(nil, msg, op)
}

func Unauthorized(op string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", "Не авторизован", op, nil)
}

func Forbidden(op string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", "", op, nil)
}

// Хелперы проверки
func IsType(err error, errType string) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Type == errType
}
