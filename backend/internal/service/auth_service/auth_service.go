package auth_service

import (
	"context"
	"math/rand/v2"
	"strconv"
	"time"
	"tplatform/internal/apperr"
	"tplatform/internal/models"
	"tplatform/internal/repository/auth_repository"
	"tplatform/pkg/logger"
	"tplatform/pkg/hash"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, registrationRequest *models.RegistrationRequest, ipAddress, userAgent string) error
	Activate(ctx context.Context, token, verifyCode string) (*models.User, error)
	VerifyJwt(tokenString string) (*models.CurrentUser, error)
	Login(ctx context.Context, email, password, ipAddress, userAgent string) (*models.SessionResponse, error)
	RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (*models.SessionResponse, error)
}

type authService struct {
	repo      auth_repository.AuthRepository
	logger    logger.AppLogger
	jwtSecret string
}

func New(repo auth_repository.AuthRepository, logger logger.AppLogger, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, registrationRequest *models.RegistrationRequest, ipAddress, userAgent string) error {
	op := "AuthService.Register"

	// Проверим не занят ли email
	existsUser, err := s.repo.GetUserByEmail(ctx, registrationRequest.Email)
	if err != nil {
		s.logger.Error(err, op, "email", registrationRequest.Email)
		return err
	}

	if existsUser != nil {
		s.logger.Warn("Попытка регистрации с занятым email", op, "email", registrationRequest.Email)
		return apperr.BadRequest(nil, "Email занят", op)
	}

	// Проверим, отправлял ли пользователь заявку на регистрацию последние 15 минут
	existsRegistration, err := s.repo.GetRegistrationActiveByEmail(ctx, registrationRequest.Email)
	if err != nil {
		s.logger.Error(err, op, "email", registrationRequest.Email)
		return err
	}

	if existsRegistration != nil {
		// Если с момента регистрации не прошло 15 минут, продублируем отправку email
		if existsRegistration.CreatedAt.After(time.Now().Add(-15 * time.Minute)) {
			s.logger.Info("Повторная отправка кода подтверждения", op, "email", registrationRequest.Email)
			go existsRegistration.SendEmailConfirm()
			return nil
		}
	}

	// Хешируем пароль
	passwordHashed, err := hash.HashPassword(registrationRequest.Password)
	if err != nil {
		s.logger.Error(err, op)
		return apperr.Internal(err, op)
	}

	// Создаем регистрацию
	registration := models.Registration{
		Name:           registrationRequest.Name,
		Email:          registrationRequest.Email,
		PasswordHashed: passwordHashed,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		CreatedAt:      time.Now(),
		IsActive:       true,
		ExpiredAt:      time.Now().Add(time.Hour * 24),
		Token:          uuid.New().String(),
		VerifyCode:     generateVerifyCode(),
	}

	// Сохраняем регистрацию
	if err := s.repo.CreateRegistration(ctx, &registration); err != nil {
		return err
	}

	// Отправляем email подтверждения
	s.logger.Info("Отправка кода подтверждения", op, "email", registration.Email)
	go registration.SendEmailConfirm()

	return nil
}

func (s *authService) Activate(ctx context.Context, token, verifyCode string) (*models.User, error) {
	op := "AuthService.Activate"

	// Получаем регистрацию по токену
	registration, err := s.repo.GetRegistrationByToken(ctx, token)
	if err != nil {
		s.logger.Error(err, op, "token", token)
		return nil, err
	}

	if registration == nil {
		s.logger.Warn("Попытка активации по несуществующему токену", op, "token", token)
		return nil, apperr.NotFound("Не найдена запись регистрации для активации", op)
	}

	// Проверяем валидность регистрации
	if !registration.IsValidForActivate() {
		s.logger.Warn("Попытка активации просроченного аккаунта", op,
			"email", registration.Email,
			"expired_at", registration.ExpiredAt)
		return nil, apperr.BadRequest(nil, "Ссылка для активации истекла, зарегистрируйтесь снова", op)
	}

	// Проверяем код подтверждения
	if registration.VerifyCode != verifyCode {
		s.logger.Warn("Неверный код подтверждения", op,
			"email", registration.Email,
			"expected_code", registration.VerifyCode,
			"received_code", verifyCode)
		return nil, apperr.BadRequest(nil, "Неверный проверочный код", op)
	}

	// Создаем пользователя из регистрации
	user, err := s.repo.CreateUserFromRegistration(ctx, registration)
	if err != nil {
		s.logger.Error(err, op, "registration_id", registration.ID)
		return nil, err
	}

	s.logger.Info("Пользователь активирован", op,
		"user_id", user.ID,
		"email", user.Email)

	return user, nil
}

func (s *authService) VerifyJwt(tokenString string) (*models.CurrentUser, error) {
	op := "AuthService.VerifyJwt"

	currentUser, err := verifyJwt(tokenString, s.jwtSecret)
	if err != nil || currentUser == nil {
		s.logger.Warn("Невалидный JWT токен", op, "error", err)
		return nil, apperr.Unauthorized(op)
	}

	return currentUser, nil
}

