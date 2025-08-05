package authz

import (
	"context"
	"time"
)

func Detach(ctx context.Context) context.Context { return detachedContext{ctx} }

type detachedContext struct {
	parent context.Context
}

func (v detachedContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (v detachedContext) Done() <-chan struct{}             { return nil }
func (v detachedContext) Err() error                        { return nil }
func (v detachedContext) Value(key interface{}) interface{} { return v.parent.Value(key) }
