package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/handler"
	authmw "github.com/Voltage11/tplatform/internal/middleware"
	"github.com/Voltage11/tplatform/internal/repository"
	"github.com/Voltage11/tplatform/internal/service"
	"github.com/Voltage11/tplatform/pkg/jwt"
	"github.com/Voltage11/tplatform/pkg/logger"
)

const (
	loggerLevel     = "info"
	shutdownTimeout = 30 * time.Second // максимальное время на завершение
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ошибка получения конфигурации: %v", err)
	}

	// Инициализация логгера
	appLogger := logger.New(loggerLevel)

	// Подключение к Postgres
	dbPostgres, err := db.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("ошибка подключения к postgres: %v", err)
	}
	defer dbPostgres.Close()
	appLogger.Info("Подключение к БД успешно")

	// Репозитории
	sessionRepo := repository.NewSessionRepository(dbPostgres)
	permissionRepo := repository.NewPermissionsRepository(dbPostgres)
	departmentRepo := repository.NewDepartmentRepository(dbPostgres)
	roleRepo := repository.NewRoleRepository(dbPostgres)
	userRepo := repository.NewUserRepository(dbPostgres)

	// Сервисы
	jwtCfg := jwt.Config{
		SecretKey:  cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	}

	permissionService := service.NewPermissionService(permissionRepo)
	defer permissionService.Shutdown()

	authService := service.NewAuthService(userRepo, sessionRepo, jwtCfg)
	departmentService := service.NewDepartmentService(departmentRepo, permissionService)
	roleService := service.NewRoleService(roleRepo, permissionService)
	userService := service.NewUserService(userRepo, permissionService, dbPostgres)
	defer userService.Shutdown() // остановка кэша пользователей

	// Создание / проверка администратора
	if err := userService.CheckOrCreateAdmin(context.Background(), cfg.Admin); err != nil {
		log.Fatalf("ошибка создания админа: %v", err)
	}

	// Роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Server.ReadTimeout))

	authMiddleware := authmw.NewAuthMiddleware(authService, userService, appLogger)
	r.Use(authMiddleware.ExtractUser)

	handler.NewAuthHandlers(r, authMiddleware, authService, userService, appLogger)
	handler.NewDepartmentHandler(r, authMiddleware, departmentService)
	handler.NewRoleHandler(r, authMiddleware, roleService)
	handler.NewUserHandler(r, authMiddleware, userService)
	handler.NewPermissionHandler(r, authMiddleware, permissionService)

	// HTTP-сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Установим уровень логирования из конфигурации
	appLogger.SetLevel(cfg.Logger.Level)

	// Запуск сервера в фоне и ожидание сигнала для остановки
	// Канал для получения ошибок от сервера
	serverErr := make(chan error, 1)

	go func() {
		appLogger.Info("Сервер запущен", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("ошибка сервера: %w", err)
		}
	}()

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем либо ошибку от сервера, либо сигнал на завершение
	select {
	case err := <-serverErr:
		appLogger.Error("Критическая ошибка сервера", "error", err)
		// не вызываем os.Exit, дадим defer отработать
		// но завершим программу с ненулевым кодом
		log.Fatalf("Сервер упал: %v", err)

	case sig := <-quit:
		appLogger.Info("Получен сигнал на остановку", "signal", sig.String())

		// Создаём контекст с таймаутом для плавного завершения
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// 1. Останавливаем HTTP-сервер (перестаём принимать новые запросы)
		if err := srv.Shutdown(ctx); err != nil {
			appLogger.Error("Ошибка при остановке HTTP-сервера", "error", err)
		} else {
			appLogger.Info("HTTP-сервер остановлен")
		}

		// 2. Останавливаем фоновые сервисы (если требуется)
		// permissionService.Shutdown() вызовется через defer, но можно и явно здесь
		// userService.Shutdown() тоже через defer
		// Дополнительно можно добавить явные вызовы с таймаутом, если нужно гарантировать порядок:
		// permissionService.Shutdown()
		// userService.Shutdown()

		// 3. Закрытие БД произойдёт через defer dbPostgres.Close()

		appLogger.Info("Приложение завершено корректно")
		// Выход с кодом 0 произойдёт после возврата из main
	}
}
