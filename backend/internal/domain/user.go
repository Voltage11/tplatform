package domain

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/types/filterbool"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	SecondName   string
	LastName     string
	Email        string
	PasswordHash string
	DepartmentID *uuid.UUID
	RoleID       *uuid.UUID
	IsActive     bool
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type UserWithDetail struct {
	User
	Department *DepartmentDetail
	Role       *RoleDetail
}

// UserFilter Входящий параметры с отбором и пагинацией
type UserFilter struct {
	FirstName    string
	SecondName   string
	LastName     string
	Email        string
	DepartmentID *uuid.UUID
	RoleID       *uuid.UUID
	IsActive     filterbool.FilterBool
	Pagination   PaginationRequest
}

// UserRepository Интерфейс для работы с репозиторием
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	SetIsActive(ctx context.Context, id uuid.UUID, isActive bool) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*UserWithDetail, error)
	GetList(ctx context.Context, filter UserFilter) ([]*UserWithDetail, int64, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailWithDetail(ctx context.Context, email string) (*UserWithDetail, error)
}

// UserService Интерфейс для работы с сервисом
type UserService interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	SetIsActive(ctx context.Context, id uuid.UUID, isActive bool) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*UserWithDetail, error)
	GetList(ctx context.Context, filter UserFilter) (*PagedResult[*UserWithDetail], error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailWithDetail(ctx context.Context, email string) (*UserWithDetail, error)
	CheckOrCreateAdmin(ctx context.Context, adminCfg config.AdminConfig) error
	Shutdown()
}
