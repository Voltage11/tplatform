package validators

import (
	"fmt"
	"tplatform/internal/apperr"

	"github.com/go-playground/validator/v10"
)

// Validator интерфейс для валидации
type Validator interface {
	Validate(i any, op string) error
}

// ValidationError структура ошибки валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
}

// ValidationErrorResponse ответ с ошибками валидации
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func (v *ValidationErrorResponse) Error() string {
	return "Validation failed"
}

// appValidator реализация валидатора
type appValidator struct {
	validate *validator.Validate
}

// NewValidator создает новый валидатор
func NewValidator() Validator {
	v := validator.New()
	return &appValidator{
		validate: v,
	}
}

// Validate проверяет структуру и возвращает AppError при ошибке
func (av *appValidator) Validate(i any, op string) error {
	err := av.validate.Struct(i)
	if err == nil {
		return nil
	}

	// Преобразуем ошибки валидации в наш формат
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationError
		for _, err := range validationErrors {
			validationErr := ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: av.getErrorMessage(err),
			}
			errors = append(errors, validationErr)
		}

		return &ValidationErrorResponse{Errors: errors}
	}

	// Любые другие ошибки валидации
	return apperr.BadRequest(nil, "Ошибка валидации", op)
}

func (av *appValidator) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Это поле обязательно для заполнения"
	case "min":
		return fmt.Sprintf("Минимальная длина: %s", err.Param())
	case "max":
		return fmt.Sprintf("Максимальная длина: %s", err.Param())
	case "email":
		return "Неверный формат email"
	case "numeric":
		return "Должно содержать только цифры"
	default:
		return "Некорректное значение"
	}
}
