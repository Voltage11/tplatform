package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	RevokedAt        *time.Time
	UserAgent        string
	ClientIP         string
}

// SessionRepository интерфейс для работы с сессиями
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByRefreshTokenHash(ctx context.Context, hash string) (*Session, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

// SessionService интерфейс для сервиса
type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, refreshToken string, userAgent, clientIP string) (*Session, error)
	RevokeSession(ctx context.Context, refreshToken string) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
}
