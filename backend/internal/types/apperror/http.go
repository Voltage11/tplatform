package apperror

import "net/http"

// HTTPStatusFromError возвращает HTTP статус код по ошибке
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch GetType(err) {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
