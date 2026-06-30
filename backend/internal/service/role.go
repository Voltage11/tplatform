package service

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
)

type roleService struct {
	repo              domain.RoleRepository
	permissionService domain.PermissionService
}

func NewRoleService(repo domain.RoleRepository, permissionService domain.PermissionService) domain.RoleService {
	return &roleService{
		repo:              repo,
		permissionService: permissionService,
	}
}

func (s *roleService) Create(ctx context.Context, role *domain.Role) error {
	if !s.permissionService.CanFromCtx(ctx, domain.EntityRoles.Name, domain.ActionCreate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return s.repo.Create(ctx, role)
}

func (s *roleService) Update(ctx context.Context, role *domain.Role) error {
	if !s.permissionService.CanFromCtx(ctx, domain.EntityRoles.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return s.repo.Update(ctx, role)
}

func (s *roleService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	if !s.permissionService.CanFromCtx(ctx, domain.EntityRoles.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return s.repo.GetByID(ctx, id)
}

func (s *roleService) GetList(ctx context.Context, filter domain.RoleFilter) (*domain.PagedResult[*domain.Role], error) {
	if !s.permissionService.CanFromCtx(ctx, domain.EntityRoles.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	departments, total, err := s.repo.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	pagination := domain.NewPaginationResponse(filter.Pagination, total)

	return &domain.PagedResult[*domain.Role]{
		Data:       departments,
		Pagination: pagination,
	}, nil
}

func (s *roleService) Delete(ctx context.Context, id uuid.UUID) error {
	if !s.permissionService.CanFromCtx(ctx, domain.EntityRoles.Name, domain.ActionHardDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return s.repo.Delete(ctx, id)
}
