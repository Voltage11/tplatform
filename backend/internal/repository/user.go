package repository

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/pkg/helpers"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type userRepository struct {
	db *db.PostgresDB
}

func NewUserRepository(postgresDB *db.PostgresDB) domain.UserRepository {
	return &userRepository{db: postgresDB}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO users 
        (first_name, second_name, last_name, email, password_hash, department_id, role_id, is_active, is_admin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id`

	executor := u.db.GetDB(ctx)

	if err := executor.QueryRow(ctx, query,
		user.FirstName,
		user.SecondName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.DepartmentID,
		user.RoleID,
		user.IsActive,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		return apperror.NewPostgresError(err)
	}
	return nil
}

// Update обновляем любые поля, кроме password_hash
func (u *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC()

	query := `UPDATE users SET 
        first_name = $1, second_name = $2, last_name = $3, email = $4, 
        department_id = $5, role_id = $6, is_active = $7, is_admin = $8, updated_at = $9
        WHERE id = $10`

	executor := u.db.GetDB(ctx)

	result, err := executor.Exec(ctx, query,
		user.FirstName, user.SecondName, user.LastName, user.Email,
		user.DepartmentID,
		user.RoleID,
		user.IsActive,
		user.IsAdmin,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Пользователь не найден", nil)
	}
	return nil
}

func (u *userRepository) SetIsActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	executor := u.db.GetDB(ctx)

	result, err := executor.Exec(ctx, query, isActive, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Пользователь не найден", nil)
	}
	return nil
}

func (u *userRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	executor := u.db.GetDB(ctx)
	result, err := executor.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Пользователь не найден", nil)
	}
	return nil
}

func (u *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	deletedAt := time.Now().UTC()
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	executor := u.db.GetDB(ctx)
	result, err := executor.Exec(ctx, query, deletedAt, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Пользователь не найден", nil)
	}
	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, first_name, second_name, last_name, email, password_hash, 
        department_id, role_id, is_active, is_admin, created_at, updated_at, deleted_at
        FROM users WHERE id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var deptID, roleID uuid.NullUUID

	executor := u.db.GetDB(ctx)

	err := executor.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&deptID,
		&roleID,
		&user.IsActive,
		&user.IsAdmin, // <-- обязательно
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	user.DepartmentID = helpers.ScanNullableUUID(deptID)
	user.RoleID = helpers.ScanNullableUUID(roleID)
	return user, nil
}

func (u *userRepository) GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*domain.UserWithDetail, error) {
	query := `SELECT 
        u.id, u.first_name, u.second_name, u.last_name, u.email, u.password_hash,
        u.department_id, u.role_id, u.is_active, u.is_admin, 
        u.created_at, u.updated_at, u.deleted_at,
        d.id AS dept_id, d.name AS dept_name,
        r.id AS role_id_val, r.name AS role_name
    FROM users u
    LEFT JOIN departments d ON u.department_id = d.id
    LEFT JOIN roles r ON u.role_id = r.id
    WHERE u.id = $1 AND u.deleted_at IS NULL`

	user := &domain.UserWithDetail{}
	var deptID, roleID uuid.NullUUID
	var deptIDPtr, roleIDPtr *uuid.UUID
	var deptName, roleName *string

	executor := u.db.GetDB(ctx)

	err := executor.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&deptID,
		&roleID,
		&user.IsActive,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&deptIDPtr,
		&deptName,
		&roleIDPtr,
		&roleName,
	)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}

	user.DepartmentID = helpers.ScanNullableUUID(deptID)
	user.RoleID = helpers.ScanNullableUUID(roleID)

	if deptIDPtr != nil && deptName != nil {
		user.Department = &domain.DepartmentDetail{
			ID:   *deptIDPtr,
			Name: *deptName,
		}
	}
	if roleIDPtr != nil && roleName != nil {
		user.Role = &domain.RoleDetail{
			ID:   *roleIDPtr,
			Name: *roleName,
		}
	}

	return user, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, first_name, second_name, last_name, email, password_hash, 
        department_id, role_id, is_active, is_admin, created_at, updated_at, deleted_at
        FROM users WHERE email = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var deptID, roleID uuid.NullUUID

	executor := u.db.GetDB(ctx)

	err := executor.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Email,
		&user.PasswordHash, &deptID, &roleID, &user.IsActive, &user.IsAdmin,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}
	user.DepartmentID = helpers.ScanNullableUUID(deptID)
	user.RoleID = helpers.ScanNullableUUID(roleID)
	return user, nil
}

func (u *userRepository) GetByEmailWithDetail(ctx context.Context, email string) (*domain.UserWithDetail, error) {
	query := `SELECT 
        u.id, u.first_name, u.second_name, u.last_name, u.email, u.password_hash,
        u.department_id, u.role_id, u.is_active, u.is_admin, 
        u.created_at, u.updated_at, u.deleted_at,
        d.id AS dept_id, d.name AS dept_name,
        r.id AS role_id_val, r.name AS role_name
    FROM users u
    LEFT JOIN departments d ON u.department_id = d.id
    LEFT JOIN roles r ON u.role_id = r.id
    WHERE u.email = $1 AND u.deleted_at IS NULL`

	user := &domain.UserWithDetail{}
	var deptID, roleID uuid.NullUUID
	var deptIDPtr, roleIDPtr *uuid.UUID
	var deptName, roleName *string

	executor := u.db.GetDB(ctx)

	err := executor.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&deptID,
		&roleID,
		&user.IsActive,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&deptIDPtr,
		&deptName,
		&roleIDPtr,
		&roleName,
	)
	if err != nil {
		return nil, apperror.NewPostgresError(err)
	}

	user.DepartmentID = helpers.ScanNullableUUID(deptID)
	user.RoleID = helpers.ScanNullableUUID(roleID)

	if deptIDPtr != nil && deptName != nil {
		user.Department = &domain.DepartmentDetail{
			ID:   *deptIDPtr,
			Name: *deptName,
		}
	}
	if roleIDPtr != nil && roleName != nil {
		user.Role = &domain.RoleDetail{
			ID:   *roleIDPtr,
			Name: *roleName,
		}
	}

	return user, nil
}

func (u *userRepository) GetList(ctx context.Context, filter domain.UserFilter) ([]*domain.UserWithDetail, int64, error) {
	sbFilter := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sbCount := sqlbuilder.PostgreSQL.NewSelectBuilder()

	// Основной фильтр с LEFT JOIN
	sbFilter.Select(
		"u.id", "u.first_name", "u.second_name", "u.last_name",
		"u.email", "u.password_hash",
		"u.department_id", "u.role_id", "u.is_active", "u.is_admin",
		"u.created_at", "u.updated_at", "u.deleted_at",
		"d.id AS dept_id", "d.name AS dept_name",
		"r.id AS role_id_val", "r.name AS role_name",
	).
		From("users u").
		JoinWithOption(sqlbuilder.LeftJoin, "departments d", "u.department_id = d.id").
		JoinWithOption(sqlbuilder.LeftJoin, "roles r", "u.role_id = r.id").
		Where(sbFilter.IsNull("u.deleted_at"))

	sbCount.Select("COUNT(*)").From("users u").
		Where(sbCount.IsNull("u.deleted_at"))

	// Фильтры
	if filter.FirstName != "" {
		pattern := "%" + filter.FirstName + "%"
		sbFilter.Where(sbFilter.ILike("u.first_name", pattern))
		sbCount.Where(sbCount.ILike("u.first_name", pattern))
	}
	if filter.SecondName != "" {
		pattern := "%" + filter.SecondName + "%"
		sbFilter.Where(sbFilter.ILike("u.second_name", pattern))
		sbCount.Where(sbCount.ILike("u.second_name", pattern))
	}
	if filter.LastName != "" {
		pattern := "%" + filter.LastName + "%"
		sbFilter.Where(sbFilter.ILike("u.last_name", pattern))
		sbCount.Where(sbCount.ILike("u.last_name", pattern))
	}
	if filter.Email != "" {
		pattern := "%" + filter.Email + "%"
		sbFilter.Where(sbFilter.ILike("u.email", pattern))
		sbCount.Where(sbCount.ILike("u.email", pattern))
	}
	if filter.DepartmentID != nil {
		sbFilter.Where(sbFilter.Equal("u.department_id", *filter.DepartmentID))
		sbCount.Where(sbCount.Equal("u.department_id", *filter.DepartmentID))
	}
	if filter.RoleID != nil {
		sbFilter.Where(sbFilter.Equal("u.role_id", *filter.RoleID))
		sbCount.Where(sbCount.Equal("u.role_id", *filter.RoleID))
	}
	if isActive := filter.IsActive.GetBool(); isActive != nil {
		sbFilter.Where(sbFilter.Equal("u.is_active", *isActive))
		sbCount.Where(sbCount.Equal("u.is_active", *isActive))
	}

	sbFilter.Limit(filter.Pagination.GetLimit()).
		Offset(filter.Pagination.GetOffset()).
		OrderBy("u.last_name")

	// Используем универсальный getList с готовой scanUserRow
	return getList(ctx, u.db, sbFilter, sbCount, scanUserRow)
}

func scanUserRow(scanner rowScanner) (*domain.UserWithDetail, error) {
	usr := &domain.UserWithDetail{}
	var deptID, roleID uuid.NullUUID
	var deptIDPtr, roleIDPtr *uuid.UUID
	var deptName, roleName *string

	err := scanner.Scan(
		&usr.ID,
		&usr.FirstName,
		&usr.SecondName,
		&usr.LastName,
		&usr.Email,
		&usr.PasswordHash,
		&deptID,
		&roleID,
		&usr.IsActive,
		&usr.IsAdmin,
		&usr.CreatedAt,
		&usr.UpdatedAt,
		&usr.DeletedAt,
		&deptIDPtr,
		&deptName,
		&roleIDPtr,
		&roleName,
	)
	if err != nil {
		return nil, err
	}

	usr.DepartmentID = helpers.ScanNullableUUID(deptID)
	usr.RoleID = helpers.ScanNullableUUID(roleID)

	if deptIDPtr != nil && deptName != nil {
		usr.Department = &domain.DepartmentDetail{ID: *deptIDPtr, Name: *deptName}
	}
	if roleIDPtr != nil && roleName != nil {
		usr.Role = &domain.RoleDetail{ID: *roleIDPtr, Name: *roleName}
	}
	return usr, nil
}
