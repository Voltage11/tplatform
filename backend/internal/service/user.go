package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Voltage11/tplatform/internal/appcontext"
	"github.com/Voltage11/tplatform/internal/cache"
	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/pkg/password"
	"github.com/google/uuid"
)

const (
	userCacheDuration = 5 * time.Minute
)

type userService struct {
	repo             domain.UserRepository
	cache            cache.UserCache
	permisionServise domain.PermissionService
	txManager        domain.Transactor
}

func NewUserService(repo domain.UserRepository, permisionServise domain.PermissionService, txManager domain.Transactor) domain.UserService {
	return &userService{
		repo:             repo,
		cache:            cache.NewUserCache(userCacheDuration),
		permisionServise: permisionServise,
		txManager:        txManager,
	}
}

func (u *userService) Create(ctx context.Context, user *domain.User) error {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionCreate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	hashed, err := password.Hash(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}
	user.PasswordHash = hashed

	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	u.cache.Set(user)
	return nil
}

func (u *userService) Update(ctx context.Context, user *domain.User) error {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return err
	}

	u.cache.Set(user)
	return nil
}

func (u *userService) SetIsActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	if err := u.repo.SetIsActive(ctx, id, isActive); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) HardDelete(ctx context.Context, id uuid.UUID) error {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionHardDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	if err := u.repo.HardDelete(ctx, id); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionSoftDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	if err := u.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
	// 	return nil, apperror.NewForbiddenWithoutErr()
	// }

	// Проверяем кэш
	if user := u.cache.GetByID(id); user != nil {
		return user, nil
	}

	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.cache.Set(user)
	return user, nil
}

func (u *userService) GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*domain.UserWithDetail, error) {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return u.repo.GetByIDWithDetail(ctx, id)
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	// Проверяем кэш
	if user := u.cache.GetByEmail(email); user != nil {
		return user, nil
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	u.cache.Set(user)
	return user, nil
}

func (u *userService) GetByEmailWithDetail(ctx context.Context, email string) (*domain.UserWithDetail, error) {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return u.repo.GetByEmailWithDetail(ctx, email)
}

func (u *userService) GetList(ctx context.Context, filter domain.UserFilter) (*domain.PagedResult[*domain.UserWithDetail], error) {
	if !u.permisionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	users, total, err := u.repo.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	pagination := domain.NewPaginationResponse(filter.Pagination, total)

	return &domain.PagedResult[*domain.UserWithDetail]{
		Data:       users,
		Pagination: pagination,
	}, nil
}

func (u *userService) CheckOrCreateAdmin(ctx context.Context, adminCfg config.AdminConfig) error {

	userToCtx := domain.User{IsAdmin: true}
	ctx = appcontext.SetUserToContext(ctx, &userToCtx)

	// Оборачиваем всю операцию проверки и создания/обновления в транзакцию
	return u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Ищем админа внутри транзакции
		user, err := u.repo.GetByEmail(txCtx, adminCfg.Email)
		if err != nil {
			if apperror.GetType(err) != apperror.ErrNotFound {
				return fmt.Errorf("ошибка при поиске администратора: %w", err)
			}
			user = nil
		}

		// 2. Если админ найден, при необходимости обновляем его статус
		if user != nil {
			if user.IsActive && user.IsAdmin {
				return nil
			}

			user.IsActive = true
			user.IsAdmin = true
			// Вызываем метод репозитория напрямую с txCtx
			if err := u.repo.Update(txCtx, user); err != nil {
				return fmt.Errorf("не удалось обновить администратора: %w", err)
			}
			// Кэш обновляем внутри замыкания, он применится только если Commit пройдет успешно
			u.cache.Set(user)
			return nil
		}

		// 3. Если админа нет — создаём нового
		createdUser := &domain.User{
			ID:           uuid.New(),
			FirstName:    "Админ",
			LastName:     "Администратор",
			Email:        adminCfg.Email,
			PasswordHash: adminCfg.Password,
			IsActive:     true,
			IsAdmin:      true,
		}

		hashed, err := password.Hash(createdUser.PasswordHash)
		if err != nil {
			return fmt.Errorf("ошибка хеширования пароля админа: %w", err)
		}
		createdUser.PasswordHash = hashed

		// Записываем в БД внутри транзакции
		if err := u.repo.Create(txCtx, createdUser); err != nil {
			return fmt.Errorf("не удалось создать администратора: %w", err)
		}

		u.cache.Set(createdUser)
		return nil
	})
}

func (u *userService) Shutdown() {
	u.cache.Stop()
}
