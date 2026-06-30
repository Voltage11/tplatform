package main

import (
	"log"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/pkg/logger"
)

func main() {
	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ошибка получения конфигурации: %v", err)
	}

	// Logger
	appLogger := logger.New("info")

	// Установим уровень логирования из конфигурации после старта сервера
	appLogger.SetLevel(cfg.Logger.Level)
}
