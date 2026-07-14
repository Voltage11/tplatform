package repository

// import (
//     "context"

//     "github.com/jackc/pgx/v5"
//     "github.com/jackc/pgx/v5/pgconn"
//     "github.com/jackc/pgx/v5/pgxpool"
// )

// // Executor — интерфейс для выполнения запросов
// // Ему удовлетворяют *pgxpool.Pool и pgx.Tx.
// type Executor interface {
//     Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
//     Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
//     QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
// }

// // GetExecutor возвращает переданную транзакцию, если она не nil, иначе возвращает пул соединений.
// func GetExecutor(pool *pgxpool.Pool, tx pgx.Tx) Executor {
//     if tx != nil {
//         return tx
//     }
//     return pool
// }
