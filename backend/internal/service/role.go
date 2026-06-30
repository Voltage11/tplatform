package service

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/google/uuid"
)

type roleService struct {
	repo domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) domain.RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) Create(ctx context.Context, role *domain.Role) error {
	return s.repo.Create(ctx, role)
}

func (s *roleService) Update(ctx context.Context, role *domain.Role) error {
	return s.repo.Update(ctx, role)
}

func (s *roleService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *roleService) GetList(ctx context.Context, filter domain.RoleFilter) (*domain.PagedResult[*domain.Role], error) {
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
	return s.repo.Delete(ctx, id)
}
