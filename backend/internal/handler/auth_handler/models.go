package auth_handler

import (
	"strings"
	"tplatform/internal/apperr"
	"tplatform/pkg/validators"
)

// loginRequest структура для входа
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (lr *loginRequest) validate() error {
	op := "loginRequest.validate"
	lr.Email = strings.ToLower(strings.TrimSpace(lr.Email))

	if !validators.IsEmailValid(lr.Email) {
		return apperr.BadRequest(nil, "Неверный формат email", op)
	}

	if len(lr.Password) < 5 || len(lr.Password) > 15 {
		return apperr.BadRequest(nil, "пароль должен быть от 5 до 15 символов", op)
	}

	return nil
}
