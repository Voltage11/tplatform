package service

import (
    "context"

    "github.com/Voltage11/tplatform/internal/appcontext"
    "github.com/Voltage11/tplatform/internal/domain"
    "github.com/Voltage11/tplatform/internal/types/apperror"
    "github.com/google/uuid"
)

type themeService struct {
    themeRepo   domain.ThemeRepository
    permissions domain.PermissionService
}

func NewThemeService(themeRepo domain.ThemeRepository, permissions domain.PermissionService) domain.ThemeService {
    return &themeService{
        themeRepo:   themeRepo,
        permissions: permissions,
    }
}

func (t *themeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionView.Name) {
        return nil, apperror.NewForbiddenWithoutErr()
    }
    return t.themeRepo.GetByID(ctx, id)
}

func (t *themeService) GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*domain.ThemeWithDetail, error) {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionView.Name) {
        return nil, apperror.NewForbiddenWithoutErr()
    }
    return t.themeRepo.GetByIDWithDetail(ctx, id)
}

func (t *themeService) Create(ctx context.Context, theme *domain.Theme) error {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionCreate.Name) {
        return apperror.NewForbiddenWithoutErr()
    }

    user := appcontext.GetUserFromContext(ctx)
    if user == nil {
        return apperror.NewUnauthorized("пользователь не аутентифицирован", nil)
    }
    theme.CreatedByID = user.ID
    return t.themeRepo.Create(ctx, theme)
}

func (t *themeService) Update(ctx context.Context, theme *domain.Theme) error {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionUpdate.Name) {
        return apperror.NewForbiddenWithoutErr()
    }
    return t.themeRepo.Update(ctx, theme)
}

func (t *themeService) Delete(ctx context.Context, id uuid.UUID) error {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionSoftDelete.Name) {
        return apperror.NewForbiddenWithoutErr()
    }
    return t.themeRepo.Delete(ctx, id)
}

func (t *themeService) SetActive(ctx context.Context, id uuid.UUID, isActive bool) error {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionUpdate.Name) {
        return apperror.NewForbiddenWithoutErr()
    }
    return t.themeRepo.SetActive(ctx, id, isActive)
}

func (t *themeService) GetList(ctx context.Context, filter domain.ThemeFilter) (*domain.PagedResult[*domain.ThemeWithDetail], error) {
    if !t.permissions.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionView.Name) {
        return nil, apperror.NewForbiddenWithoutErr()
    }

    themes, total, err := t.themeRepo.GetList(ctx, filter)
    if err != nil {
        return nil, err
    }

    pagination := domain.NewPaginationResponse(filter.Pagination, total)
    return &domain.PagedResult[*domain.ThemeWithDetail]{
        Data:       themes,
        Pagination: pagination,
    }, nil
}