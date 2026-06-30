package repository

import (
    "context"
    "github.com/huandu/go-sqlbuilder"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Voltage11/tplatform/internal/types/apperror"
)

type rowScanner interface {
    Scan(dest ...any) error
}

func getList[T any](
    ctx context.Context,
    pool *pgxpool.Pool,
    sbFilter, sbCount *sqlbuilder.SelectBuilder,
    scanRow func(scanner rowScanner) (T, error),
) ([]T, int64, error) {
    queryFilter, argsFilter := sbFilter.Build()
    queryCount, argsCount := sbCount.Build()

    tx, err := pool.Begin(ctx)
    if err != nil {
        return nil, 0, apperror.NewPostgresError(err)
    }
    defer tx.Rollback(ctx)

    var total int64
    if err := tx.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
        return nil, 0, apperror.NewPostgresError(err)
    }

    if total == 0 {
        return []T{}, 0, nil
    }

    rows, err := tx.Query(ctx, queryFilter, argsFilter...)
    if err != nil {
        return nil, 0, apperror.NewPostgresError(err)
    }
    defer rows.Close()

    result := make([]T, 0, total)
    for rows.Next() {
        item, err := scanRow(rows)
        if err != nil {
            return nil, 0, apperror.NewPostgresError(err)
        }
        result = append(result, item)
    }
    if err := rows.Err(); err != nil {
        return nil, 0, apperror.NewPostgresError(err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, 0, apperror.NewPostgresError(err)
    }

    return result, total, nil
}