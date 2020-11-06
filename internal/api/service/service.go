package service

import "context"

type serviceKey struct{}

var key *serviceKey = (*serviceKey)(nil)

func WithService(partent context.Context, serviceName string) context.Context {
	return context.WithValue(partent, key, serviceName)
}

func FromContext(ctx context.Context) string {
	value := ctx.Value(key)
	if name, ok := value.(string); ok {
		return name
	}

	return ""
}
