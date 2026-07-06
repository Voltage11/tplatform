package domain

import (
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID            uuid.UUID
	Name          string
	Description   string
	IsActive      bool
	CreatedByID   uuid.UUID
	CreatedAt     time.Time
	DateBegin     *time.Time
	DateEnd       *time.Time
	MaxPoint      int
	CheckPoint    int
	ImgPath       string
	CorrectPoints int
	DeletedAt     *time.Time
}

type ThemeWithDetail struct {
	Theme
	CreatedBy User
}
