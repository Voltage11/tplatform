package repository

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) domain.RoleRepository {
	return &roleRepository{pool: pool}
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	role.CreatedAt = time.Now().UTC()

	query := `INSERT INTO roles (name, description, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.pool.QueryRow(ctx, query, role.Name, role.Description, role.CreatedAt).Scan(&role.ID)

	return apperror.NewPostgresError(err)
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`

	result, err := r.pool.Exec(ctx, query, role.Name, role.Description, role.ID)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("роль не найдена", nil)
	}
	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`
	role := &domain.Role{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("роль не найдена", nil)
	}
	return nil
}

func (r *roleRepository) GetList(ctx context.Context, filter domain.RoleFilter) ([]*domain.Role, int64, error) {

	sbFilter := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sbCount := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sbFilter.Select("id", "name", "description", "created_at").
		From("roles").
		Where(sbFilter.IsNull("deleted_at"))

	sbCount.Select("COUNT(*)").
		From("roles").
		Where(sbCount.IsNull("deleted_at"))

	if filter.Name != "" {
		pattern := "%" + filter.Name + "%"
		sbFilter.Where(sbFilter.ILike("name", pattern))
		sbCount.Where(sbCount.ILike("name", pattern))
	}

	sbFilter.Limit(filter.Pagination.GetLimit()).
		Offset(filter.Pagination.GetOffset()).
		OrderBy("name")

	queryFilter, argsFilter := sbFilter.Build()
	queryCount, argsCount := sbCount.Build()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}
	defer tx.Rollback(ctx)

	var total int64
	if err := tx.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	if total == 0 {
		return []*domain.Role{}, 0, nil
	}

	rows, err := tx.Query(ctx, queryFilter, argsFilter...)
	if err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}
	defer rows.Close()

	roles := make([]*domain.Role, 0)
	for rows.Next() {
		var dep domain.Role
		if err := rows.Scan(&dep.ID, &dep.Name, &dep.Description, &dep.CreatedAt); err != nil {
			return nil, 0, apperror.NewPostgresError(err)
		}
		roles = append(roles, &dep)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	return roles, total, nil
}
