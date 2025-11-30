package logger

import "log/slog"

type sensitiveValue struct {
	key string
	val any
}

func Sensitive(key string, val any) sensitiveValue {
	return sensitiveValue{key: key, val: val}
}

func maskValue(v any) any {
	return "****"
}

func filterSensitive(kvs []any) []any {
	result := make([]any, 0, len(kvs))
	for _, kv := range kvs {
		switch v := kv.(type) {
		case sensitiveValue:
			if IsProd() {
				result = append(result, slog.Any(v.key, maskValue(v.val)))
			} else {
				result = append(result, slog.Any(v.key, v.val))
			}
		default:
			result = append(result, kv)
		}
	}
	return result
}
