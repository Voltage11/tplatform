package domain

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Name       string
	IsCorrect  bool
	SortOrder  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
