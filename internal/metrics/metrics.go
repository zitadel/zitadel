package metrics

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

type Meter interface {
	RegisterCounter(name string) error
	AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error
}

type Config interface {
	NewMetrics() error
}

var M Meter

func RegisterCounter(name string) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.RegisterCounter(name)
}

func AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.AddCount(ctx, name, value, labels)
}
