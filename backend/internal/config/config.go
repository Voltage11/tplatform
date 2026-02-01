package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB struct {
		Host          string `env:"DB_HOST" env-required:"true"`
		Port          int    `env:"DB_PORT" env-required:"true"`
		User          string `env:"DB_USER" env-required:"true"`
		Password      string `env:"DB_PASSWORD" env-required:"true"`
		Database      string `env:"DB_DATABASE" env-required:"true"`
		MigrationPath string `env:"DB_MIGRATION_PATH" env-required:"true"`
	}
	Secret struct {
		Jwt  string `env:"SECRET_JWT" env-required:"true"`
		Hash string `env:"SECRET_HASH" env-required:"true"`
	}
	Server struct {
		Port int      `env:"SERVER_PORT" env-required:"true"`
		Host string   `env:"SERVER_HOST" env-required:"true"`
		Cors []string `env:"SERVER_CORS" env-required:"true" env-separator:","`
	}
	Email struct {
		Smtp struct {
			Host     string `env:"EMAIL_SMTP_HOST" env-required:"true"`
			Port     int    `env:"EMAIL_SMTP_PORT" env-required:"true"`
			Username string `env:"EMAIL_SMTP_USERNAME" env-required:"true"`
			Password string `env:"EMAIL_SMTP_PASSWORD" env-required:"true"`
		}
	}
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	// Media    struct {
	// 	FavIconsPath string `env:"ICONS_DIR" env-default:"./media/favicons"`
	// }
}

func New() (*Config, error) {

	instance := &Config{}

	envPath := ".env"

	// Проверяем существование файла .env
	if _, err := os.Stat(envPath); err == nil {
		err = cleanenv.ReadConfig(envPath, instance)
		if err != nil {
			return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
		}
	} else {
		// Загружаем только из переменных окружения
		err = cleanenv.ReadEnv(instance)
		if err != nil {
			return nil, fmt.Errorf("ошика чтения переменных окружения: %w", err)
		}
	}

	// Валидация дополнительных правил
	if err := instance.validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации значений переменных окружения: %w", err)
	}

	// Проверка существования директории с иконками, если отсутствует - создаем
	// if _, err := os.Stat(instance.Media.FavIconsPath); os.IsNotExist(err) {
	// 	err = os.MkdirAll(instance.Media.FavIconsPath, os.ModePerm)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("ошибка создания директории с иконками: %w", err)
	// 	}
	// }

	return instance, nil
}

func (c *Config) validate() error {
	// Валидация портов
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("invalid DB port: %d", c.DB.Port)
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Email.Smtp.Port <= 0 || c.Email.Smtp.Port > 65535 {
		return fmt.Errorf("invalid SMTP port: %d", c.Email.Smtp.Port)
	}

	// Валидация JWT секрета (минимум 32 символа)
	if len(c.Secret.Jwt) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// Валидация CORS
	if len(c.Server.Cors) == 0 {
		return fmt.Errorf("at least one CORS origin must be specified")
	}

	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Database,
	)
}
