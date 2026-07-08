package domain

import "context"

// Transactor управляет границами бизнес-транзакций
// Сервисы не будут зависеть от драйвера субд
type Transactor interface {
	// WithinTransaction оборачивает выполнение функции fn в транзакцию.
	// Если fn возвращает ошибку, транзакция откатывается. Если nil — коммитится
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
