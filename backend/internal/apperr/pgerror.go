package apperr

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleDBError(err error, op string, entityName string) error {
	if err == nil {
		return nil
	}

	// Если уже AppError, просто пробрасываем (или оборачиваем Op)
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return NotFound(entityName+" не найден(а)", op)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return New(409, "CONFLICT", entityName+" уже существует", op, err)
		case "23503":
			return BadRequest(err, "Нарушена связь с другой записью", op)
		case "23502": // not_null_violation
			return BadRequest(err, "Обязательное поле не заполнено", op)
		case "23514": // check_violation
			return BadRequest(err, "Некорректное значение поля", op)
		case "25P02": // in_failed_sql_transaction
			return BadRequest(err, "Ошибка в транзакции", op)
		case "40001": // serialization_failure
			return BadRequest(err, "Конфликт параллельных операций", op)
		default:
			Internal(err, op)
		}
	}

	return Internal(err, op)
}
