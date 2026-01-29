package logger

import (
	"context"
	"log/slog"
	"os"
)

type AppLoggerLevel string

const (
	AppLoggerLevelDebug AppLoggerLevel = "debug"
	AppLoggerLevelInfo  AppLoggerLevel = "info"
	AppLoggerLevelWarn  AppLoggerLevel = "warn"
	AppLoggerLevelError AppLoggerLevel = "error"
)

type AppLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, op string, args ...any)
	Warn(msg string, op string, args ...any)
	Error(err error, op string, args ...any)
	InfoContext(ctx context.Context, msg string, op string, args ...any)
	SetLevel(level AppLoggerLevel) 
}

type appLogger struct {
	logger *slog.Logger
	level  *slog.LevelVar 
}

func getLoggerLevel(level AppLoggerLevel) slog.Level {
	switch level {
	case AppLoggerLevelDebug:
		return slog.LevelDebug
	case AppLoggerLevelInfo:
		return slog.LevelInfo
	case AppLoggerLevelWarn:
		return slog.LevelWarn
	case AppLoggerLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func New(levelCfg AppLoggerLevel) AppLogger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(getLoggerLevel(levelCfg))

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelVar,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "timestamp", Value: a.Value}
			}
			if a.Key == slog.LevelKey {
				return slog.Attr{Key: "level", Value: a.Value}
			}
			if a.Key == slog.MessageKey {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return a
		},
	})

	logger := slog.New(handler).With(slog.String("env", "tplatform"))
	
	// Устанавливаем как стандартный для системы
	slog.SetDefault(logger)

	return &appLogger{
		logger: logger,
		level:  levelVar,
	}
}

func (l *appLogger) SetLevel(level AppLoggerLevel) {
	// Просто меняем значение в LevelVar, логгер подхватит его автоматически
	l.level.Set(getLoggerLevel(level))
}

func (l *appLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *appLogger) Info(msg string, op string, args ...any) {
	l.logger.With(slog.String("op", op)).Info(msg, args...)
}

func (l *appLogger) Warn(msg string, op string, args ...any) {
	l.logger.With(slog.String("op", op)).Warn(msg, args...)
}

func (l *appLogger) Error(err error, op string, args ...any) {
	if err == nil {
		return
	}
	l.logger.With(
		slog.String("op", op),
		slog.String("error", err.Error()),
	).Error("operation failed", args...)
}

func (l *appLogger) InfoContext(ctx context.Context, msg string, op string, args ...any) {
	l.logger.With(slog.String("op", op)).InfoContext(ctx, msg, args...)
}