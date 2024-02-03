package slogf

import (
	"log/slog"

	"github.com/ssgreg/logf"
)

func logfField(attr slog.Attr) (logf.Field, bool) {
	if attr.Equal(slog.Attr{}) {
		return logf.Field{}, false
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return logf.Bool(attr.Key, attr.Value.Bool()), true
	case slog.KindDuration:
		return logf.Duration(attr.Key, attr.Value.Duration()), true
	case slog.KindFloat64:
		return logf.Float64(attr.Key, attr.Value.Float64()), true
	case slog.KindInt64:
		return logf.Int64(attr.Key, attr.Value.Int64()), true
	case slog.KindString:
		return logf.String(attr.Key, attr.Value.String()), true
	case slog.KindTime:
		return logf.Time(attr.Key, attr.Value.Time()), true
	case slog.KindUint64:
		return logf.Uint64(attr.Key, attr.Value.Uint64()), true
	case slog.KindGroup:
		return logf.Object(attr.Key, &object{logfFields(attr.Value.Group()...)}), true
	case slog.KindLogValuer:
		return logfField(slog.Attr{Key: attr.Key, Value: attr.Value.Resolve()})
	case slog.KindAny:
		fallthrough
	default:
		return logf.Any(attr.Key, attr.Value.Any()), true
	}
}

func logfFields(attrs ...slog.Attr) []logf.Field {
	fields := make([]logf.Field, 0, len(attrs))

	for _, attr := range attrs {
		if field, ok := logfField(attr); ok {
			fields = append(fields, field)
		}
	}

	return fields
}

// ---

type object struct {
	fields []logf.Field
}

func (o *object) EncodeLogfObject(enc logf.FieldEncoder) error {
	for _, field := range o.fields {
		field.Accept(enc)
	}

	return nil
}

// ---

var _ logf.ObjectEncoder = (*object)(nil)
