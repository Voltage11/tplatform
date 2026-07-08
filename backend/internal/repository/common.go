package repository

import (
	"context"

	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/huandu/go-sqlbuilder"
)

type rowScanner interface {
    Scan(dest ...any) error
}

func getList[T any](
	ctx context.Context,
	postgresDB *db.PostgresDB,
	sbFilter, sbCount *sqlbuilder.SelectBuilder,
	scanRow func(scanner rowScanner) (T, error),
) ([]T, int64, error) {
	queryFilter, argsFilter := sbFilter.Build()
	queryCount, argsCount := sbCount.Build()

	// 2. Автоматически получаем нужный исполнитель (транзакция или пул)
	executor := postgresDB.GetDB(ctx)

	// 3. Выполняем подсчет общего количества
	var total int64
	if err := executor.QueryRow(ctx, queryCount, argsCount...).Scan(&total); err != nil {
		return nil, 0, apperror.NewPostgresError(err)
	}

	if total == 0 {
		return []T{}, 0, nil
	}

	// 4. Выполняем фильтрованный запрос списка
	rows, err := executor.Query(ctx, queryFilter, argsFilter...)
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
	return result, total, nil
}