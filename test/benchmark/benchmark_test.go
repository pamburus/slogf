package benchmark_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/ssgreg/logf"

	"github.com/pamburus/slogf"
	"github.com/pamburus/slogx"
)

func BenchmarkAtLevel(b *testing.B) {
	benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
		b.Helper()
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
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
				}
				b.StopTimer()
			})
		})
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfX(b, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.Info("test", logf.String("key", "value"))
				}
				b.StopTimer()
			})
		})
	})
	b.Run("With3x", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlogLevel(b, slog.LevelDebug, false, func(_ context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
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
			benchSlogfLevel(b, logf.LevelDebug, false, func(_ context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
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
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfXLevel(b, logf.LevelDebug, false, func(_ context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
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
			benchLogfLevel(b, logf.LevelDebug, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
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
	b.Run("WithAndLog1x", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
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
				b.Helper()
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
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfX(b, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
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
			benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
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
	b.Run("WithAndLog3x", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
					logger.LogAttrs(ctx, slog.LevelInfo, "m1")
					logger.LogAttrs(ctx, slog.LevelInfo, "m2")
					logger.LogAttrs(ctx, slog.LevelInfo, "m3")
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
					logger.LogAttrs(ctx, slog.LevelInfo, "m1")
					logger.LogAttrs(ctx, slog.LevelInfo, "m2")
					logger.LogAttrs(ctx, slog.LevelInfo, "m3")
				}
				b.StopTimer()
			})
		})
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfX(b, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						slog.String("a", "a1"),
						slog.Int("b", 42),
						slog.String("x", "x1"),
					)
					logger.LogAttrs(ctx, slog.LevelInfo, "m1")
					logger.LogAttrs(ctx, slog.LevelInfo, "m2")
					logger.LogAttrs(ctx, slog.LevelInfo, "m3")
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger := logger.With(
						logf.String("a", "a1"),
						logf.Int("b", 42),
						logf.String("x", "x1"),
					)
					logger.Info("m1")
					logger.Info("m2")
					logger.Info("m3")
				}
				b.StopTimer()
			})
		})
	})
	b.Run("LogAfterWith3x0", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
				logger = logger.With(
					slog.String("a", "a1"),
					slog.Int("b", 42),
					slog.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test")
				}
				b.StopTimer()
			})
		})
		b.Run("slogf", func(b *testing.B) {
			benchSlogf(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
				logger = logger.With(
					slog.String("a", "a1"),
					slog.Int("b", 42),
					slog.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test")
				}
				b.StopTimer()
			})
		})
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfX(b, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
				logger = logger.WithLongTerm(
					slog.String("a", "a1"),
					slog.Int("b", 42),
					slog.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.LogAttrs(ctx, slog.LevelInfo, "test")
				}
				b.StopTimer()
			})
		})
		b.Run("logf", func(b *testing.B) {
			benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
				logger = logger.With(
					logf.String("a", "a1"),
					logf.Int("b", 42),
					logf.String("x", "x1"),
				)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					logger.Info("test")
				}
				b.StopTimer()
			})
		})
	})
	b.Run("LogAfterWith3x3", func(b *testing.B) {
		b.Run("slog", func(b *testing.B) {
			benchSlog(b, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
				b.Helper()
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
				b.Helper()
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
		b.Run("slogf+x", func(b *testing.B) {
			benchSlogfX(b, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
				b.Helper()
				logger = logger.WithLongTerm(
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
			benchLogf(b, func(_ context.Context, b *testing.B, logger *logf.Logger) {
				b.Helper()
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
}

func benchSlogf(b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	b.Helper()

	test := func(b *testing.B, withCaller bool, f func(context.Context, *testing.B, *slog.Logger)) {
		b.Helper()

		b.Run("Pass", func(b *testing.B) {
			benchSlogfLevel(b, logf.LevelDebug, withCaller, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchSlogfLevel(b, logf.LevelWarn, withCaller, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, true, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
			b.Helper()
			f(ctx, b, logger)
		})
	})
	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, false, func(ctx context.Context, b *testing.B, logger *slog.Logger) {
			b.Helper()
			f(ctx, b, logger)
		})
	})
}

func benchSlogfX(b *testing.B, f func(context.Context, *testing.B, *slogx.Logger)) {
	b.Helper()

	test := func(b *testing.B, withCaller bool, f func(context.Context, *testing.B, *slogx.Logger)) {
		b.Helper()

		b.Run("Pass", func(b *testing.B) {
			benchSlogfXLevel(b, logf.LevelDebug, withCaller, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchSlogfXLevel(b, logf.LevelWarn, withCaller, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, true, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
			b.Helper()
			f(ctx, b, logger.WithSource(true))
		})
	})
	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, false, func(ctx context.Context, b *testing.B, logger *slogx.Logger) {
			b.Helper()
			f(ctx, b, logger.WithSource(false))
		})
	})
}

func benchSlogfLevel(b *testing.B, level logf.Level, withCaller bool, f func(context.Context, *testing.B, *slog.Logger)) {
	b.Helper()
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
	logger := slog.New(handler)
	f(ctx, b, logger)
	_ = appender.Flush()
}

func benchSlogfXLevel(b *testing.B, level logf.Level, withCaller bool, f func(context.Context, *testing.B, *slogx.Logger)) {
	b.Helper()

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
	logger := slogx.New(handler)
	f(ctx, b, logger)
	_ = appender.Flush()
}

func benchSlog(b *testing.B, f func(context.Context, *testing.B, *slog.Logger)) {
	b.Helper()

	test := func(b *testing.B, addSource bool, f func(context.Context, *testing.B, *slog.Logger)) {
		b.Helper()

		b.Run("Pass", func(b *testing.B) {
			benchSlogLevel(b, slog.LevelDebug, addSource, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchSlogLevel(b, slog.LevelWarn, addSource, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, true, f)
	})

	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, false, f)
	})
}

func benchSlogLevel(b *testing.B, level slog.Level, addSource bool, f func(context.Context, *testing.B, *slog.Logger)) {
	b.Helper()

	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	handler := slog.NewJSONHandler(io.Discard, options)
	logger := slog.New(handler)
	f(context.Background(), b, logger)
}

func benchLogf(b *testing.B, f func(context.Context, *testing.B, *logf.Logger)) {
	b.Helper()

	test := func(b *testing.B, f func(context.Context, *testing.B, *logf.Logger)) {
		b.Helper()

		b.Run("Pass", func(b *testing.B) {
			benchLogfLevel(b, logf.LevelDebug, f)
		})
		b.Run("Drop", func(b *testing.B) {
			benchLogfLevel(b, logf.LevelWarn, f)
		})
	}

	b.Run("WithCaller", func(b *testing.B) {
		test(b, func(ctx context.Context, b *testing.B, logger *logf.Logger) {
			b.Helper()
			f(ctx, b, logger.WithCaller())
		})
	})

	b.Run("WithoutCaller", func(b *testing.B) {
		test(b, f)
	})
}

func benchLogfLevel(b *testing.B, level logf.Level, f func(context.Context, *testing.B, *logf.Logger)) {
	b.Helper()

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
