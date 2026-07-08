package handler

import (
	"net/http"

	"github.com/Voltage11/tplatform/internal/appcontext"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/dto"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"
	"github.com/Voltage11/tplatform/internal/service"
	"github.com/Voltage11/tplatform/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	authService *service.AuthService
	userService domain.UserService
	logger      logger.Logger
	validate    *validator.Validate
}

func NewAuthHandlers(r chi.Router, authMW *middleware.AuthMiddleware, authService *service.AuthService, userService domain.UserService, log logger.Logger) {
	h := authHandler{
		authService: authService,
		userService: userService,
		logger:      log,
		validate:    validator.New(),
	}

	// Публичные маршруты
	r.Post("/api/v1/auth/login", h.Login)
	r.Post("/api/v1/auth/refresh", h.Refresh)
	r.Post("/api/v1/auth/logout", h.Logout)

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth) // Любая авторизация принимается, не важно админ или нет
		r.Get("/api/v1/profile", h.Profile)
	})
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest, err := httputils.DecodeJSONBodyWithValidate[dto.LoginRequest](r, a.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	token, err := a.authService.Login(r.Context(), loginRequest.Email, loginRequest.Password, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	})
}

func (a *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshRequest, err := httputils.DecodeJSONBodyWithValidate[dto.RefreshRequest](r, a.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	tokenPair, err := a.authService.Refresh(r.Context(), refreshRequest.RefreshToken)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logoutRequest, err := httputils.DecodeJSONBodyWithValidate[dto.LogoutRequest](r, a.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := a.authService.Logout(r.Context(), logoutRequest.RefreshToken); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "успешно вышли из системы"})
}

// Profile возвращает информацию о текущем пользователе
func (a *authHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user := appcontext.GetUserFromContext(r.Context())
	if user == nil {
		httputils.WriteErrorString(w, http.StatusUnauthorized, "пользователь не авторизован")
		return
	}

	userWithDetail, err := a.userService.GetByIDWithDetail(r.Context(), user.ID)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, map[string]any{
		"user": dto.UserToResponse(userWithDetail),
	})
}
