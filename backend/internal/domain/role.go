package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ********************* РОЛЬ *********************

// Role — роль пользователя
type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

type RoleDetail struct {
	ID   uuid.UUID
	Name string
}

type RoleFilter struct {
	Name       string
	Pagination PaginationRequest
}

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetList(ctx context.Context, filter RoleFilter) ([]*Role, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RoleService interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetList(ctx context.Context, filter RoleFilter) (*PagedResult[*Role], error)
	Delete(ctx context.Context, id uuid.UUID) error
}
