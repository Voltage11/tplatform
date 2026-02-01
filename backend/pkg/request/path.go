package request

import (
	"net/http"
	"strconv"
)

// GetIntFromRequest - получение значения пути из запроса
func GetIntFromRequest(r *http.Request, key string) (int, bool) {
	valueStr := r.PathValue(key)
	if valueStr == "" {
		return 0, false
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, false
	}
	return valueInt, true
}

// GetQueryBoolValueFromRequest - получение bool значения query из запроса
func GetQueryBoolValueFromRequest(r *http.Request, key string) (bool, bool) {
	query := r.URL.Query()
	if !query.Has(key) {
		return false, false
	}

	valueStr := query.Get(key)
	if valueStr == "" {
		return false, false
	}
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, false
	}
	return valueBool, true
}

// GetQueryValueFromRequest - получение query из запроса
func GetQueryValueFromRequest(r *http.Request, key string) (string, bool) {
	value := r.URL.Query()
	return value.Get(key), value.Has(key)
}

// GetQueryIntValueFromRequest - получение int query из запроса
func GetQueryIntValueFromRequest(r *http.Request, key string) (int, bool) {
	valueStr, ok := GetQueryValueFromRequest(r, key)
	if !ok {
		return 0, false
	}
	valueInt, err := strconv.Atoi(valueStr)

	if err != nil {
		return 0, false
	}

	return valueInt, true
}

// GetPaginateFromRequest - получение пагинации из запроса
func GetPaginateFromRequest(r *http.Request) (int, int) {
	page, ok := GetQueryIntValueFromRequest(r, "page")
	if !ok || page < 1 {
		page = 1
	}

	pageSize, ok := GetQueryIntValueFromRequest(r, "page_size")
	if !ok || pageSize < 1 {
		pageSize = 30
	}

	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}
