package auth_repository

import (
	"context"
	"errors"
	"time"
	"tplatform/internal/apperr"
	"tplatform/internal/models"
	"tplatform/pkg/database"
	"tplatform/pkg/logger"

	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	CreateRegistration(ctx context.Context, registration *models.Registration) error
	GetRegistrationByToken(ctx context.Context, token string) (*models.Registration, error)
	GetRegistrationActiveByEmail(ctx context.Context, email string) (*models.Registration, error)
	CreateUserFromRegistration(ctx context.Context, registration *models.Registration) (*models.User, error)
	SetUserLastLogin(ctx context.Context, userID int)

	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	DeleteSessionByRefreshToken(ctx context.Context, refreshToken string) error
	RefreshSession(ctx context.Context, oldRefreshToken string, newSession *models.Session) error

	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateSessionAndLogin(ctx context.Context, session *models.Session) error
}

type authRepository struct {
	db     *database.Database
	logger logger.AppLogger
}

func New(db *database.Database, logger logger.AppLogger) AuthRepository {
	return &authRepository{
		db:     db,
		logger: logger,
	}
}

func (r *authRepository) CreateRegistration(ctx context.Context, registration *models.Registration) error {
	op := "auth_repository.CreateRegistration"

	// До создания новой регистрации - деактивируем предыдущие заявки по email
	if err := r.deactivateUnvalidRegistrationsByEmail(ctx, registration.Email); err != nil {
		return err
	}

	query := `
        INSERT INTO registrations (name, email, password_hashed, ip_address, user_agent, created_at, is_active, expired_at, activated_at, token, verify_code)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `

	err := r.db.Pool.QueryRow(ctx, query,
		registration.Name,
		registration.Email,
		registration.PasswordHashed,
		registration.IPAddress,
		registration.UserAgent,
		registration.CreatedAt,
		registration.IsActive,
		registration.ExpiredAt,
		registration.ActivatedAt,
		registration.Token,
		registration.VerifyCode).Scan(&registration.ID)

	if err != nil {
		r.logger.Error(err, op, "email", registration.Email)
		return apperr.HandleDBError(err, "не удалось создать регистрацию", op)
	}

	r.logger.Info("Регистрация создана", op, "registration_id", registration.ID)
	return nil
}

func (r *authRepository) GetRegistrationByToken(ctx context.Context, token string) (*models.Registration, error) {
	op := "auth_repository.GetRegistrationByToken"

	query := `
        SELECT id, name, email, password_hashed, ip_address, user_agent, created_at, is_active, expired_at, activated_at, token, verify_code
        FROM registrations
        WHERE token = $1
    `

	row := r.db.Pool.QueryRow(ctx, query, token)

	var registration models.Registration
	err := row.Scan(
		&registration.ID,
		&registration.Name,
		&registration.Email,
		&registration.PasswordHashed,
		&registration.IPAddress,
		&registration.UserAgent,
		&registration.CreatedAt,
		&registration.IsActive,
		&registration.ExpiredAt,
		&registration.ActivatedAt,
		&registration.Token,
		&registration.VerifyCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("Регистрация не найдена", op, "token", token)
			return nil, nil
		}
		r.logger.Error(err, op, "token", token)
		return nil, apperr.HandleDBError(err, "Получение регистрации по токену", op)
	}

	return &registration, nil
}

func (r *authRepository) GetRegistrationActiveByEmail(ctx context.Context, email string) (*models.Registration, error) {
	op := "auth_repository.GetRegistrationActiveByEmail"

	query := `
        SELECT id, name, email, password_hashed, ip_address, user_agent, created_at, is_active, expired_at, activated_at, token, verify_code
        FROM registrations
        WHERE email = $1 AND
              is_active = TRUE AND
              activated_at IS NULL AND
              expired_at > CURRENT_TIMESTAMP
        ORDER BY created_at DESC
        LIMIT 1
    `

	row := r.db.Pool.QueryRow(ctx, query, email)

	var registration models.Registration
	err := row.Scan(
		&registration.ID,
		&registration.Name,
		&registration.Email,
		&registration.PasswordHashed,
		&registration.IPAddress,
		&registration.UserAgent,
		&registration.CreatedAt,
		&registration.IsActive,
		&registration.ExpiredAt,
		&registration.ActivatedAt,
		&registration.Token,
		&registration.VerifyCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error(err, op, "email", email)
		return nil, apperr.HandleDBError(err, "Получение активной регистрации по email", op)
	}

	return &registration, nil
}

