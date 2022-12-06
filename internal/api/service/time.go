package service

import (
	"context"
	"time"
)

type callTimeKey struct{}

var callKey *callTimeKey = (*callTimeKey)(nil)

func WithCallTime(parent context.Context, t time.Time) context.Context {
	if existing := parent.Value(callKey); existing != nil {
		return parent
	}
	return context.WithValue(parent, callKey, t)
}

func WithCallTimeNow(parent context.Context) context.Context {
	return WithCallTime(parent, time.Now())
}

func CallTimeFromContext(ctx context.Context) (callTime time.Time) {
	value := ctx.Value(callKey)
	callTime, _ = value.(time.Time)

	return callTime
}
