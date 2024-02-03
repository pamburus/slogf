package slogf_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/pamburus/slogf"
	"github.com/ssgreg/logf"
)

func BenchmarkAtLevel(b *testing.B) {
	benchLogf(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
		b.ResetTimer()
		counter := 0
		for i := 0; i < b.N; i++ {
			logger.AtLevel(logf.LevelInfo, func(logf.LogFunc) {
				counter++
			})
		}
		b.StopTimer()
		logger.Info("test", logf.Int("counter", counter))
	})
}

func BenchmarkLogging(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.Info("test", logf.String("key", "value"))
				}
				b.StopTimer()
			})
		})
	})
	b.Run("Medium", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				logger = logger.With(
					slog.String("a", "a1"),
					slog.Int("b", 42),
					slog.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("c", "d"), slog.Int("e", 10), slog.String("f", "g"))
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				logger = logger.With(
					slog.String("a", "a1"),
					slog.Int("b", 42),
					slog.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("c", "d"), slog.Int("e", 10), slog.String("f", "g"))
				}
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
				logger = logger.With(
					logf.String("a", "a1"),
					logf.Int("b", 42),
					logf.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.Info("test", logf.String("c", "d"), logf.Int("e", 10), logf.String("f", "g"))
				}
				b.StopTimer()
			})
		})
	})
	b.Run("With", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlogLevel(slog.LevelDebug, b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogfLevel(logf.LevelDebug, false, b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogfLevel(logf.LevelDebug, b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.With(
						logf.String("a", "a1"),
						logf.Int("b", 42),
						logf.String("x", "x1"),
					)
				}
				b.StopTimer()
			})
		})
	})
	b.Run("With+Log", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
					logger.LogAttrs(ctx, slog.LevelInfo, "test")
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
					logger.LogAttrs(ctx, slog.LevelInfo, "test")
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						logf.String("a", "a1"),
						logf.Int("b", 42),
						logf.String("x", "x1"),
					)
					logger.Info("test")
				}
				b.StopTimer()
			})
		})
	})
}

func benchSlogf(b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	test := func(b *testing.B, withCaller bool, f func(context.Context, *testing.B, *slog.Logger)) {
		b.Run("Pass", func(b *testing.B) {
			benchSlogfLevel(logf.LevelDebug, withCaller, b, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchSlogfLevel(logf.LevelWarn, withCaller, b, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, true, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
			f(ctx, b, logger)
		})
	})
	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, false, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
			f(ctx, b, logger)
		})
	})
}

func benchSlogfLevel(level logf.Level, withCaller bool, b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	handler := slogf.NewHandler()
	appender := logf.NewWriteAppender(io.Discard, logf.NewJSONEncoder(logf.JSONEncoderConfig{
		EncodeDuration:     logf.NanoDurationEncoder,
		EncodeTime:         logf.RFC3339NanoTimeEncoder,
		DisableFieldCaller: true,
	}))

	logfLogger := logf.NewLogger(
		level,
		logf.NewUnbufferedEntryWriter(appender),
	)

	if withCaller {
		logfLogger = logfLogger.WithCaller()
	}

	ctx := logf.NewContext(context.Background(), logfLogger)
	logger := slog.New(handler.WithGroup("").WithAttrs(nil))
	f(ctx, b, logger)
	_ = appender.Flush()
}

func benchSlog(b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	b.Run("Pass", func(b *testing.B) {
		benchSlogLevel(slog.LevelDebug, b, f)
	})
	b.Run("Drop", func(b *testing.B) {
		benchSlogLevel(slog.LevelWarn, b, f)
	})
}

func benchSlogLevel(level slog.Level, b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	options := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(io.Discard, options)
	logger := slog.New(handler)
	f(context.Background(), b, logger)
}

func benchLogf(b *testing.B, f func(context.Context, *testing.B, *logf.Logger)) {
	test := func(b *testing.B, f func(context.Context, *testing.B, *logf.Logger)) {
		b.Run("Pass", func(b *testing.B) {
			benchLogfLevel(logf.LevelDebug, b, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchLogfLevel(logf.LevelWarn, b, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
			f(ctx, b, logger.WithCaller())
		})
	})

	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, f)
	})
}

func benchLogfLevel(level logf.Level, b *testing.B, f func(context.Context, *testing.B, *logf.Logger)) {
	appender := logf.NewWriteAppender(io.Discard, logf.NewJSONEncoder(logf.JSONEncoderConfig{
		EncodeDuration:     logf.NanoDurationEncoder,
		EncodeTime:         logf.RFC3339NanoTimeEncoder,
		DisableFieldCaller: true,
	}))

	logger := logf.NewLogger(
		level,
		logf.NewUnbufferedEntryWriter(appender),
	)

	f(context.Background(), b, logger)
	_ = appender.Flush()
}