func (r *authRepository) CreateUserFromRegistration(ctx context.Context, registration *models.Registration) (*models.User, error) {
	op := "auth_repository.CreateUserFromRegistration"

	queryCreateUser := `
        INSERT INTO users (name, email, password_hashed, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	queryDeactivateRegistration := `
        UPDATE registrations
            SET is_active = FALSE,
                activated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	user := models.User{
		Name:           registration.Name,
		Email:          registration.Email,
		PasswordHashed: registration.PasswordHashed,
		IsActive:       true,
		IsAdmin:        false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error(err, op, "registration_id", registration.ID)
		return nil, apperr.Internal(err, op)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// 1. Создаем пользователя
	err = tx.QueryRow(ctx, queryCreateUser,
		user.Name,
		user.Email,
		user.PasswordHashed,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		r.logger.Error(err, op, "email", user.Email)
		return nil, apperr.HandleDBError(err, "Создание пользователя на основании регистрации", op)
	}

	// 2. Деактивируем регистрацию
	if _, err = tx.Exec(ctx, queryDeactivateRegistration, registration.ID); err != nil {
		r.logger.Error(err, op, "registration_id", registration.ID)
		return nil, apperr.HandleDBError(err, "Деактивация регистрации", op)
	}

	// 3. Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		r.logger.Error(err, op, "user_id", user.ID)
		return nil, apperr.Internal(err, op)
	}

	r.logger.Info("Пользователь создан из регистрации", op,
		"user_id", user.ID,
		"registration_id", registration.ID)

	return &user, nil
}

func (r *authRepository) deactivateUnvalidRegistrationsByEmail(ctx context.Context, email string) error {
	op := "auth_repository.deactivateUnvalidRegistrationsByEmail"

	query := `
        UPDATE registrations
            SET is_active = FALSE
        WHERE is_active = TRUE AND
              expired_at > CURRENT_TIMESTAMP AND
              email = $1
    `

	_, err := r.db.Pool.Exec(ctx, query, email)
	if err != nil {
		r.logger.Error(err, op, "email", email)
		return apperr.HandleDBError(err, "Деактивация регистраций по email", op)
	}

	return nil
}

func (r *authRepository) SetUserLastLogin(ctx context.Context, userID int) {
	op := "auth_repository.SetUserLastLogin"

	if ctx == nil {
		ctx = context.Background()
	}

	query := `
        UPDATE users
            SET last_login_at = CURRENT_TIMESTAMP
        WHERE id = $1
	`
	if _, err := r.db.Pool.Exec(ctx, query, userID); err != nil {
		r.logger.Error(err, op, "user_id", userID)
	}
}

// SESSIONS
func (r *authRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	op := "auth_repository.GetSessionByRefreshToken"

	query := `
        SELECT id, user_id, refresh_token, user_agent, ip_address, expired_at, created_at
        FROM sessions
        WHERE refresh_token = $1
    `

	row := r.db.Pool.QueryRow(ctx, query, refreshToken)

	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiredAt,
		&session.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("Сессия не найдена", op, "refresh_token", refreshToken)
			return nil, nil
		}
		r.logger.Error(err, op, "refresh_token", refreshToken)
		return nil, apperr.HandleDBError(err, "Получение сессии по refresh_token", op)
	}

	return &session, nil
}

func (r *authRepository) DeleteSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	op := "auth_repository.DeleteSessionByRefreshToken"

	query := `DELETE FROM sessions WHERE refresh_token = $1`
	_, err := r.db.Pool.Exec(ctx, query, refreshToken)

	if err != nil {
		r.logger.Error(err, op, "refresh_token", refreshToken)
		return apperr.HandleDBError(err, "Удаление сессии по refresh_token", op)
	}

	r.logger.Info("Сессия удалена", op, "refresh_token", refreshToken)
	return nil
}

