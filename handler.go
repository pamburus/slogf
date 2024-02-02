// Package slogf provides a bridge between the slog and logf packages.
package slogf

import (
	"context"
	"log/slog"
	"slices"

	"github.com/ssgreg/logf"
	"github.com/ssgreg/logf/logfc"
)

// NewHandler returns a new slog.Handler which uses logf.Logger to log records.
func NewHandler() *Handler {
	return &Handler{nil, nil, logfc.Get}
}

// ---

// Handler is a slog.Handler implementation which uses logf.Logger to log records.
type Handler struct {
	fields []logf.Field
	groups []group
	logger func(context.Context) *logf.Logger
}

// WithLogger returns a new Handler with the given logger.
func (h *Handler) WithLogger(logger *logf.Logger) *Handler {
	return h.WithLoggerFunc(func(context.Context) *logf.Logger {
		return logger
	})
}

// WithLoggerFunc returns a new Handler with the given logger provider function.
func (h *Handler) WithLoggerFunc(logger func(context.Context) *logf.Logger) *Handler {
	h = h.fork()
	h.logger = logger

	return h
}

// Enabled returns true if the given level is enabled.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	var enabled bool

	h.logger(ctx).AtLevel(LogfLevel(level), func(logf.LogFunc) {
		enabled = true
	})

	return enabled
}

// Handle logs the given record.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	collectAttrs := func(fields []logf.Field) []logf.Field {
		record.Attrs(func(attr slog.Attr) bool {
			if field, ok := logfField(attr); ok {
				fields = append(fields, field)
			}

			return true
		})

		return fields
	}

	var fields []logf.Field
	if len(h.fields)+record.NumAttrs() != 0 {
		if len(h.groups) == 0 {
			fields = make([]logf.Field, 0, record.NumAttrs()+len(h.fields))
			fields = append(fields, h.fields...)
			fields = collectAttrs(fields)
		} else {
			g := groupEncoder{h, 0, nil}
			fields = make([]logf.Field, 0, record.NumAttrs()+h.groups[0].i+1)
			fields = append(fields, h.fields[:h.groups[0].i]...)
			fields = append(fields, logf.Object(h.groups[0].name, &g))
			i := len(fields)
			fields = collectAttrs(fields)
			g.suffix = fields[i:]
			fields = fields[:i]
		}
	}

	logfLog(record.Level, h.logger(ctx), record.Message, fields...)

	return nil
}

// WithAttrs returns a new Handler with the given attributes.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	h = h.fork()
	h.fields = slices.Grow(h.fields, len(attrs))

	for _, attr := range attrs {
		if field, ok := logfField(attr); ok {
			h.fields = append(h.fields, field)
		}
	}

	return h
}

// WithGroup returns a new Handler with the given group.
func (h *Handler) WithGroup(key string) slog.Handler {
	if key == "" {
		return h
	}

	h = h.fork()
	h.groups = append(h.groups, group{len(h.fields), key})

	return h
}

func (h *Handler) fork() *Handler {
	h = &Handler{
		slices.Clip(h.fields),
		slices.Clip(h.groups),
		h.logger,
	}

	return h
}

func (h *Handler) groupAttrRange(i int) (int, int) {
	begin := h.groups[i].i
	end := len(h.fields)
	if i+1 < len(h.groups) {
		end = h.groups[i+1].i
	}

	return begin, end
}

// ---

type group struct {
	i    int
	name string
}

// ---

type groupEncoder struct {
	h      *Handler
	i      int
	suffix []logf.Field
}

func (g *groupEncoder) EncodeLogfObject(enc logf.FieldEncoder) error {
	begin, end := g.h.groupAttrRange(g.i)
	for i := begin; i != end; i++ {
		g.h.fields[i].Accept(enc)
	}

	if g.i+1 < len(g.h.groups) {
		g.i++
		if len(g.suffix) != 0 || g.h.groups[g.i].i < len(g.h.fields) {
			enc.EncodeFieldObject(g.h.groups[g.i].name, g)
		}
	} else {
		for i := range g.suffix {
			g.suffix[i].Accept(enc)
		}
	}

	return nil
}

// ---

var (
	_ slog.Handler = (*Handler)(nil)
)
