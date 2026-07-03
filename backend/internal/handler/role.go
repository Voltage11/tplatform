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

type roleHandler struct {
    roleService domain.RoleService
    validate    *validator.Validate
}

func NewRoleHandler(r chi.Router, authMW *middleware.AuthMiddleware, roleService domain.RoleService) {
    h := roleHandler{
        roleService: roleService,
        validate:    validator.New(),
    }

    r.Group(func(r chi.Router) {
        r.Use(authMW.RequireAuth)
        r.Get("/api/v1/roles", h.GetList)
        r.Post("/api/v1/roles", h.Create)
        r.Put("/api/v1/roles/{id}", h.Update)
        r.Delete("/api/v1/roles/{id}", h.Delete)
    })
}

func (h *roleHandler) Create(w http.ResponseWriter, r *http.Request) {
    req, err := httputils.DecodeJSONBodyWithValidate[dto.RoleCreateRequest](r, h.validate)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    role := domain.Role{
        Name:        req.Name,
        Description: req.Description,
    }

    if err := h.roleService.Create(r.Context(), &role); err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, dto.RoleToResponse(&role))
}

func (h *roleHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := httputils.ParseUUID(r, "id")
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    req, err := httputils.DecodeJSONBodyWithValidate[dto.RoleUpdateRequest](r, h.validate)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    role, err := h.roleService.GetByID(r.Context(), id)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    role.Name = req.Name
    role.Description = req.Description

    if err := h.roleService.Update(r.Context(), role); err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, dto.RoleToResponse(role))
}

func (h *roleHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := httputils.ParseUUID(r, "id")
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    if err := h.roleService.Delete(r.Context(), id); err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, map[string]string{"message": "Роль удалена"})
}

func (h *roleHandler) GetList(w http.ResponseWriter, r *http.Request) {
    paginationRequest := httputils.ParsePagination(r)

    // Безопасное получение фильтра по имени
    nameFilter, ok := httputils.GetQueryValue(r, "name")
    if ok && nameFilter == "" {
        nameFilter = "" // игнорируем пустой параметр
    } else if !ok {
        nameFilter = ""
    }

    roleFilter := domain.RoleFilter{
        Name:       nameFilter,
        Pagination: paginationRequest,
    }

    rolesDomain, err := h.roleService.GetList(r.Context(), roleFilter)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    // Преобразуем в DTO
    roleDTOs := dto.RolesToResponseSlice(rolesDomain.Data)
    response := dto.NewPagedResponse(roleDTOs, rolesDomain.Pagination)

    httputils.WriteOk(w, response)
}