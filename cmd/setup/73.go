package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

var (
	//go:embed 73.sql
	fixUserGrantRoles string
)

// FixUserGrantRoles repairs grant-based user grants that kept roles which
// should have been cascade-removed (GHSA-v859-c572-qh5p). Only grant-based
// user grants (ug.grant_id set) are in scope: the bug can only manifest via
// ChangeProjectGrant's multi-role cascade to those grants (see 73.sql for
// details), so direct user grants are intentionally left untouched even if
// their roles also happen to mismatch, since that mismatch would be unrelated
// pre-existing data drift, not this CVE. The corruption lives in the
// eventstore event payloads, so we fix the source of truth by pushing a
// corrective user.grant.cascade.changed event with the reconciled role set,
// then let the projection catch up.
type FixUserGrantRoles struct {
	eventstore *eventstore.Eventstore
}

func (mig *FixUserGrantRoles) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceAddedEventType).
			Builder().ExcludeAggregateIDs().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}

	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "SETUP"})
	for i, instanceID := range instances {
		ctx = authz.WithInstanceID(ctx, instanceID)
		logging.Info(ctx, "fix user grant roles", "instance", instanceID, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances)))

		// The finder reads projection tables, which may not yet be populated during
		// setup. Bring the sources up to date for this instance before querying.
		for _, source := range []*handler.Handler{
			projection.ProjectGrantProjection,
			projection.UserGrantProjection,
		} {
			if _, err = source.Trigger(ctx); err != nil {
				return err
			}
		}

		fixedCount, err := mig.fixInstanceUserGrants(ctx, instanceID)
		if err != nil {
			return err
		}
		if fixedCount == 0 {
			continue
		}
		logging.Info(ctx, "fixed user grant roles", "instance", instanceID, "migration", mig.String(), "grantsFixed", fixedCount)
		// Reflect the corrective events we just pushed in the read model.
		_, err = projection.UserGrantProjection.Trigger(ctx)
		logging.OnError(ctx, err).Debug("failed triggering user grant projection to update roles")
	}
	return nil
}

func (mig *FixUserGrantRoles) fixInstanceUserGrants(ctx context.Context, instanceID string) (fixedCount int, err error) {
	var correctedGrants []eventstore.Command

	tx, err := mig.eventstore.Client().BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	rows, err := tx.QueryContext(ctx, fixUserGrantRoles, instanceID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            string
			resourceOwner string
			correctRoles  = database.TextArray[string]{}
		)
		if err := rows.Scan(&id, &resourceOwner, &correctRoles); err != nil {
			return 0, err
		}
		aggregate := &eventstore.Aggregate{
			ID:            id,
			Type:          usergrant.AggregateType,
			Version:       usergrant.AggregateVersion,
			ResourceOwner: resourceOwner,
			InstanceID:    instanceID,
		}
		correctedGrants = append(correctedGrants, usergrant.NewUserGrantCascadeChangedEvent(ctx, aggregate, correctRoles))
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	_, err = mig.eventstore.PushWithClient(ctx, tx, correctedGrants...)
	return len(correctedGrants), err
}

func (*FixUserGrantRoles) String() string {
	return "73_fix_user_grant_roles"
}
