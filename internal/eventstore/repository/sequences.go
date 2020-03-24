package repository

import (
	"context"

	"github.com/caos/eventstore-lib/pkg/models"
)

type sequence struct {
	Sequence uint64 `gorm:"column:event_sequence"`
}

func (sql *SQL) ValidateLatestSequence(ctx context.Context, aggregates ...models.Aggregate) (err error) {
	return nil
}
