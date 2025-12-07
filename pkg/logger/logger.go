package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

type Config struct {
	Level       string // debug, info, warn, error
	Format      string // json, text
	Environment string // development, production
}

func Init(cfg Config) *slog.Logger {
	once.Do(func() {
		var level slog.Level
		switch cfg.Level {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}

		opts := &slog.HandlerOptions{
			Level:     level,
			AddSource: cfg.Environment == "development",
		}

		var handler slog.Handler
		if cfg.Format == "json" || cfg.Environment == "production" {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

		instance = slog.New(handler)
		slog.SetDefault(instance)
	})

	return instance
}

func Get() *slog.Logger {
	if instance == nil {
		return slog.Default()
	}
	return instance
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func With(args ...any) *slog.Logger {
	return Get().With(args...)
}

func WithContext(ctx context.Context) *slog.Logger {
	return Get()
}
