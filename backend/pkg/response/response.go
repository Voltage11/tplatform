package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"tplatform/internal/apperr"
	"tplatform/pkg/validators"
)

// WriteSuccess записывает успешный ответ
func WriteSuccess(w http.ResponseWriter, data any) {
	response := SuccessResponse{
		Success: true,
		Result:  data,
	}
	writeJSON(w, http.StatusOK, response)
}

// WriteError записывает ошибку в ответ
func WriteError(w http.ResponseWriter, err error) {
	// Проверяем, является ли ошибка ошибкой валидации
	var validationErr *validators.ValidationErrorResponse
	if errors.As(err, &validationErr) {
		WriteValidationErrors(w, validationErr)
		return
	}

	// Проверяем, является ли ошибка AppError
	var appErr *apperr.AppError
	if errors.As(err, &appErr) {
		// Логируем внутренние ошибки
		if appErr.Type == "INTERNAL_ERROR" && appErr.Err != nil {
			// logger.Error("Internal error", "op", appErr.Op, "error", appErr.Err)
		}

		errorResponse := ErrorResponse{
			Success: false,
			Error: map[string]interface{}{
				"type":    appErr.Type,
				"message": appErr.Message,
				"code":    appErr.Type,
			},
		}
		writeJSON(w, appErr.Code, errorResponse)
		return
	}

	// Неизвестная ошибка - создаем внутреннюю
	appErr = apperr.Internal(err, "unknown")
	WriteError(w, appErr)
}

func WriteValidationErrors(w http.ResponseWriter, validationErr *validators.ValidationErrorResponse) {
	errorResponse := ErrorResponse{
		Success: false,
		Error: map[string]interface{}{
			"type":    "VALIDATION_ERROR",
			"message": "Ошибка валидации",
			"details": validationErr.Errors,
		},
	}
	writeJSON(w, http.StatusBadRequest, errorResponse)
}

// writeJSON вспомогательная функция для записи JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// logger.Error("Failed to encode JSON response", "error", err)
	}
}
