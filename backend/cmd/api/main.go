package main

import (
	"log"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/config/db"
	"github.com/Voltage11/tplatform/pkg/logger"
)

const (
	LOGGER_LEVEL = "info" // Стартовый уровень логирования до запуска сервера
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ошибка получения конфигурации: %v", err)
	}

	// Инициализация логгера
	appLogger := logger.New(LOGGER_LEVEL)

	// Подключение к Postgres
	dbPostgres, err := db.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("ошибка подключения к postgres: %v", err)
	}
	appLogger.Info("Подключение к БД успешно")

	// Установим уровень логирования из конфигурации после старта сервера
	appLogger.SetLevel(cfg.Logger.Level)
}
