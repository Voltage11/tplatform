package domain

import "math"

// Пагинация request ********************
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

// GetOffset рассчитывает сдвиг для SQL-запроса (OFFSET)
func (p PaginationRequest) GetOffset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.GetLimit()
}

// Пагинация response ********************

// PaginationResponse метаданные для клиента
type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalCount  int64 `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
}

// NewPaginationResponse — конструктор метаданных ответа
func NewPaginationResponse(req PaginationRequest, totalCount int64) PaginationResponse {
	limit := req.GetLimit()
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	return PaginationResponse{
		CurrentPage: page,
		Limit:       limit,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	}
}

// PagedResult — объединяем данные и пагинацию
type PagedResult[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
