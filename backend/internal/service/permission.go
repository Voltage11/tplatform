package service

import (
	"context"
	"sync"
	"time"

	"github.com/Voltage11/tplatform/internal/appcontext"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/google/uuid"
)

type permissionService struct {
	repo     domain.PermissionsRepository
	cache    map[uuid.UUID][]domain.RolePermission
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
	cachedAt map[uuid.UUID]time.Time
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewPermissionService(repo domain.PermissionsRepository) domain.PermissionService {
	s := &permissionService{
		repo:     repo,
		cache:    make(map[uuid.UUID][]domain.RolePermission),
		cacheTTL: 5 * time.Minute,
		cachedAt: make(map[uuid.UUID]time.Time),
		stopCh:   make(chan struct{}),
	}
	// горутина очистки просроченного кэша
	go s.startCleaner()
	return s
}

// --- основные методы ---

func (s *permissionService) SetAction(ctx context.Context, roleID uuid.UUID, entity domain.PermissionEntity, actionName string) error {
	if err := s.repo.SetAction(ctx, roleID, entity.Name, actionName); err != nil {
		return err
	}
	s.invalidate(roleID)
	return nil
}

func (s *permissionService) RemoveAction(ctx context.Context, roleID uuid.UUID, entity domain.PermissionEntity, actionName string) error {
	if err := s.repo.RemoveAction(ctx, roleID, entity.Name, actionName); err != nil {
		return err
	}
	s.invalidate(roleID)
	return nil
}

func (s *permissionService) GetForRole(ctx context.Context, roleID uuid.UUID) ([]domain.PermissionEntityWithActions, error) {
	perms, err := s.getCached(roleID, func() ([]domain.RolePermission, error) {
		return s.repo.GetRolePermissions(ctx, roleID)
	})
	if err != nil {
		return nil, err
	}

	active := make(map[string]bool, len(perms))
	for _, p := range perms {
		active[p.EntityName+":"+p.ActionName] = true
	}

	base := domain.GetPermissionsEntitiesWithActions()
	result := make([]domain.PermissionEntityWithActions, len(base))
	for i, entity := range base {
		result[i].Entity = entity.Entity
		result[i].Actions = make([]domain.PermissionAction, len(entity.Actions))
		for j, action := range entity.Actions {
			result[i].Actions[j] = action.WithActive(active[entity.Entity.Name+":"+action.Name])
		}
	}
	return result, nil
}

func (s *permissionService) Can(ctx context.Context, user *domain.User, entityName, actionName string) bool {
	if user == nil {
		return false
	}

	if user.IsAdmin {
		return true
	}
	if user.RoleID == nil {
		return false
	}

	perms, err := s.getCached(*user.RoleID, func() ([]domain.RolePermission, error) {
		return s.repo.GetRolePermissions(ctx, *user.RoleID)
	})
	if err != nil {
		return false
	}

	for _, p := range perms {
		if p.EntityName == entityName && p.ActionName == actionName {
			return true
		}
	}
	return false
}

// CanFromCtx оптимизировал получение прав, пользователя берем из контекста
func (s *permissionService) CanFromCtx(ctx context.Context, entityName, actionName string) bool {
	user := appcontext.GetUserFromContext(ctx)

	return s.Can(ctx, user, entityName, actionName)
}

// --- кэш ---

func (s *permissionService) getCached(roleID uuid.UUID, loader func() ([]domain.RolePermission, error)) ([]domain.RolePermission, error) {
	// быстрое чтение
	s.cacheMu.RLock()
	if cachedAt, ok := s.cachedAt[roleID]; ok && time.Since(cachedAt) < s.cacheTTL {
		val := s.cache[roleID]
		s.cacheMu.RUnlock()
		return val, nil
	}
	s.cacheMu.RUnlock()

	// загрузка под write-lock
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	if cachedAt, ok := s.cachedAt[roleID]; ok && time.Since(cachedAt) < s.cacheTTL {
		return s.cache[roleID], nil
	}

	perms, err := loader()
	if err != nil {
		return nil, err
	}

	s.cache[roleID] = perms
	s.cachedAt[roleID] = time.Now()
	return perms, nil
}

func (s *permissionService) invalidate(roleID uuid.UUID) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	delete(s.cache, roleID)
	delete(s.cachedAt, roleID)
}

func (s *permissionService) startCleaner() {
	ticker := time.NewTicker(s.cacheTTL)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.clearExpired()
		case <-s.stopCh:
			return
		}
	}
}

func (s *permissionService) clearExpired() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	now := time.Now()
	for id, at := range s.cachedAt {
		if now.Sub(at) >= s.cacheTTL {
			delete(s.cache, id)
			delete(s.cachedAt, id)
		}
	}
}

// Shutdown вызывается при остановке приложения
func (s *permissionService) Shutdown() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
}
