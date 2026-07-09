package handler

import (
    "net/http"
    "time"

    "github.com/Voltage11/tplatform/internal/domain"
    "github.com/Voltage11/tplatform/internal/handler/dto"
    "github.com/Voltage11/tplatform/internal/handler/httputils"
    "github.com/Voltage11/tplatform/internal/middleware"
    "github.com/go-chi/chi/v5"
    "github.com/go-playground/validator/v10"
)

type themeHandlers struct {
    themeService domain.ThemeService
    validate     *validator.Validate
}

func NewThemeHandlers(r chi.Router, authMW *middleware.AuthMiddleware, themeService domain.ThemeService) {
    h := themeHandlers{
        themeService: themeService,
        validate:     validator.New(),
    }

    r.Group(func(r chi.Router) {
        r.Use(authMW.RequireAuth) // все действия требуют авторизацию
        r.Get("/api/v1/themes", h.GetList)
        r.Post("/api/v1/themes", h.Create)
        r.Put("/api/v1/themes/{id}", h.Update)
        r.Delete("/api/v1/themes/{id}", h.Delete)
    })
}

func (h *themeHandlers) GetList(w http.ResponseWriter, r *http.Request) {
    pagination := httputils.ParsePagination(r)
    name := httputils.GetQueryValue(r, "name")

    var createdAtFrom, createdAtTo *time.Time
    if fromStr, ok := httputils.GetQueryValueWithExist(r, "created_at_from"); ok && fromStr != "" {
        t, err := time.Parse(time.RFC3339, fromStr)
        if err == nil {
            createdAtFrom = &t
        }
    }
    if toStr, ok := httputils.GetQueryValueWithExist(r, "created_at_to"); ok && toStr != "" {
        t, err := time.Parse(time.RFC3339, toStr)
        if err == nil {
            createdAtTo = &t
        }
    }

    isActive := httputils.ParseFilterBool(r, "is_active")
    createdByID := httputils.ParseUUIDQuery(r, "created_by_id")

    filter := domain.ThemeFilter{
        Name:          name,
        IsActive:      isActive,
        CreatedByID:   createdByID,
        CreatedAtFrom: createdAtFrom,
        CreatedAtTo:   createdAtTo,
        Pagination:    pagination,
    }

    themes, err := h.themeService.GetList(r.Context(), filter)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    themesDTO := dto.ThemesToResponseSlice(themes.Data)
    response := dto.NewPagedResponse(themesDTO, themes.Pagination)
    httputils.WriteOk(w, response)
}

func (h *themeHandlers) Create(w http.ResponseWriter, r *http.Request) {
    req, err := httputils.DecodeJSONBodyWithValidate[dto.ThemeCreateRequest](r, h.validate)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    theme := domain.Theme{
        Name:        req.Name,
        Description: req.Description,
        IsActive:    req.IsActive,
        DateBegin:   req.DateBegin,
        DateEnd:     req.DateEnd,
        CheckPoint:  req.CheckPoint,
        ImgPath:     req.ImgPath,
    }

    if err := h.themeService.Create(r.Context(), &theme); err != nil {
        httputils.WriteError(w, err)
        return
    }

    themeWithDetail, err := h.themeService.GetByIDWithDetail(r.Context(), theme.ID)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, dto.ThemeToResponse(themeWithDetail))
}

func (h *themeHandlers) Update(w http.ResponseWriter, r *http.Request) {
    req, err := httputils.DecodeJSONBodyWithValidate[dto.ThemeUpdateRequest](r, h.validate)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    id, err := httputils.ParseUUID(r, "id")
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    theme, err := h.themeService.GetByID(r.Context(), id)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    theme.Name = req.Name
    theme.Description = req.Description
    theme.IsActive = req.IsActive
    theme.DateBegin = req.DateBegin
    theme.DateEnd = req.DateEnd
    theme.CheckPoint = req.CheckPoint
    theme.ImgPath = req.ImgPath

    if err := h.themeService.Update(r.Context(), theme); err != nil {
        httputils.WriteError(w, err)
        return
    }

    themeWithDetail, err := h.themeService.GetByIDWithDetail(r.Context(), id)
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, dto.ThemeToResponse(themeWithDetail))
}

func (h *themeHandlers) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := httputils.ParseUUID(r, "id")
    if err != nil {
        httputils.WriteError(w, err)
        return
    }

    if err := h.themeService.Delete(r.Context(), id); err != nil {
        httputils.WriteError(w, err)
        return
    }

    httputils.WriteOk(w, map[string]string{"message": "Тема удалена"})
}