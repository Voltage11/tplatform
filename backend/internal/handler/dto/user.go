package dto

import "github.com/Voltage11/tplatform/internal/domain"

type UserResponse struct {
	ID           string              `json:"id"`
	FirstName    string              `json:"first_name"`
	SecondName   string              `json:"second_name,omitempty"`
	LastName     string              `json:"last_name"`
	Email        string              `json:"email"`
	IsActive     bool                `json:"is_active"`
	IsAdmin      bool                `json:"is_admin"`
	// Department   *DepartmentResponse `json:"department,omitempty"`
	// Role         *RoleResponse       `json:"role,omitempty"`      
	CreatedAt    string              `json:"created_at,omitempty"`
	UpdatedAt    string              `json:"updated_at,omitempty"`
}

// type RoleResponse struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

func NewUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsAdmin:    user.IsAdmin,
	}
}