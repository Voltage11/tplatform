package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Voltage11/tplatform/internal/appcontext"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/pkg/logger"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	validator   TokenValidator
	userService domain.UserService
	logger      logger.Logger
}

// TokenValidator интерфейс для валидации токена
type TokenValidator interface {
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

func NewAuthMiddleware(validator TokenValidator, userService domain.UserService, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		validator:   validator,
		userService: userService,
		logger:      log,
	}
}

func (u *AuthMiddleware) ExtractUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			next.ServeHTTP(w, r)
			return
		}

		token := parts[1]

		userID, err := u.validator.ValidateAccessToken(token)
		if err != nil {
			// Логируем причину
			u.logger.Error("ExtractUser: invalid token", "err", err)
			next.ServeHTTP(w, r)
			return
		}

		user, err := u.userService.GetByID(r.Context(), userID)
		if err != nil {
			u.logger.Error("ExtractUser: GetByID error", "userID", userID, "err", err)
			next.ServeHTTP(w, r)
			return
		}

		if !user.IsActive {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), appcontext.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (u *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appcontext.GetUserFromContext(r.Context())

		if user == nil {
			httputils.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		if !user.IsActive {
			httputils.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (u *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appcontext.GetUserFromContext(r.Context())

		if user == nil {
			httputils.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		if !(user.IsActive && user.IsAdmin) {
			httputils.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		next.ServeHTTP(w, r)
	})
}
