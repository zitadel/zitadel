package command

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// Retry the whole operation: replaying previously computed events would retain
// stale membership differences and reservation releases.
func (c *Commands) retryGroupChange(ctx context.Context, change func(context.Context) (*domain.ObjectDetails, error)) (*domain.ObjectDetails, error) {
	for range 3 {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var details *domain.ObjectDetails
		err := c.eventstore.WithInstanceLock(ctx, func(ctx context.Context) error {
			var err error
			details, err = change(ctx)
			return err
		})
		if !errors.Is(err, eventstore.ErrConcurrentModification) {
			return details, err
		}
	}
	return nil, eventstore.ErrConcurrentModification
}
