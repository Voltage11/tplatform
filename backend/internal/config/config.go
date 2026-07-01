package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JWTConfig
	Admin    AdminConfig
}

type AdminConfig struct {
	Email    string `env:"ADMIN_EMAIL" env-required:"true"`
	Password string `env:"ADMIN_PASSWORD" env-required:"true"`
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"60s"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" env-separator:"," env-required:"true"` 
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
	MaxConns int32  `env:"DB_MAX_CONNS" env-default:"10"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// Validate проверяет корректность конфигурации БД
func (c *DatabaseConfig) Validate() error {
	if c.User == "" {
		return fmt.Errorf("DB_USER не может быть пустым")
	}
	if c.Password == "" {
		return fmt.Errorf("DB_PASSWORD не может быть пустым")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME не может быть пустым")
	}
	return nil
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type JWTConfig struct {
	Secret     string        `env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"720h"` // 30 дней
}

// Validate проверяет корректность JWT конфигурации
func (c *JWTConfig) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("JWT_SECRET не может быть пустым")
	}
	// if c.ExpireHours <= 0 {
	// 	return fmt.Errorf("JWT_EXPIRE_HOURS должен быть больше 0")
	// }
	return nil
}

func New() (*Config, error) {
	var cfg Config

	// Пробуем прочитать .env, если он есть
	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("ошибка чтения файла .env: %w", err)
		}
	}

	// Читаем переменные окружения (перезаписывают значения из .env)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}

	// Валидация
	if err := cfg.Database.Validate(); err != nil {
		return nil, fmt.Errorf("невалидная конфигурация БД: %w", err)
	}
	if err := cfg.JWT.Validate(); err != nil {
		return nil, fmt.Errorf("невалидная JWT конфигурация: %w", err)
	}

	return &cfg, nil
}