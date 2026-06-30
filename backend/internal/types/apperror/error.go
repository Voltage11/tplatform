package apperror

import (
	"errors"
	"fmt"
)

type ErrType string

const (
	ErrNotFound      ErrType = "NOT_FOUND"
	ErrBadRequest    ErrType = "BAD_REQUEST"
	ErrAlreadyExists ErrType = "ALREADY_EXISTS"
	ErrUnauthorized  ErrType = "UNAUTHORIZED"
	ErrForbidden     ErrType = "FORBIDDEN"
	ErrInternal      ErrType = "INTERNAL"
)

var ErrForbiddenGeneric = errors.New("доступ запрещён")

type AppError struct {
	Type    ErrType `json:"type"`
	Message string  `json:"message"`
	Err     error   `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Конструкторы

func NewNotFound(msg string, err error) error {
	return &AppError{Type: ErrNotFound, Message: msg, Err: err}
}

func NewConflict(msg string, err error) error {
	return &AppError{Type: ErrAlreadyExists, Message: msg, Err: err}
}

func NewBadRequest(msg string, err error) error {
	return &AppError{Type: ErrBadRequest, Message: msg, Err: err}
}

func NewUnauthorized(msg string, err error) error {
	return &AppError{Type: ErrUnauthorized, Message: msg, Err: err}
}

func NewForbidden(msg string, err error) error {
	return &AppError{Type: ErrForbidden, Message: msg, Err: err}
}

func NewForbiddenWithoutErr() error {
	return &AppError{Type: ErrForbidden, Message: "нет прав", Err: ErrForbiddenGeneric}
}

func NewInternal(msg string, err error) error {
	return &AppError{Type: ErrInternal, Message: msg, Err: err}
}

// GetType извлекает тип ошибки
func GetType(err error) ErrType {
	if err == nil {
		return ""
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrInternal
}