func (s *authService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*models.SessionResponse, error) {
	op := "AuthService.Login"

	// Всегда выполняем хеширование для консистентного времени
	dummyHash := "$2a$10$dummyhashfordummycomparison123456"

	// Получаем пользователя
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(err, op, "email", email)
		return nil, err
	}

	// Используем фиктивный хеш если пользователь не найден
	hashToCheck := dummyHash
	if user != nil {
		hashToCheck = user.PasswordHashed
	}

	// ВСЕГДА проверяем пароль (консистентное время)
	if !hash.CheckPasswordHash(password, hashToCheck) {
		s.logger.Warn("Неверные учетные данные", op, "email", email)
		return nil, apperr.BadRequest(nil, "Неверные учетные данные", op)
	}

	// Только ПОСЛЕ проверки пароля проверяем существование пользователя
	if user == nil {
		s.logger.Warn("Пользователь не найден", op, "email", email)
		return nil, apperr.BadRequest(nil, "Неверные учетные данные", op)
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		s.logger.Warn("Попытка входа заблокированного пользователя", op,
			"user_id", user.ID,
			"email", user.Email)
		return nil, apperr.BadRequest(nil, "Пользователь заблокирован", op)
	}

	// Генерируем access token
	accessToken, err := generateJwt(user, time.Hour, s.jwtSecret)
	if err != nil {
		s.logger.Error(err, op, "user_id", user.ID)
		return nil, apperr.Internal(err, op)
	}

	// Генерируем refresh token
	refreshTokenUuid, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(err, op, "user_id", user.ID)
		return nil, apperr.Internal(err, op)
	}
	refreshToken := refreshTokenUuid.String()

	// Создаем сессию
	session := models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 7),
		CreatedAt:    time.Now(),
	}

	// Сохраняем сессию и обновляем время входа
	if err := s.repo.CreateSessionAndLogin(ctx, &session); err != nil {
		s.logger.Error(err, op, "user_id", user.ID)
		return nil, err
	}

	s.logger.Info("Успешный вход", op,
		"user_id", user.ID,
		"email", user.Email,
		"ip", ipAddress)

	go s.repo.SetUserLastLogin(nil, user.ID)

	return &models.SessionResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (*models.SessionResponse, error) {
	op := "AuthService.RefreshToken"

	// Получаем текущую сессию
	currentSession, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error(err, op, "refresh_token", refreshToken)
		return nil, err
	}

	if currentSession == nil {
		s.logger.Warn("Сессия не найдена", op, "refresh_token", refreshToken)
		return nil, apperr.BadRequest(nil, "Неверные учетные данные", op)
	}

	// Проверяем срок действия refresh token
	if time.Now().After(currentSession.ExpiredAt) {
		s.logger.Warn("Refresh token истек", op,
			"refresh_token", refreshToken,
			"expired_at", currentSession.ExpiredAt)
		return nil, apperr.BadRequest(nil, "Токен истек", op)
	}

	// Получаем пользователя
	currentUser, err := s.repo.GetUserByID(ctx, currentSession.UserID)
	if err != nil {
		s.logger.Error(err, op, "user_id", currentSession.UserID)
		return nil, err
	}

	if currentUser == nil {
		s.logger.Warn("Пользователь не найден", op, "user_id", currentSession.UserID)
		return nil, apperr.Unauthorized(op)
	}

	// Проверяем активность пользователя
	if !currentUser.IsActive {
		s.logger.Warn("Попытка обновления токена заблокированным пользователем", op,
			"user_id", currentUser.ID,
			"email", currentUser.Email)
		return nil, apperr.BadRequest(nil, "Пользователь заблокирован", op)
	}

	// Генерируем новый access token
	accessToken, err := generateJwt(currentUser, time.Hour, s.jwtSecret)
	if err != nil {
		s.logger.Error(err, op, "user_id", currentUser.ID)
		return nil, apperr.Internal(err, op)
	}

	// Генерируем новый refresh token
	newRefreshToken, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(err, op, "user_id", currentUser.ID)
		return nil, apperr.Internal(err, op)
	}

	// Создаем новую сессию
	newSession := models.Session{
		UserID:       currentUser.ID,
		RefreshToken: newRefreshToken.String(),
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 7),
		CreatedAt:    time.Now(),
	}

	// Обновляем сессию
	if err := s.repo.RefreshSession(ctx, currentSession.RefreshToken, &newSession); err != nil {
		s.logger.Error(err, op, "user_id", currentUser.ID)
		return nil, err
	}

	s.logger.Info("Токен обновлен", op,
		"user_id", currentUser.ID,
		"email", currentUser.Email,
		"new_refresh_token", newRefreshToken.String())

	return &models.SessionResponse{
		User:         *currentUser,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.String(),
	}, nil
}

// generateVerifyCode генерирует 5-значный код подтверждения
func generateVerifyCode() string {
	minValue := 10000
	maxValue := 99999
	randomNumber := rand.IntN(maxValue-minValue+1) + minValue
	return strconv.Itoa(randomNumber)
}
