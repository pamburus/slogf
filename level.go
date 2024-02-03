package slogf

import (
	"log/slog"

	"github.com/ssgreg/logf"
)

// LogfLevel converts slog.Level to logf.Level.
func LogfLevel(level slog.Level) logf.Level {
	switch {
	case level >= slog.LevelError:
		return logf.LevelError
	case level >= slog.LevelWarn:
		return logf.LevelWarn
	case level >= slog.LevelInfo:
		return logf.LevelInfo
	default:
		return logf.LevelDebug
	}
}

func logfLog(level slog.Level, logger *logf.Logger, text string, fields ...logf.Field) {
	switch {
	case level >= slog.LevelError:
		logger.Error(text, fields...)
	case level >= slog.LevelWarn:
		logger.Warn(text, fields...)
	case level >= slog.LevelInfo:
		logger.Info(text, fields...)
	default:
		logger.Debug(text, fields...)
	}
}
