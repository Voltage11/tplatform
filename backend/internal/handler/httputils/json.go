package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/go-playground/validator/v10"
)

const maxBodySize = 1_048_576 // 1 MB

// DecodeJSONBody читает тело запроса, ограничивает размер и десериализует JSON
// Возвращает указатель на готовый объект или ошибку
func DecodeJSONBody[T any](r *http.Request) (*T, error) {
    // Ограничиваем размер во избежание атак
    r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

    var obj T

    decoder := json.NewDecoder(r.Body)    

    if err := decoder.Decode(&obj); err != nil {
        var syntaxErr *json.SyntaxError
        var unmarshalTypeErr *json.UnmarshalTypeError

        switch {
        case errors.Is(err, io.EOF):
            return nil, apperror.NewBadRequest("тело запроса пустое", err)
        case errors.As(err, &syntaxErr):
            return nil, apperror.NewBadRequest("некорректный JSON", err)
        case errors.As(err, &unmarshalTypeErr):
            return nil, apperror.NewBadRequest(
                "неверный тип поля "+unmarshalTypeErr.Field, err,
            )
        case err.Error() == "http: request body too large":
            return nil, apperror.NewBadRequest("тело запроса превышает допустимый размер", err)
        default:
            return nil, apperror.NewBadRequest("ошибка разбора JSON", err)
        }
    }
    
    if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
        return nil, apperror.NewBadRequest("тело запроса содержит лишние данные", err)
    }

    return &obj, nil
}

func DecodeJSONBodyWithValidate[T any](r *http.Request, validate *validator.Validate) (*T, error) {
    out, err := DecodeJSONBody[T](r)
    if err != nil {
        return nil, err
    }

    if err := validate.Struct(out); err != nil {		
		return nil, err
	}

    return out, nil
}