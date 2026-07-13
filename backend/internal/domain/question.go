package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type QuestionType string

func (q QuestionType) String() string {
	return string(q)
}

const (
	QuestionTypeSingle QuestionType = "single"
	QuestionTypeMulti  QuestionType = "multi"
)

type Question struct {
	ID           uuid.UUID
	ThemeID      uuid.UUID
	QuestionType QuestionType
	Name         string
	PointCorrect int
	SortOrder    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type QuestionRepository interface{
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
	Create(ctx context.Context, question *Question) error
	GetLastPositionInTheme(ctx context.Context, themeID uuid.UUID) (int, error)
	Update(ctx context.Context, question *Question) error
	SetSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetList(ctx context.Context, themeID uuid.UUID) ([]*Question, error)
}

type QuestionService interface{
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
	Create(ctx context.Context, question *Question) error
	Update(ctx context.Context, question *Question) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetList(ctx context.Context, themeID uuid.UUID) ([]*Question, error)
	SetSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error
}