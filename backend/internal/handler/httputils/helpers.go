package httputils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/internal/types/filterbool"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// WriteJSON отправляет JSON-ответ с заданным статусом
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteOk успешный ответ
func WriteOk(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, data)
}

// WriteErrorString отправляет ошибку в формате JSON
func WriteErrorString(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, map[string]string{"error": message})
}

// WriteError отправляет ошибку в формате JSON
func WriteError(w http.ResponseWriter, err error) {
	var code int
	var message string

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		message = validationErrs.Error()
	} else if appErr, ok := err.(*apperror.AppError); ok {
		code = apperror.HTTPStatusFromError(appErr)
		message = appErr.Message
	} else {
		code = http.StatusInternalServerError
		message = err.Error()
	}

	WriteJSON(w, code, map[string]string{"error": message})
}

// ParseUUID извлекает и валидирует UUID из параметра URL
func ParseUUID(r *http.Request, paramName string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, paramName)
	return uuid.Parse(idStr)
}

// ParsePagination извлекает page и limit из query-параметров
func ParsePagination(r *http.Request) domain.PaginationRequest {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return domain.PaginationRequest{Page: page, Limit: limit}
}

func GetQueryValue(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryValueWithExist(r *http.Request, key string) (string, bool) {
	ok := r.URL.Query().Has(key)

	if !ok {
		return "", false
	}

	return r.URL.Query().Get(key), true
}

func GetQueryValueInt(r *http.Request, key string) (int, bool) {
	valueStr, ok := GetQueryValueWithExist(r, key)

	if !ok {
		return 0, false
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, false
	}

	return valueInt, true
}

func ParseFilterBool(r *http.Request, key string) filterbool.FilterBool {
	val, ok := GetQueryValueWithExist(r, key)
	if !ok {
		return filterbool.FilterBoolAll
	}
	switch val {
	case "true":
		return filterbool.FilterBoolTrue
	case "false":
		return filterbool.FilterBoolFalse
	default:
		return filterbool.FilterBoolAll
	}
}

func ParseUUIDQuery(r *http.Request, key string) *uuid.UUID {
	val, ok := GetQueryValueWithExist(r, key)
	if !ok || val == "" {
		return nil
	}
	id, err := uuid.Parse(val)
	if err != nil {
		return nil
	}
	return &id
}
