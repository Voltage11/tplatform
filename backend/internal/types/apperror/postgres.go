package apperror

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// NewPostgresError конвертирует ошибку pgx в кастомную AppError
func NewPostgresError(err error) error {
	if err == nil {
		return nil
	}

	// Проверка на отсутствия строк в pgx
	if errors.Is(err, pgx.ErrNoRows) {
		return NewNotFound("запись не найдена", err)
	}

	// Проверка на ошибки самого PostgreSQL (нарушение ограничений и т.д.)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		// Дубликат ключа (например, email уже занят)
		case pgerrcode.UniqueViolation: // Код "23505"
			return NewConflict("такая запись уже существует", err)

		// Нарушение внешнего ключа (попытка привязать к несуществующей сущности)
		case pgerrcode.ForeignKeyViolation: // Код "23250"
			return NewBadRequest("нарушение связности данных (внешний ключ)", err)

		// Нарушение NOT NULL или CHECK ограничений
		case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
			return NewBadRequest("переданы некорректные данные", err)

			// Превышение таймаута выполнения запроса в БД
			// case pgerrcode.QueryCanceled:
			// 	return NewTimeout("превышено время ожидания базы данных", err)
			// }
		}
	}

	// 3. Дефолтный ответ для всех остальных непредвиденных ошибок БД
	return NewInternal("внутренняя ошибка базы данных", err)
}
