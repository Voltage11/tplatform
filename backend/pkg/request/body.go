package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseRequestBody Универсальная функция распарсить тело запроса, вернуть структуру, с дженериками
func ParseRequestBody[T any](r *http.Request) (*T, error) {
	if r.Method == http.MethodGet {
		return nil, nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("Content-Type должен быть application/json")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела запроса: %v", err)
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil, nil
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге тела запроса: %v", err)
	}
	return &result, nil
}
