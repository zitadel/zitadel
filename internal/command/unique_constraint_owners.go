package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var _ eventstore.QueryReducer = (*uniqueConstraintOwnersBackfillState)(nil)

func (c *Commands) isUniqueConstraintOwnerDeleteReady(ctx context.Context) (bool, error) {
	if c.ownerDeleteReadyCached.Load() {
		return true, nil
	}
	if c.ownerDeleteReady == nil {
		return false, nil
	}
	ready, err := c.ownerDeleteReady(ctx)
	if err != nil {
		return false, err
	}
	if ready {
		c.ownerDeleteReadyCached.Store(true)
	}
	return ready, nil
}

func (c *Commands) uniqueConstraintOwnersBackfillFinalized(ctx context.Context) (bool, error) {
	var state uniqueConstraintOwnersBackfillState
	if err := c.eventstore.FilterToQueryReducer(ctx, &state); err != nil {
		return false, err
	}
	return state.finalized, nil
}

type uniqueConstraintOwnersBackfillState struct {
	eventstore.ReadModel
	finalized bool
}

func (*uniqueConstraintOwnersBackfillState) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID("").
		AddQuery().
		AggregateTypes(migration.SystemAggregate).
		AggregateIDs(migration.SystemAggregateID).
		EventTypes(eventstore.EventType("system.migration.repeatable.done")).
		Builder()
}

func (s *uniqueConstraintOwnersBackfillState) Reduce() error {
	for _, event := range s.Events {
		if event.Type() != eventstore.EventType("system.migration.repeatable.done") {
			continue
		}
		step, ok := event.(*migration.SetupStep)
		if !ok || step.Name != eventstore.UniqueConstraintOwnersBackfillStep {
			continue
		}
		lastRun, _ := step.LastRun.(map[string]interface{})
		if lastRun == nil {
			s.finalized = false
			continue
		}
		finalized, _ := lastRun["finalized"].(bool)
		s.finalized = finalized
	}
	return s.ReadModel.Reduce()
}
