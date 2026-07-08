package repository

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type roleRepository struct {
	db *db.PostgresDB
}

func NewRoleRepository(db *db.PostgresDB) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	role.CreatedAt = time.Now().UTC()

	query := `INSERT INTO roles (name, description, created_at) VALUES ($1, $2, $3) RETURNING id`

	executor := r.db.GetDB(ctx)
	err := executor.QueryRow(ctx, query, role.Name, role.Description, role.CreatedAt).Scan(&role.ID)

	return apperror.NewPostgresError(err)
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`

	executor := r.db.GetDB(ctx)
	result, err := executor.Exec(ctx, query, role.Name, role.Description, role.ID)
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
	
	executor := r.db.GetDB(ctx)
	err := executor.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`

	executor := r.db.GetDB(ctx)
	result, err := executor.Exec(ctx, query, id)
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
		From("roles")
		// Where(sbFilter.IsNull("deleted_at"))

	sbCount.Select("COUNT(*)").From("roles")
		// Where(sbCount.IsNull("deleted_at"))

	if filter.Name != "" {
		pattern := "%" + filter.Name + "%"
		sbFilter.Where(sbFilter.ILike("name", pattern))
		sbCount.Where(sbCount.ILike("name", pattern))
	}

	sbFilter.Limit(filter.Pagination.GetLimit()).
		Offset(filter.Pagination.GetOffset()).
		OrderBy("name")

	return getList(ctx, r.db, sbFilter, sbCount, func(scanner rowScanner) (*domain.Role, error) {
		var role domain.Role
		err := scanner.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		return &role, err
	})
}
