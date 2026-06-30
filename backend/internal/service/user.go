package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Voltage11/tplatform/internal/cache"
	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/pkg/hash"
	"github.com/google/uuid"
)

const (
	userCacheDuration = 5 * time.Minute
)

type userService struct {
	repo  domain.UserRepository
	cache cache.UserCache
	permissionServise domain.PermissionService
}

func NewUserService(repo domain.UserRepository, permissionServise domain.PermissionService) domain.UserService {
	return &userService{
		repo:  repo,
		cache: cache.NewUserCache(userCacheDuration),
		permissionServise: permissionServise,
	}
}

func (u *userService) Create(ctx context.Context, user *domain.User) error {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionCreate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	user.PasswordHash = hash.ToHash(user.PasswordHash)

	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	u.cache.Set(user)
	return nil
}

func (u *userService) Update(ctx context.Context, user *domain.User) error {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}
	
	if err := u.repo.Update(ctx, user); err != nil {
		return err
	}

	u.cache.Set(user)
	return nil
}

func (u *userService) SetIsActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}
	
	if err := u.repo.SetIsActive(ctx, id, isActive); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) HardDelete(ctx context.Context, id uuid.UUID) error {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionHardDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}
	
	if err := u.repo.HardDelete(ctx, id); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionSoftDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}
	
	if err := u.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	u.cache.Delete(id)
	return nil
}

func (u *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}
	
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
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}
	
	return u.repo.GetByIDWithDetail(ctx, id)
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
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
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}
	
	return u.repo.GetByEmailWithDetail(ctx, email)
}

func (u *userService) GetList(ctx context.Context, filter domain.UserFilter) (*domain.PagedResult[*domain.UserWithDetail], error) {
	if ! u.permissionServise.CanFromCtx(ctx, domain.EntityUsers.Name, domain.ActionView.Name) {
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
	user, err := u.repo.GetByEmail(ctx, adminCfg.Email)
	if err != nil {
		if apperror.GetType(err) != apperror.ErrNotFound {
			return fmt.Errorf("ошибка при поиске администратора: %w", err)
		}
		// пользователь не найден — значит, создадим нового
		user = nil
	}

	// 1. Если пользователь уже существует
	if user != nil {
		// Уже активный админ — ничего не делаем
		if user.IsActive && user.IsAdmin {
			return nil
		}

		// Иначе обновляем: делаем активным и админом, не трогая пароль и остальные данные
		user.IsActive = true
		user.IsAdmin = true
		if err := u.Update(ctx, user); err != nil {
			return fmt.Errorf("не удалось обновить администратора: %w", err)
		}
		return nil
	}

	// 2. Пользователя нет — создаём нового администратора
	createdUser := &domain.User{
		FirstName:    "Админ",
		LastName:     "Администратор",
		Email:        adminCfg.Email,
		PasswordHash: adminCfg.Password, // пароль будет захэширован внутри Create
		IsActive:     true,
		IsAdmin:      true,
	}
	if err := u.Create(ctx, createdUser); err != nil {
		return fmt.Errorf("не удалось создать администратора: %w", err)
	}

	return nil
}
