package dto

import "github.com/Voltage11/tplatform/internal/domain"

// PagedResponse – универсальная обёртка для списков с пагинацией
type PagedResponse[T any] struct {
    Data       []T                     `json:"data"`
    Pagination domain.PaginationResponse `json:"pagination"`
}

func NewPagedResponse[T any](data []T, pagination domain.PaginationResponse) *PagedResponse[T] {
    return &PagedResponse[T]{
        Data:       data,
        Pagination: pagination,
    }
}