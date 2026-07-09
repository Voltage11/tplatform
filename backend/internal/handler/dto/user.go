package dto

import (
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
)

// ---------- Запросы ----------

type UserCreateRequest struct {
	FirstName    string  `json:"first_name" validate:"required,min=1"`
	SecondName   *string `json:"second_name,omitempty"`
	LastName     string  `json:"last_name" validate:"required,min=1"`
	Email        string  `json:"email" validate:"required,email"`
	Password     string  `json:"password" validate:"required,min=6"`
	DepartmentID *string `json:"department_id,omitempty" validate:"omitempty,uuid4"`
	RoleID       *string `json:"role_id,omitempty" validate:"omitempty,uuid4"`
	IsActive     *bool   `json:"is_active,omitempty"`
	IsAdmin      *bool   `json:"is_admin,omitempty"`
}

type UserUpdateRequest struct {
	FirstName    *string `json:"first_name,omitempty"`
	SecondName   *string `json:"second_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	DepartmentID *string `json:"department_id,omitempty" validate:"omitempty,uuid4"`
	RoleID       *string `json:"role_id,omitempty" validate:"omitempty,uuid4"`
	IsActive     *bool   `json:"is_active,omitempty"`
	IsAdmin      *bool   `json:"is_admin,omitempty"`
}

type UserSetActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// ---------- Ответы ----------

type UserResponse struct {
	ID         string              `json:"id"`
	FirstName  string              `json:"first_name"`
	SecondName *string             `json:"second_name,omitempty"`
	LastName   string              `json:"last_name"`
	Email      string              `json:"email"`
	Department *DepartmentResponse `json:"department,omitempty"`
	Role       *RoleResponse       `json:"role,omitempty"`
	IsActive   bool                `json:"is_active"`
	IsAdmin    bool                `json:"is_admin"`
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
}

// UserToResponse конвертирует доменную модель (с деталями) в DTO
func UserToResponse(u *domain.UserWithDetail) *UserResponse {
	if u == nil {
		return nil
	}

	resp := &UserResponse{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}

	if u.SecondName != "" {
		resp.SecondName = &u.SecondName
	}

	if u.Department != nil {
		resp.Department = &DepartmentResponse{
			ID:   u.Department.ID.String(),
			Name: u.Department.Name,
		}
	}

	if u.Role != nil {
		resp.Role = &RoleResponse{
			ID:   u.Role.ID.String(),
			Name: u.Role.Name,
		}
	}

	return resp
}

// UsersToResponseSlice конвертирует список пользователей
func UsersToResponseSlice(users []*domain.UserWithDetail) []*UserResponse {
	if users == nil {
		return nil
	}
	out := make([]*UserResponse, len(users))
	for i, u := range users {
		out[i] = UserToResponse(u)
	}
	return out
}

type UserShortResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
