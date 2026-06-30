package service

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
)

type departmentService struct {
	repo              domain.DepartmentRepository
	permissionService domain.PermissionService
}

func NewDepartmentService(repo domain.DepartmentRepository, permissionService domain.PermissionService) domain.DepartmentService {
	return &departmentService{
		repo:              repo,
		permissionService: permissionService,
	}
}

func (d *departmentService) Create(ctx context.Context, department *domain.Department) error {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionCreate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return d.repo.Create(ctx, department)
}

func (d *departmentService) Update(ctx context.Context, department *domain.Department) error {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return d.repo.Update(ctx, department)
}

func (d *departmentService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionSoftDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return d.repo.SoftDelete(ctx, id)
}

func (d *departmentService) HardDelete(ctx context.Context, id uuid.UUID) error {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionHardDelete.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return d.repo.HardDelete(ctx, id)
}

func (d *departmentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return d.repo.GetByID(ctx, id)
}

func (d *departmentService) GetList(ctx context.Context, filter domain.DepartmentFilter) (*domain.PagedResult[*domain.Department], error) {
	if !d.permissionService.CanFromCtx(ctx, domain.EntityDepartments.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	departments, total, err := d.repo.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	pagination := domain.NewPaginationResponse(filter.Pagination, total)

	return &domain.PagedResult[*domain.Department]{
		Data:       departments,
		Pagination: pagination,
	}, nil
}
