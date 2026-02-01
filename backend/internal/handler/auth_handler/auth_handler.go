package auth_handler

import (
	"net/http"
	"time"
	"tplatform/internal/apperr"
	"tplatform/internal/middleware"
	"tplatform/internal/models"
	"tplatform/internal/service/auth_service"
	"tplatform/pkg/logger"
	"tplatform/pkg/request"
	"tplatform/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
)

type authHandler struct {
	service auth_service.AuthService
	logger  logger.AppLogger
}

func New(r *chi.Mux, service auth_service.AuthService, logger logger.AppLogger) {
	if r == nil {
		panic("auth_handler.New: получен nil router")
	}

	if service == nil {
		panic("auth_handler.New: получен nil service")
	}

	h := &authHandler{
		service: service,
		logger:  logger,
	}

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, 1*time.Minute))
		r.Post("/register", h.register)
		r.Post("/activate/{token}", h.activate)
		r.Post("/login", h.login)
		r.Get("/profile", h.profile)
		r.Post("/refresh-token", h.refreshToken)
	})
}

func (a *authHandler) register(w http.ResponseWriter, r *http.Request) {
	op := "auth_handler.register"

	// Парсим запрос
	registration, err := request.ParseRequestBody[models.RegistrationRequest](r)
	if err != nil {
		response.WriteError(w, apperr.BadRequest(err, "Неверный формат запроса", op))
		return
	}
	if registration == nil {
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный формат запроса", op))
		return
	}

	// Валидируем запрос
	if err := registration.Validate(); err != nil {
		response.WriteError(w, err)
		return
	}

	// Вызываем сервис
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	if err := a.service.Register(r.Context(), registration, ipAddress, userAgent); err != nil {
		a.logger.Error(err, op, "email", registration.Email)
		response.WriteError(w, err)
		return
	}

	a.logger.Info("Регистрация успешна", op, "email", registration.Email)
	response.WriteSuccess(w, "Успешная регистрация, проверьте email для подтверждения")
}

func (a *authHandler) activate(w http.ResponseWriter, r *http.Request) {
	op := "auth_handler.activate"

	// Парсим код подтверждения из тела запроса
	type VerifyCode struct {
		Code string `json:"code" validate:"required"`
	}

	verifyCode, err := request.ParseRequestBody[VerifyCode](r)
	if err != nil {
		response.WriteError(w, apperr.BadRequest(err, "Неверный формат запроса", op))
		return
	}
	if verifyCode == nil {
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный формат запроса", op))
		return
	}

	// Валидируем код
	if verifyCode.Code == "" {
		response.WriteError(w, apperr.BadRequestWithoutError("Код подтверждения обязателен", op))
		return
	}

	// Получаем токен из URL
	tokenFromPath := chi.URLParam(r, "token")
	if _, err := uuid.Parse(tokenFromPath); err != nil {
		a.logger.Warn("Невалидный токен в URL", op, "token", tokenFromPath)
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный URL", op))
		return
	}

	// Вызываем сервис
	user, err := a.service.Activate(r.Context(), tokenFromPath, verifyCode.Code)
	if err != nil {
		a.logger.Error(err, op, "token", tokenFromPath)
		response.WriteError(w, err)
		return
	}

	a.logger.Info("Активация успешна", op,
		"user_id", user.ID,
		"email", user.Email)

	response.WriteSuccess(w, user)
}

func (a *authHandler) login(w http.ResponseWriter, r *http.Request) {
	op := "auth_handler.login"

	// Парсим запрос
	loginReq, err := request.ParseRequestBody[loginRequest](r)
	if err != nil {
		response.WriteError(w, apperr.BadRequest(err, "Неверный формат запроса", op))
		return
	}
	if loginReq == nil {
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный формат запроса", op))
		return
	}

	// Валидируем запрос
	if err := loginReq.validate(); err != nil {
		response.WriteError(w, err)
		return
	}

	// Вызываем сервис
	sessionResponse, err := a.service.Login(r.Context(),
		loginReq.Email,
		loginReq.Password,
		r.RemoteAddr,
		r.UserAgent())

	if err != nil {
		a.logger.Error(err, op, "email", loginReq.Email)
		response.WriteError(w, err)
		return
	}

	a.logger.Info("Вход успешен", op, "email", loginReq.Email)
	response.WriteSuccess(w, sessionResponse)
}

func (a *authHandler) refreshToken(w http.ResponseWriter, r *http.Request) {
	op := "auth_handler.refreshToken"

	// Парсим запрос
	type refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	refreshTokenReq, err := request.ParseRequestBody[refreshTokenRequest](r)
	if err != nil {
		response.WriteError(w, apperr.BadRequest(err, "Неверный формат запроса", op))
		return
	}
	if refreshTokenReq == nil {
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный формат запроса", op))
		return
	}

	// Валидируем refresh token
	if _, err := uuid.Parse(refreshTokenReq.RefreshToken); err != nil {
		a.logger.Warn("Невалидный refresh token", op, "token", refreshTokenReq.RefreshToken)
		response.WriteError(w, apperr.BadRequestWithoutError("Неверный refresh token", op))
		return
	}

	// Вызываем сервис
	sessionResponse, err := a.service.RefreshToken(r.Context(),
		refreshTokenReq.RefreshToken,
		r.RemoteAddr,
		r.UserAgent())

	if err != nil {
		a.logger.Error(err, op, "refresh_token", refreshTokenReq.RefreshToken)
		response.WriteError(w, err)
		return
	}

	a.logger.Info("Refresh token успешен", op, "user_id", sessionResponse.User.ID)
	response.WriteSuccess(w, sessionResponse)
}

func (a *authHandler) profile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUserFromContext(r.Context())
	if user != nil && user.IsActive {
		response.WriteSuccess(w, user)
		return
	}
	response.WriteError(w, apperr.Unauthorized("Пользователь не авторизован"))
}
