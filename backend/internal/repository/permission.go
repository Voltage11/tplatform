package repository

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type permissionsRepository struct {
	pool *pgxpool.Pool
}

func NewPermissionsRepository(pool *pgxpool.Pool) domain.PermissionsRepository {
	return &permissionsRepository{pool: pool}
}

func (r *permissionsRepository) SetAction(ctx context.Context, roleID uuid.UUID, entityName, actionName string) error {
	query := `
        INSERT INTO role_permissions (role_id, entity_name, action_name)
        VALUES ($1, $2, $3)
        ON CONFLICT DO NOTHING
    `
	_, err := r.pool.Exec(ctx, query, roleID, entityName, actionName)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	return nil
}

func (r *permissionsRepository) RemoveAction(ctx context.Context, roleID uuid.UUID, entityName, actionName string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND entity_name = $2 AND action_name = $3`
	result, err := r.pool.Exec(ctx, query, roleID, entityName, actionName)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Право не найдено", nil)
	}
	return nil
}

func (r *permissionsRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.RolePermission, error) {
	query := `SELECT role_id, entity_name, action_name FROM role_permissions WHERE role_id = $1`
	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	defer rows.Close()

	var perms []domain.RolePermission
	for rows.Next() {
		var p domain.RolePermission
		if err := rows.Scan(&p.RoleID, &p.EntityName, &p.ActionName); err != nil {
			return nil, apperror.NewPostgresError(err)
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}
