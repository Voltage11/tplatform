package jwt

import (
	"fmt"
	"time"

	"github.com/Voltage11/tplatform/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	SecretKey  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	cfg Config
}

func NewJWTService(cfg Config) *JWTService {
	return &JWTService{cfg: cfg}
}

// GenerateAccessToken создаёт короткоживущий access-токен
func (j *JWTService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.SecretKey))
}

// GenerateRefreshToken генерирует случайную строку
func (j *JWTService) GenerateRefreshToken() (string, error) {
	return hash.GenerateNewRandomStr()
}

// ValidateAccessToken проверяет access-токен и возвращает userID как строку
func (j *JWTService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.SecretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil // возвращаем uuid.UUID
	}
	return uuid.Nil, fmt.Errorf("invalid token")
}

// GetAccessTTL возвращает время жизни access токена
func (j *JWTService) GetAccessTTL() time.Duration {
	return j.cfg.AccessTTL
}

// GetRefreshTTL возвращает время жизни refresh токена
func (j *JWTService) GetRefreshTTL() time.Duration {
	return j.cfg.RefreshTTL
}