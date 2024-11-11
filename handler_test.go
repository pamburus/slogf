package slogf_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/ssgreg/logf"

	. "github.com/pamburus/go-tst/tst"
	"github.com/pamburus/slogf"
)

func TestHandler(tt *testing.T) {
	t := New(tt)

	type test struct {
		lineTag  LineTag
		name     string
		log      func(context.Context, *slog.Logger)
		expected []string
	}

	tests := []test{
		{
			lineTag: ThisLine(),
			name:    "Simple",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","key":"value"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "With",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger = logger.With(
					slog.String("a", "a1"),
					slog.Int("b", 42),
				)

				logger.LogAttrs(ctx, slog.LevelInfo, "test 1", slog.String("key 1", "value 1"))
				logger.LogAttrs(ctx, slog.LevelInfo, "test 2", slog.String("key 2", "value 2"))
			},
			expected: []string{
				`{"level":"info","msg":"test 1","a":"a1","b":42,"key 1":"value 1"}`,
				`{"level":"info","msg":"test 2","a":"a1","b":42,"key 2":"value 2"}`,
			},
		},
		{
			lineTag: ThisLine(),
			name:    "WithGroup",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"key":"value"}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "WithGroup:x2",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").WithGroup("g2").LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"g2":{"key":"value"}}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "WithGroup:x3:0:1",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").WithGroup("g2").WithGroup("g3").With(slog.Int("x", 42)).LogAttrs(ctx, slog.LevelInfo, "test")
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"g2":{"g3":{"x":42}}}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "WithGroup:x2:0",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").WithGroup("g2").LogAttrs(ctx, slog.LevelInfo, "test")
			},
			expected: []string{`{"level":"info","msg":"test"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "Mix1",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").With(slog.Bool("b", true)).WithGroup("g2").With(slog.Int("c", 43)).LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"b":true,"g2":{"c":43,"key":"value"}}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "Mix2",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.With(slog.Float64("f", 30)).WithGroup("g1").With(slog.Bool("b", true)).WithGroup("g2").With(slog.Int("c", 43)).LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","f":30,"g1":{"b":true,"g2":{"c":43,"key":"value"}}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "Mix3",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("g1").With(slog.Bool("b", true)).WithGroup("g2").LogAttrs(ctx, slog.LevelInfo, "test")
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"b":true}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "LevelDebug",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelDebug, "test")
			},
			expected: []string{`{"level":"debug","msg":"test"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "LevelWarning",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelWarn, "test")
			},
			expected: []string{`{"level":"warn","msg":"test"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "LevelError",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelError, "test")
			},
			expected: []string{`{"level":"error","msg":"test"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "EmptyAttrs",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Attr{})
			},
			expected: []string{`{"level":"info","msg":"test"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "EmptyGroup",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.WithGroup("").LogAttrs(ctx, slog.LevelInfo, "test", slog.String("key", "value"))
			},
			expected: []string{`{"level":"info","msg":"test","key":"value"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueDuration",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Duration("key", time.Second))
			},
			expected: []string{`{"level":"info","msg":"test","key":1000000000}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueTime",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Time("key", time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)))
			},
			expected: []string{`{"level":"info","msg":"test","key":"2020-01-02T03:04:05.000000006Z"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueUint",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Uint64("key", 42))
			},
			expected: []string{`{"level":"info","msg":"test","key":42}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueGroup",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Group("g1", slog.Uint64("key", 42)))
			},
			expected: []string{`{"level":"info","msg":"test","g1":{"key":42}}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueAny",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Any("e", errors.New("err")))
			},
			expected: []string{`{"level":"info","msg":"test","e":"err"}`},
		},
		{
			lineTag: ThisLine(),
			name:    "ValueValuer",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, "test", slog.Any("v", testValuer{42}))
			},
			expected: []string{`{"level":"info","msg":"test","v":42}`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t Test) {
			t.Run("slog", func(t Test) {
				t.AddLineTags(test.lineTag)
				t.Expect(testLog(testSlog(test.log))).To(Equal(test.expected))
			})
			t.Run("slogf", func(t Test) {
				t.AddLineTags(test.lineTag)
				t.Expect(testLog(testSlogf(test.log))).To(Equal(test.expected))
			})
		})
	}

	t.Run("WithLogger", func(t Test) {
		t.Expect(
			testLog(testLogf(func(logfLogger *logf.Logger) {
				logger := slog.New(slogf.NewHandler().WithLogger(logfLogger))
				logger.Info("test", slog.String("key", "value"))
			})),
		).To(Equal([]string{`{"level":"info","msg":"test","key":"value"}`}))
	})
}

func testLog(f func(io.Writer)) []string {
	buffer := bytes.NewBuffer(nil)
	f(buffer)

	return strings.Split(strings.TrimSpace(buffer.String()), "\n")
}

func testSlogf(f func(context.Context, *slog.Logger)) func(io.Writer) {
	return func(writer io.Writer) {
		handler := slogf.NewHandler()
		appender := logf.NewWriteAppender(writer, logf.NewJSONEncoder(logf.JSONEncoderConfig{
			DisableFieldTime: true,
			EncodeDuration:   logf.NanoDurationEncoder,
			EncodeTime:       logf.RFC3339NanoTimeEncoder,
		}))

		ctx := logf.NewContext(
			context.Background(),
			logf.NewLogger(
				logf.LevelDebug,
				logf.NewUnbufferedEntryWriter(appender),
			),
		)

		logger := slog.New(handler.WithGroup("").WithAttrs(nil))
		f(ctx, logger)
		err := appender.Flush()
		if err != nil {
			panic(err)
		}
	}
}

func testSlog(f func(context.Context, *slog.Logger)) func(io.Writer) {
	options := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			if len(groups) == 0 && a.Key == slog.LevelKey {
				return slog.Attr{Key: "level", Value: slog.StringValue(strings.ToLower(a.Value.String()))}
			}

			return a
		},
	}

	return func(writer io.Writer) {
		handler := slog.NewJSONHandler(writer, options)
		logger := slog.New(handler)
		f(context.Background(), logger)
	}
}

func testLogf(f func(*logf.Logger)) func(io.Writer) {
	return func(writer io.Writer) {
		appender := logf.NewWriteAppender(writer, logf.NewJSONEncoder(logf.JSONEncoderConfig{
			DisableFieldTime: true,
			EncodeDuration:   logf.NanoDurationEncoder,
			EncodeTime:       logf.RFC3339NanoTimeEncoder,
		}))

		logger := logf.NewLogger(
			logf.LevelDebug,
			logf.NewUnbufferedEntryWriter(appender),
		)

		f(logger)
		err := appender.Flush()
		if err != nil {
			panic(err)
		}
	}
}

// ---

type testValuer struct {
	value int
}

func (v testValuer) LogValue() slog.Value {
	return slog.IntValue(v.value)
}

// ---

var (
	_ slog.LogValuer = testValuer{}
)
