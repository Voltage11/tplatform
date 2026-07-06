package dto

// UpdatePermissionsRequest – тело запроса PUT /api/v1/roles/{id}/permissions
type UpdatePermissionsRequest struct {
	Permissions []PermissionItem `json:"permissions" validate:"required,dive"`
}

// PermissionItem – одна пара сущность‑действие
type PermissionItem struct {
	EntityName string `json:"entity_name" validate:"required"`
	ActionName string `json:"action_name" validate:"required"`
}