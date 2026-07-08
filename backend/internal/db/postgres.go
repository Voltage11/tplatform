package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBClient объединяет методы pgxpool.Pool и pgx.Tx, неймин методов идентичен
type DBClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Приватный тип ключа для контекста
type txKey struct{}

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() error {
	if db.Pool != nil {
		db.Pool.Close()
	}
	return nil
}

// WithinTransaction реализует интерфейс domain.Transactor
func (db *PostgresDB) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Если транзакция уже открыта выше по стеку (вложенный вызов) 
	// повторно открывать её не нужно — просто выполняем функцию в текущем контексте
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	// Открываем новую транзакцию в pgx
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось открыть транзакцию: %w", err)
	}

	// Гарантируем откат (Rollback) при возникновении паники или ошибки до вызова Commit
	// Если транзакция уже закрыта через Commit(), метод Rollback вернет pgx.ErrTxClosed, который мы игнорируем.
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			// Здесь можно вызвать логгер, например: logger.Error("rollback error", err)
		}
	}()

	// Обогащаем контекст, сохраняя в него открытую транзакцию
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Выполняем переданную бизнес-логику. Если внутри fn будет ошибка, возвращаем её.
	if err := fn(txCtx); err != nil {
		return err 
	}

	// Если вся цепочка операций выполнена без ошибок — фиксируем изменения в БД
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("не удалось закоммитить транзакцию: %w", err)
	}

	return nil
}

// GetDB проверяет контекст. Если в нем есть активная транзакция — возвращает её.
// Если транзакции нет — возвращает стандартный пул соединений.
func (db *PostgresDB) GetDB(ctx context.Context) DBClient {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return db.Pool
}