func (r *authRepository) RefreshSession(ctx context.Context, oldRefreshToken string, newSession *models.Session) error {
	op := "auth_repository.RefreshSession"

	queryDeleteOld := `DELETE FROM sessions WHERE refresh_token = $1`
	queryCreateNew := `
        INSERT INTO sessions (user_id, refresh_token, user_agent, ip_address, expired_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	queryUpdateLastLogin := `
        UPDATE users
            SET last_login_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error(err, op, "old_refresh_token", oldRefreshToken)
		return apperr.Internal(err, op)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// 1. Удаляем старую сессию
	if _, err := tx.Exec(ctx, queryDeleteOld, oldRefreshToken); err != nil {
		r.logger.Error(err, op, "old_refresh_token", oldRefreshToken)
		return apperr.HandleDBError(err, "Удаление старой сессии", op)
	}

	// 2. Создаем новую сессию
	err = tx.QueryRow(ctx, queryCreateNew,
		newSession.UserID,
		newSession.RefreshToken,
		newSession.UserAgent,
		newSession.IPAddress,
		newSession.ExpiredAt,
		newSession.CreatedAt).Scan(&newSession.ID)

	if err != nil {
		r.logger.Error(err, op, "user_id", newSession.UserID)
		return apperr.HandleDBError(err, "Создание новой сессии", op)
	}

	// 3. Обновляем время последнего входа
	if _, err := tx.Exec(ctx, queryUpdateLastLogin, newSession.UserID); err != nil {
		r.logger.Error(err, op, "user_id", newSession.UserID)
		return apperr.HandleDBError(err, "Обновление времени входа", op)
	}

	// 4. Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		r.logger.Error(err, op, "user_id", newSession.UserID)
		return apperr.Internal(err, op)
	}

	r.logger.Info("Сессия обновлена", op,
		"user_id", newSession.UserID,
		"new_refresh_token", newSession.RefreshToken)

	return nil
}

// USERS
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	op := "auth_repository.GetUserByEmail"

	query := `
        SELECT id, name, email, password_hashed, is_active, is_admin, last_login_at, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHashed,
		&user.IsActive,
		&user.IsAdmin,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error(err, op, "email", email)
		return nil, apperr.HandleDBError(err, "Получение пользователя по email", op)
	}

	return &user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	op := "auth_repository.GetUserByID"

	query := `
        SELECT id, name, email, password_hashed, is_active, is_admin, last_login_at, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHashed,
		&user.IsActive,
		&user.IsAdmin,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("Пользователь не найден", op, "user_id", id)
			return nil, nil
		}
		r.logger.Error(err, op, "user_id", id)
		return nil, apperr.HandleDBError(err, "Получение пользователя по ID", op)
	}

	return &user, nil
}

func (r *authRepository) CreateSessionAndLogin(ctx context.Context, session *models.Session) error {
	op := "auth_repository.CreateSessionAndLogin"

	queryCreateSession := `
        INSERT INTO sessions (user_id, refresh_token, user_agent, ip_address, expired_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	queryUpdateLastLogin := `
        UPDATE users
            SET last_login_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error(err, op, "user_id", session.UserID)
		return apperr.Internal(err, op)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// 1. Создаем сессию
	err = tx.QueryRow(ctx, queryCreateSession,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiredAt,
		session.CreatedAt).Scan(&session.ID)

	if err != nil {
		r.logger.Error(err, op, "user_id", session.UserID)
		return apperr.HandleDBError(err, "Создание сессии пользователя", op)
	}

	// 2. Обновляем время последнего входа
	if _, err = tx.Exec(ctx, queryUpdateLastLogin, session.UserID); err != nil {
		r.logger.Error(err, op, "user_id", session.UserID)
		return apperr.HandleDBError(err, "Установка времени входа пользователя", op)
	}

	// 3. Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		r.logger.Error(err, op, "user_id", session.UserID)
		return apperr.Internal(err, op)
	}

	r.logger.Info("Сессия создана", op,
		"user_id", session.UserID,
		"session_id", session.ID)

	return nil
}
