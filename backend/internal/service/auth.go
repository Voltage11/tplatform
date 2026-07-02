package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/pkg/hash"
	jwtpkg "github.com/Voltage11/tplatform/pkg/jwt"
)

// AuthService - экспортируемая структура (с большой буквы)
type AuthService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	jwt         *jwtpkg.JWTService
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// NewAuthService - экспортируемый конструктор (с большой буквы)
func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	jwtCfg jwtpkg.Config,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwt:         jwtpkg.NewJWTService(jwtCfg),
	}
}

// Login проверяет email/пароль, создаёт сессию и возвращает токены
func (a *AuthService) Login(ctx context.Context, email, password, userAgent, clientIP string) (*TokenPair, error) {
	// 1. Найти пользователя по email
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.NewUnauthorized("неверный email или пароль", nil)
	}
	// 2. Проверить пароль
	if !hash.IsValidHash(user.PasswordHash, password) {
		return nil, apperror.NewUnauthorized("неверный email или пароль", nil)
	}
	// 3. Проверить активен ли пользователь
	if !user.IsActive {
		return nil, apperror.NewForbidden("учётная запись деактивирована", nil)
	}

	// 4. Сгенерировать токены
	accessToken, err := a.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, apperror.NewInternal("ошибка генерации токена", err)
	}
	refreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, apperror.NewInternal("ошибка генерации refresh токена", err)
	}

	// 5. Сохранить сессию
	refreshHash := hash.ToHash(refreshToken)
	session := &domain.Session{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().UTC().Add(a.jwt.GetRefreshTTL()),
		CreatedAt:        time.Now().UTC(),
		UserAgent:        userAgent,
		ClientIP:         clientIP,
	}
	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessTTL().Seconds()),
	}, nil
}

// Logout отзывает сессию по refresh токену
func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hashToken := hash.ToHash(refreshToken)
	session, err := a.sessionRepo.GetByRefreshTokenHash(ctx, hashToken)
	if err != nil {
		return apperror.NewNotFound("сессия не найдена", nil)
	}
	return a.sessionRepo.Revoke(ctx, session.ID)
}

// Refresh выдаёт новую пару токенов по refresh токену
func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	hashToken := hash.ToHash(refreshToken)
	session, err := a.sessionRepo.GetByRefreshTokenHash(ctx, hashToken)
	if err != nil {
		return nil, apperror.NewUnauthorized("недействительный refresh токен", nil)
	}
	if session.RevokedAt != nil {
		return nil, apperror.NewUnauthorized("сессия отозвана", nil)
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, apperror.NewUnauthorized("refresh токен истёк", nil)
	}

	// Генерируем новый refresh
	newRefreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, apperror.NewInternal("ошибка генерации refresh токена", err)
	}
	newRefreshHash := hash.ToHash(newRefreshToken)

	newSession := &domain.Session{
		ID:               uuid.New(),
		UserID:           session.UserID,
		RefreshTokenHash: newRefreshHash,
		ExpiresAt:        time.Now().UTC().Add(a.jwt.GetRefreshTTL()),
		CreatedAt:        time.Now().UTC(),
		UserAgent:        session.UserAgent,
		ClientIP:         session.ClientIP,
	}
	if err := a.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	// Отзываем старую сессию
	_ = a.sessionRepo.Revoke(ctx, session.ID)

	// Новый access
	accessToken, err := a.jwt.GenerateAccessToken(session.UserID)
	if err != nil {
		return nil, apperror.NewInternal("ошибка генерации токена", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessTTL().Seconds()),
	}, nil
}

// ValidateAccessToken валидирует access токен и возвращает userID как строку (для middleware)
func (a *AuthService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwt.ValidateAccessToken(tokenString)
}
