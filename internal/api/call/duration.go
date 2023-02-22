package call

import (
	"context"
	"time"
)

type durationKey struct{}

var key *durationKey = (*durationKey)(nil)

// WithTimestamp sets [time.Now()] adds the call field to the context
// if it's not already set
func WithTimestamp(parent context.Context) context.Context {
	if parent.Value(key) != nil {
		return parent
	}
	return context.WithValue(parent, key, time.Now())
}

// FromContext returns the [time.Time] the call hit the api
func FromContext(ctx context.Context) (t time.Time) {
	value := ctx.Value(key)
	if t, ok := value.(time.Time); ok {
		return t
	}

	return t
}

// Took returns the time the call took so far
func Took(ctx context.Context) time.Duration {
	start := FromContext(ctx)
	if start.IsZero() {
		return 0
	}
	return time.Since(start)
}
