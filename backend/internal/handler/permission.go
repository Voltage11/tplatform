package handler

import (
	"net/http"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/dto"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type permissionHandler struct {
	permissionService domain.PermissionService
	validate          *validator.Validate
}

func NewPermissionHandler(r chi.Router, authMW *middleware.AuthMiddleware, permissionService domain.PermissionService) {
	h := &permissionHandler{
		permissionService: permissionService,
		validate:          validator.New(),
	}

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth, authMW.RequireAdmin)
		r.Get("/api/v1/roles/{id}/permissions", h.GetPermissions)
		r.Put("/api/v1/roles/{id}/permissions", h.UpdatePermissions)
	})
}

func (h *permissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	perms, err := h.permissionService.GetForRole(r.Context(), roleID)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, perms)
}

func (h *permissionHandler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	roleID, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	req, err := httputils.DecodeJSONBodyWithValidate[dto.UpdatePermissionsRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Преобразуем DTO в доменные объекты
	targets := make([]domain.PermissionTarget, len(req.Permissions))
	for i, p := range req.Permissions {
		targets[i] = domain.PermissionTarget{
			EntityName: p.EntityName,
			ActionName: p.ActionName,
		}
	}

	if err := h.permissionService.ReplacePermissions(r.Context(), roleID, targets); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "права обновлены"})
}