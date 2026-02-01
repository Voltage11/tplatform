package response

// SuccessResponse успешный ответ
type SuccessResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}

type ListResponse[T any] struct {
	Data       []*T `json:"data"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
}

func NewListResponse[T any](data []*T, total, page, pageSize int) *ListResponse[T] {
	var totalPages int
	if pageSize == 0 {
		totalPages = 1
	} else {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return &ListResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
