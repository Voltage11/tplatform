package auth_service

import (
	"time"
	"tplatform/internal/models"

	//"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	models.CurrentUser
	jwt.RegisteredClaims
}

// generateJwt создает JWT токен
func generateJwt(user *models.User, expiresIn time.Duration, jwtSecret string) (string, error) {
	claims := userClaims{
		CurrentUser: models.CurrentUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_session",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// verifyJwt проверяет и расшифровывает JWT токен с учетом срока действия
func verifyJwt(tokenString string, jwtSecret string) (*models.CurrentUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err // автоматически включает проверку срока действия
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return &claims.CurrentUser, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
