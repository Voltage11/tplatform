package domain

import (
	"context"

	"github.com/google/uuid"
)

// PermissionAction — действие
type PermissionAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func (a PermissionAction) WithActive(isActive bool) PermissionAction {
	return PermissionAction{
		Name:        a.Name,
		Description: a.Description,
		IsActive:    isActive,
	}
}

var (
	ActionView       = PermissionAction{Name: "view", Description: "Просмотр"}
	ActionCreate     = PermissionAction{Name: "create", Description: "Создание"}
	ActionUpdate     = PermissionAction{Name: "update", Description: "Обновление"}
	ActionSoftDelete = PermissionAction{Name: "soft_delete", Description: "Пометка на удаление"}
	ActionHardDelete = PermissionAction{Name: "hard_delete", Description: "Удаление"}
)

type PermissionEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	EntityDepartments = PermissionEntity{Name: "departments", Description: "Отделы"}
	EntityUsers       = PermissionEntity{Name: "users", Description: "Пользователи"}
	EntityRoles       = PermissionEntity{Name: "roles", Description: "Роли"}
	EntityThemes      = PermissionEntity{Name: "themes", Description: "Темы"}
)

type PermissionEntityWithActions struct {
	Entity  PermissionEntity   `json:"entity"`
	Actions []PermissionAction `json:"actions"`
}

func GetPermissionsEntitiesWithActions() []PermissionEntityWithActions {
	return []PermissionEntityWithActions{
		{
			Entity:  EntityDepartments,
			Actions: []PermissionAction{ActionView, ActionCreate, ActionUpdate, ActionSoftDelete},
		},
		{
			Entity:  EntityUsers,
			Actions: []PermissionAction{ActionView, ActionCreate, ActionUpdate, ActionSoftDelete, ActionHardDelete},
		},
		{
			Entity:  EntityRoles,
			Actions: []PermissionAction{ActionView, ActionCreate, ActionUpdate, ActionHardDelete},
		},
	}
}

// RolePermission — связка роли и конкретного права (БД)
type RolePermission struct {
	RoleID     uuid.UUID
	EntityName string
	ActionName string
}

// PermissionTarget — доменная структура для указания желаемых прав
type PermissionTarget struct {
	EntityName string
	ActionName string
}

// PermissionsRepository работа с таблицей role_permissions
type PermissionsRepository interface {
	SetAction(ctx context.Context, roleID uuid.UUID, entityName, actionName string) error
	RemoveAction(ctx context.Context, roleID uuid.UUID, entityName, actionName string) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]RolePermission, error)
	ClearRolePermissions(ctx context.Context, roleID uuid.UUID) error
}

// PermissionService сервис управления правами ролей
type PermissionService interface {
	SetAction(ctx context.Context, roleID uuid.UUID, entity PermissionEntity, actionName string) error
	RemoveAction(ctx context.Context, roleID uuid.UUID, entity PermissionEntity, actionName string) error
	GetForRole(ctx context.Context, roleID uuid.UUID) ([]PermissionEntityWithActions, error)
	Can(ctx context.Context, user *User, entityName, actionName string) bool
	CanFromCtx(ctx context.Context, entityName, actionName string) bool
	ReplacePermissions(ctx context.Context, roleID uuid.UUID, targets []PermissionTarget) error
	Shutdown()
}
