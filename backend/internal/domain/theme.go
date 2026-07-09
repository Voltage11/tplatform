package domain

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/types/filterbool"
	"github.com/google/uuid"
)

type Theme struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsActive    bool
	CreatedByID uuid.UUID
	CreatedAt   time.Time
	DateBegin   *time.Time
	DateEnd     *time.Time
	MaxPoint    int
	CheckPoint  int
	ImgPath     string
	DeletedAt   *time.Time
}

type ThemeWithDetail struct {
	Theme
	CreatedBy UserShort
}

type ThemeFilter struct {
	Name          string
	IsActive      filterbool.FilterBool
	CreatedByID   *uuid.UUID
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	Pagination    PaginationRequest
}

// Repository Репозиторий для работы с темами
type ThemeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Theme, error)
	GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*ThemeWithDetail, error)
	Create(ctx context.Context, theme *Theme) error
	Update(ctx context.Context, theme *Theme) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActive(ctx context.Context, id uuid.UUID, isActive bool) error
	GetList(ctx context.Context, filter ThemeFilter) ([]*ThemeWithDetail, int64, error)
}

// Service Сервис для работы с темами
type ThemeService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Theme, error)
	GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*ThemeWithDetail, error)
	Create(ctx context.Context, theme *Theme) error
	Update(ctx context.Context, theme *Theme) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActive(ctx context.Context, id uuid.UUID, isActive bool) error
	GetList(ctx context.Context, filter ThemeFilter) (*PagedResult[*ThemeWithDetail], error)
}
