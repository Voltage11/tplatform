package domain

import (
	"time"

	"github.com/google/uuid"
)

type QuestionType string

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
