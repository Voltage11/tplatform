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

type themeRepository struct {
    db *db.PostgresDB
}

func NewThemeRepository(db *db.PostgresDB) domain.ThemeRepository {
    return &themeRepository{db: db}
}

func (t *themeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
    query := `SELECT id, name, description, is_active, created_by_id, created_at,
        date_begin, date_end, max_point, check_point, img_path
        FROM themes WHERE id = $1`

    var theme domain.Theme
    exec := t.db.GetDB(ctx)
    err := exec.QueryRow(ctx, query, id).Scan(
        &theme.ID, &theme.Name, &theme.Description, &theme.IsActive,
        &theme.CreatedByID, &theme.CreatedAt, &theme.DateBegin, &theme.DateEnd,
        &theme.MaxPoint, &theme.CheckPoint, &theme.ImgPath,
    )
    if err != nil {
        return nil, apperror.NewPostgresError(err)
    }
    return &theme, nil
}

func (t *themeRepository) GetByIDWithDetail(ctx context.Context, id uuid.UUID) (*domain.ThemeWithDetail, error) {
    query := `SELECT t.id, t.name, t.description, t.is_active, t.created_by_id, t.created_at,
        t.date_begin, t.date_end, t.max_point, t.check_point, t.img_path,
        u.id AS user_id, CONCAT_WS(' ', u.first_name, u.last_name) AS user_name
        FROM themes t JOIN users u ON t.created_by_id = u.id
        WHERE t.id = $1`

    var theme domain.ThemeWithDetail
    exec := t.db.GetDB(ctx)
    err := exec.QueryRow(ctx, query, id).Scan(
        &theme.ID, &theme.Name, &theme.Description, &theme.IsActive,
        &theme.CreatedByID, &theme.CreatedAt, &theme.DateBegin, &theme.DateEnd,
        &theme.MaxPoint, &theme.CheckPoint, &theme.ImgPath,
        &theme.CreatedBy.ID, &theme.CreatedBy.Name,
    )
    if err != nil {
        return nil, apperror.NewPostgresError(err)
    }
    return &theme, nil
}

func (t *themeRepository) Create(ctx context.Context, theme *domain.Theme) error {
    query := `INSERT INTO themes(name, description, is_active, created_by_id, created_at,
        date_begin, date_end, check_point, img_path)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id`

    theme.CreatedAt = time.Now().UTC()

    exec := t.db.GetDB(ctx)
    err := exec.QueryRow(ctx, query,
        theme.Name, theme.Description, theme.IsActive, theme.CreatedByID,
        theme.CreatedAt, theme.DateBegin, theme.DateEnd, theme.CheckPoint, theme.ImgPath,
    ).Scan(&theme.ID)
    return apperror.NewPostgresError(err)
}

func (t *themeRepository) Update(ctx context.Context, theme *domain.Theme) error {
    query := `UPDATE themes SET name = $1, description = $2, is_active = $3,
        date_begin = $4, date_end = $5, check_point = $6, img_path = $7
        WHERE id = $8`

    exec := t.db.GetDB(ctx)
    result, err := exec.Exec(ctx, query,
        theme.Name, theme.Description, theme.IsActive,
        theme.DateBegin, theme.DateEnd, theme.CheckPoint, theme.ImgPath, theme.ID,
    )
    if err != nil {
        return apperror.NewPostgresError(err)
    }
    if result.RowsAffected() == 0 {
        return apperror.NewNotFound("Тема не найдена", nil)
    }
    return nil
}

func (t *themeRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM themes WHERE id = $1`
    exec := t.db.GetDB(ctx)
    result, err := exec.Exec(ctx, query, id)
    if err != nil {
        return apperror.NewPostgresError(err)
    }
    if result.RowsAffected() == 0 {
        return apperror.NewNotFound("Тема не найдена", nil)
    }
    return nil
}

func (t *themeRepository) SetActive(ctx context.Context, id uuid.UUID, isActive bool) error {
    query := `UPDATE themes SET is_active = $1 WHERE id = $2`
    exec := t.db.GetDB(ctx)
    result, err := exec.Exec(ctx, query, isActive, id)
    if err != nil {
        return apperror.NewPostgresError(err)
    }
    if result.RowsAffected() == 0 {
        return apperror.NewNotFound("Тема не найдена", nil)
    }
    return nil
}

func (t *themeRepository) GetList(ctx context.Context, filter domain.ThemeFilter) ([]*domain.ThemeWithDetail, int64, error) {
    sbFilter := sqlbuilder.PostgreSQL.NewSelectBuilder()
    sbCount := sqlbuilder.PostgreSQL.NewSelectBuilder()

    sbFilter.Select(
        "t.id", "t.name", "t.description", "t.is_active", "t.created_by_id",
        "t.created_at", "t.date_begin", "t.date_end", "t.max_point",
        "t.check_point", "t.img_path",
        "u.id AS user_id", "CONCAT_WS(' ', u.first_name, u.last_name) AS user_name",
    ).
        From("themes t").
        JoinWithOption(sqlbuilder.LeftJoin, "users u", "t.created_by_id = u.id").
        Where(sbFilter.IsNull("t.deleted_at"))

    sbCount.Select("COUNT(*)").From("themes t").
        Where(sbCount.IsNull("t.deleted_at"))

    if filter.Name != "" {
        pattern := "%" + filter.Name + "%"
        sbFilter.Where(sbFilter.ILike("t.name", pattern))
        sbCount.Where(sbCount.ILike("t.name", pattern))
    }
    if filter.CreatedByID != nil {
        sbFilter.Where(sbFilter.Equal("t.created_by_id", *filter.CreatedByID))
        sbCount.Where(sbCount.Equal("t.created_by_id", *filter.CreatedByID))
    }
    if isActive := filter.IsActive.GetBool(); isActive != nil {
        sbFilter.Where(sbFilter.Equal("t.is_active", *isActive))
        sbCount.Where(sbCount.Equal("t.is_active", *isActive))
    }
    if filter.CreatedAtFrom != nil {
        sbFilter.Where(sbFilter.GTE("t.created_at", filter.CreatedAtFrom))
        sbCount.Where(sbCount.GTE("t.created_at", filter.CreatedAtFrom))
    }
    if filter.CreatedAtTo != nil {
        sbFilter.Where(sbFilter.LTE("t.created_at", filter.CreatedAtTo))
        sbCount.Where(sbCount.LTE("t.created_at", filter.CreatedAtTo))
    }

    sbFilter.Limit(filter.Pagination.GetLimit()).
        Offset(filter.Pagination.GetOffset()).
        OrderBy("t.created_at")

    return getList(ctx, t.db, sbFilter, sbCount, func(scanner rowScanner) (*domain.ThemeWithDetail, error) {
        var theme domain.ThemeWithDetail
        err := scanner.Scan(
            &theme.ID, &theme.Name, &theme.Description, &theme.IsActive,
            &theme.CreatedByID, &theme.CreatedAt, &theme.DateBegin, &theme.DateEnd,
            &theme.MaxPoint, &theme.CheckPoint, &theme.ImgPath,
            &theme.CreatedBy.ID, &theme.CreatedBy.Name,
        )
        if err != nil {
            return nil, apperror.NewPostgresError(err)
        }
        return &theme, nil
    })
}