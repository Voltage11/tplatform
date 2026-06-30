// Применение миграций
// go run cmd/migrate/main.go -cmd up

// Откат всех миграций
// go run cmd/migrate/main.go -cmd down

// Откат на 2 шага
// go run cmd/migrate/main.go -cmd down -step 2

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Инициализируем флаги напрямую
	command := flag.String("cmd", "up", "Команда миграции: up (накатить) или down (откатить)")
	step := flag.Int("step", 0, "Количество шагов для отката (используется с cmd=down)")
	flag.Parse()

	// Валидация входных данных на старте
	if *command != "up" && *command != "down" {
		log.Fatalf("Ошибка: недопустимая команда %q. Используйте 'up' или 'down'", *command)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Передаем значения (разъименовываем указатели через *)
	if err := runMigrations(cfg.Database.DSN(), *command, *step); err != nil {
		log.Fatalf("Миграция завершилась с ошибкой: %v", err)
	}
}

func runMigrations(databaseURL string, command string, step int) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать мигратор: %w", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("База данных в актуальном состоянии (нет новых миграций)")
				return nil
			}
			return fmt.Errorf("ошибка при выполнении UP: %w", err)
		}
		log.Println("Все новые миграции успешно применены!")

	case "down":
		if step > 0 {
			if err := m.Steps(-step); err != nil {
				return fmt.Errorf("ошибка при откате на %d шагов: %w", step, err)
			}
			log.Printf("Успешно откатано миграций: %d шагов\n", step)
		} else {
			if err := m.Down(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					log.Println("Нечего откатывать, база пуста")
					return nil
				}
				return fmt.Errorf("ошибка при полном откате DOWN: %w", err)
			}
			log.Println("Все миграции успешно откатаны (база очищена)!")
		}
	}

	return nil
}
