package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tplatform/internal/config"
	"tplatform/internal/handler/auth_handler"
	"tplatform/internal/middleware"
	"tplatform/internal/repository/auth_repository"
	"tplatform/internal/service/auth_service"
	"tplatform/pkg/database"
	"tplatform/pkg/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

const (
	MAX_CONNS = 25
	MIN_CONNS = 5
)

func main() {
	const op = "main"
	// Logger
	appLogger := logger.New(logger.AppLoggerLevelInfo)

	//Config
	cfg, err := config.New()
	if err != nil {
		appLogger.Error(err, op, "Ошибка загрузки конфигурации")
		return
	}

	//Database
	dbConfig := database.Config{
		DSN:           cfg.DSN(),
		MigrationPath: cfg.DB.MigrationPath,
		MaxConns:      MAX_CONNS,
		MinConns:      MIN_CONNS,
	}
	db, err := database.New(dbConfig, appLogger)
	if err != nil {
		appLogger.Error(err, op, "Ошибка подключения к БД")
		return
	}
	defer db.Close()

	// Обработаем падение приложения
	defer func() {
		if err := recover(); err != nil {
			appLogger.Error(fmt.Errorf("PANIC: %v", err), "main")
		}
	}()

	// Repositories
	authRepo := auth_repository.New(db, appLogger)

	// Services
	authService := auth_service.New(authRepo, appLogger, cfg.Secret.Jwt)

	// Server
	router := chi.NewRouter()

	// Middleware - ПРАВИЛЬНЫЙ ПОРЯДОК
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.SecurityHeaders)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Timeout(30 * time.Second))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(middleware.NewCORSMiddleware(cfg.Server.Cors).Handler)
	router.Use(middleware.RequireJSONContentType)
	router.Use(middleware.AuthMiddleware(authService, appLogger))

	// Handlers
	auth_handler.New(router, authService, appLogger)

	// Перед запуском приложеня устанавливаем уровень логирования
	appLogger.SetLevel(getLoggerLevelFromConfig(cfg.LogLevel))

	// Run
	runServer := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	appLogger.Warn(fmt.Sprintf("Запуск сервера на: %s", runServer), "main")
	log.Fatal(http.ListenAndServe(runServer, router))

}

// Полученное значение из .env преобразуем в уровень логирования
func getLoggerLevelFromConfig(level string) logger.AppLoggerLevel {
	switch level {
	case "debug":
		return logger.AppLoggerLevelDebug
	case "info":
		return logger.AppLoggerLevelInfo
	case "warn":
		return logger.AppLoggerLevelWarn
	case "error":
		return logger.AppLoggerLevelError
	default:
		return logger.AppLoggerLevelError
	}
}
