package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"tplatform/pkg/logger"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Pool   *pgxpool.Pool
	logger logger.AppLogger
}

type Config struct {
	DSN           string
	MigrationPath string
	MaxConns      int32
	MinConns      int32
}

func New(cfg Config, appLogger logger.AppLogger) (*Database, error) {
	const op = "database.New"

	// Устанавливаем значения по умолчанию
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 25
	}
	if cfg.MinConns == 0 {
		cfg.MinConns = 5
	}

	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка парсинга DSN: %w", op, err)
	}

	// Настройки пула соединений
	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 30 * time.Second
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания пула: %w", op, err)
	}

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: ошибка ping к бд: %w", op, err)
	}

	database := &Database{
		Pool:   pool,
		logger: appLogger,
	}

	// Выполняем миграции если указан путь
	if cfg.MigrationPath != "" {
		appLogger.Info("Запуск миграций БД", op, "path", cfg.MigrationPath)
		if err := database.migrate(cfg.MigrationPath); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return database, nil
}

func (d *Database) migrate(migrationPath string) error {
	const op = "database.migrate"

	// Создаем соединение через stdlib для мигратора
	sqlDB, err := d.createSQLDB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("ошибка драйвера миграций: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	d.logger.Info("Миграции выполнены успешно (или изменений нет)", op)
	return nil
}

func (d *Database) createSQLDB() (*sql.DB, error) {
	dsn := d.Pool.Config().ConnConfig.ConnString()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sql.DB: %w", err)
	}

	return db, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}