package logger

import (
	"context"
	"log/slog"
)

// Converts []slog.Attr to []any for variadic log.With
func attrsToAny(attrs []slog.Attr) []any {
	res := make([]any, len(attrs))
	for i, a := range attrs {
		res[i] = a
	}
	return res
}

type ctxKey struct{}

type Attributes map[string]any

func AddAttrs(ctx context.Context, attrs Attributes) context.Context {
	existing, _ := ctx.Value(ctxKey{}).(Attributes)
	if existing == nil {
		existing = Attributes{}
	}
	for k, v := range attrs {
		existing[k] = v
	}
	return context.WithValue(ctx, ctxKey{}, existing)
}

func extractAttrs(ctx context.Context) []slog.Attr {
	raw := ctx.Value(ctxKey{})
	if raw == nil {
		return nil
	}
	attrs := raw.(Attributes)

	slogAttrs := make([]slog.Attr, 0, len(attrs))
	for k, v := range attrs {
		slogAttrs = append(slogAttrs, slog.Any(k, v))
	}
	return slogAttrs
}

func Info(ctx context.Context, msg string, kv ...any) {
	kv = filterSensitive(kv)
	log.With(attrsToAny(extractAttrs(ctx))...).Info(msg, kv...)
}

func Debug(ctx context.Context, msg string, kv ...any) {
	kv = filterSensitive(kv)
	log.With(attrsToAny(extractAttrs(ctx))...).Debug(msg, kv...)
}

func Warn(ctx context.Context, msg string, kv ...any) {
	kv = filterSensitive(kv)
	log.With(attrsToAny(extractAttrs(ctx))...).Warn(msg, kv...)
}

func Error(ctx context.Context, msg string, kv ...any) {
	kv = filterSensitive(kv)
	log.With(attrsToAny(extractAttrs(ctx))...).Error(msg, kv...)
}
