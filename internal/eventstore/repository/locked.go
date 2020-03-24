package repository

import (
	"context"

	"github.com/caos/eventstore-lib/pkg/models"
)

func (sql *SQL) LockAggregates(ctx context.Context, aggregates ...models.Aggregate) (err error) {
	return nil
}

func (sql *SQL) UnlockAggregates(ctx context.Context, aggregates ...models.Aggregate) (err error) {
	return nil
}
