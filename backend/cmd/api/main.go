package main

import (
	"log"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/repository"
	"github.com/Voltage11/tplatform/internal/service"
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
	defer dbPostgres.Close()
	appLogger.Info("Подключение к БД успешно")

	// Репозитории
	permissionRepo := repository.NewPermissionsRepository(dbPostgres.Pool)
	departmentRepo := repository.NewDepartmentRepository(dbPostgres.Pool)
	roleRepo := repository.NewRoleRepository(dbPostgres.Pool)
	userRepo := repository.NewUserRepository(dbPostgres.Pool)

	// Сервис
	permissionService := service.NewPermissionService(permissionRepo)
	defer permissionService.Shutdown()

	departmentService := service.NewDepartmentService(departmentRepo, permissionService)
	roleService := service.NewRoleService(roleRepo, permissionService)
	userService := service.NewUserService(userRepo, permissionService)

	_ = departmentService
	_ = roleService
	_ = userService

	// Установим уровень логирования из конфигурации после старта сервера
	appLogger.SetLevel(cfg.Logger.Level)
}
