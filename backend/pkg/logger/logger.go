package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	// Стандартные методы
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	// Методы с поддержкой контекста
	DebugCtx(ctx context.Context, msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	ErrorCtx(ctx context.Context, msg string, args ...any)

	SetLevel(levelStr string)
}

type appLogger struct {
	logger   *slog.Logger
	levelVar *slog.LevelVar
}

func New(initialLevel string) Logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(parseLevel(initialLevel))

	// Базовый JSON обработчик
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelVar,
	})
	

	return &appLogger{
		logger:   slog.New(handler),
		levelVar: levelVar,
	}
}

func (l *appLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}
func (l *appLogger) Debug(msg string, args ...any) { l.logger.Debug(msg, args...) }
func (l *appLogger) Warn(msg string, args ...any)  { l.logger.Warn(msg, args...) }
func (l *appLogger) Error(msg string, args ...any) { l.logger.Error(msg, args...) }

func (l *appLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}
func (l *appLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}
func (l *appLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}
func (l *appLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *appLogger) SetLevel(levelStr string) { 
	l.logger.Info("уровень логирования", "установлен на", levelStr)
	l.levelVar.Set(parseLevel(levelStr)) 	
}

func parseLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
