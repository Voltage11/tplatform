package dto

import (
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
)

type RoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description,omitempty"`
}

type RoleUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description,omitempty"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func RoleToResponse(role *domain.Role) *RoleResponse {
	return &RoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
	}
}

func RolesToResponseSlice(roles []*domain.Role) []*RoleResponse {
	if roles == nil {
		return nil
	}

	out := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		out[i] = RoleToResponse(role)
	}
	return out
}
