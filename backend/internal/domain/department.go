package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Department отдел
type Department struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// DepartmentDetail сокр. отдел для связей с внешними таблицами
type DepartmentDetail struct {
	ID   uuid.UUID
	Name string
}

// DepartmentFilter Входящий параметры с отбором и пагинацией
type DepartmentFilter struct {
	Name       string
	Pagination PaginationRequest
}

// DepartmentRepository интерфейс для работы с репозиторием
type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	Update(ctx context.Context, department *Department) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetList(ctx context.Context, filter DepartmentFilter) ([]*Department, int64, error)
}

// DepartmentService интерфейс для работы с сервисом
type DepartmentService interface {
	Create(ctx context.Context, department *Department) error
	Update(ctx context.Context, department *Department) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetList(ctx context.Context, filter DepartmentFilter) (*PagedResult[*Department], error)
}
