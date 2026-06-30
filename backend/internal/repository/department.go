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

type departmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) domain.DepartmentRepository {
	return &departmentRepository{pool: pool}
}

func (d *departmentRepository) Create(ctx context.Context, department *domain.Department) error {
	now := time.Now().UTC()
	department.CreatedAt = now
	department.UpdatedAt = now

	query := `INSERT INTO departments (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`

	if err := d.pool.QueryRow(ctx, query, department.Name, department.CreatedAt, department.UpdatedAt).Scan(&department.ID); err != nil {
		return apperror.NewPostgresError(err)
	}
	return nil
}

func (d *departmentRepository) Update(ctx context.Context, department *domain.Department) error {
	department.UpdatedAt = time.Now().UTC()

	query := `UPDATE departments SET name = $1, updated_at = $2 WHERE id = $3`

	result, err := d.pool.Exec(ctx, query, department.Name, department.UpdatedAt, department.ID)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Отдел не найден", nil)
	}
	return nil
}

func (d *departmentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	deletedAt := time.Now().UTC()
	query := `UPDATE departments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := d.pool.Exec(ctx, query, deletedAt, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Отдел не найден", nil)
	}
	return nil
}

func (d *departmentRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM departments WHERE id = $1`

	result, err := d.pool.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Отдел не найден", nil)
	}
	return nil
}

// GetByID в слое сервиса нужно проверить кто запрашивает, если запись мягко удалена и не админ, то отказать
func (d *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	query := `SELECT id, name, created_at, updated_at, deleted_at FROM departments WHERE id = $1`

	department := &domain.Department{}
	if err := d.pool.QueryRow(ctx, query, id).Scan(
		&department.ID, &department.Name, &department.CreatedAt, &department.UpdatedAt, &department.DeletedAt,
	); err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	return department, nil
}

func (d *departmentRepository) GetList(ctx context.Context, filter domain.DepartmentFilter) ([]*domain.Department, int64, error) {
	sbFilter := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sbCount := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sbFilter.Select("id", "name", "created_at", "updated_at", "deleted_at").
		From("departments").
		Where(sbFilter.IsNull("deleted_at"))

	sbCount.Select("COUNT(*)").
		From("departments").
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

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}
	defer tx.Rollback(ctx)

	var total int64
	if err := tx.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	if total == 0 {
		return []*domain.Department{}, 0, nil
	}

	rows, err := tx.Query(ctx, queryFilter, argsFilter...)
	if err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}
	defer rows.Close()

	departments := make([]*domain.Department, 0)
	for rows.Next() {
		var dep domain.Department
		if err := rows.Scan(&dep.ID, &dep.Name, &dep.CreatedAt, &dep.UpdatedAt, &dep.DeletedAt); err != nil {
			return nil, 0, apperror.NewPostgresError(err)
		}
		departments = append(departments, &dep)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	return departments, total, nil
}
