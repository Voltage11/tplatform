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

	sbCount.Select("COUNT(*)").From("departments").
		Where(sbCount.IsNull("deleted_at"))

	if filter.Name != "" {
		pattern := "%" + filter.Name + "%"
		sbFilter.Where(sbFilter.ILike("name", pattern))
		sbCount.Where(sbCount.ILike("name", pattern))
	}

	sbFilter.Limit(filter.Pagination.GetLimit()).
		Offset(filter.Pagination.GetOffset()).
		OrderBy("name")

	return getList(ctx, d.pool, sbFilter, sbCount, func(scanner rowScanner) (*domain.Department, error) {
		var dep domain.Department
		err := scanner.Scan(&dep.ID, &dep.Name, &dep.CreatedAt, &dep.UpdatedAt, &dep.DeletedAt)
		return &dep, err
	})
}
