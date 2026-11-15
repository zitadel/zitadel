package command

import (
	"context"
	"errors"
	"sync"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Retry the whole operation: replaying previously computed events would retain
// stale membership differences and reservation releases.
func (c *Commands) retryGroupChange(ctx context.Context, change func(context.Context) (*domain.ObjectDetails, error)) (*domain.ObjectDetails, error) {
	refreshed := false
	for range 3 {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		details, err := c.groupChangeOnce(ctx, change)
		// Projections are not triggered inside the lock transaction, so a membership
		// granted moments ago is not visible to the permission check there.
		if !refreshed && c.refreshMemberProjections != nil && isMembershipNotFound(err) {
			refreshed = true
			c.refreshMemberProjections(ctx)
			details, err = c.groupChangeOnce(ctx, change)
		}
		if !errors.Is(err, eventstore.ErrConcurrentModification) {
			return details, err
		}
	}
	return nil, eventstore.ErrConcurrentModification
}

func (c *Commands) groupChangeOnce(ctx context.Context, change func(context.Context) (*domain.ObjectDetails, error)) (*domain.ObjectDetails, error) {
	var details *domain.ObjectDetails
	err := c.eventstore.WithInstanceLock(ctx, func(ctx context.Context) error {
		var err error
		details, err = change(ctx)
		return err
	})
	return details, err
}

// isMembershipNotFound matches the error of [authz] when no membership of the caller is projected yet.
func isMembershipNotFound(err error) bool {
	zitadelErr, ok := zerrors.AsZitadelError(err)
	return ok && zitadelErr.ID == "AUTHZ-cdgFk"
}

func triggerMemberProjections(ctx context.Context) {
	var wg sync.WaitGroup
	for _, h := range []*handler.Handler{
		projection.OrgMemberProjection,
		projection.InstanceMemberProjection,
		projection.ProjectMemberProjection,
		projection.ProjectGrantMemberProjection,
	} {
		wg.Go(func() { _, _ = h.Trigger(ctx, handler.WithAwaitRunning()) })
	}
	wg.Wait()
}
