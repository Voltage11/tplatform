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
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewDepartmentToResponse(department *domain.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:        department.ID.String(),
		Name:      department.Name,
		CreatedAt: department.CreatedAt.Format(time.RFC3339),
		UpdatedAt: department.UpdatedAt.Format(time.RFC3339),
	}
}

func NewDepartmentsToResponse(departments []*domain.Department) []*DepartmentResponse {
	if departments == nil {
		return nil
	}

	out := make([]*DepartmentResponse, len(departments))
	for i, department := range departments {
		out[i] = NewDepartmentToResponse(department)
	}
	return out
}
