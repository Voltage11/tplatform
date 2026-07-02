package dto

import (
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
)

type DepartmentCreateRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

type DepartmentUpdateRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

type DepartmentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DepartmentToResponse(department *domain.Department) DepartmentResponse {
	return DepartmentResponse{
		ID:        department.ID.String(),
		Name:      department.Name,
		CreatedAt: department.CreatedAt,
		UpdatedAt: department.UpdatedAt,
	}
}
