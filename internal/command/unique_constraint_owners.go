package command

import (
	"context"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

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

func (c *Commands) uniqueConstraintOwnersBackfillMatchesVersion(ctx context.Context) (bool, error) {
	var states migration.StepStates
	if err := c.eventstore.FilterToQueryReducer(ctx, &states); err != nil {
		return false, err
	}
	for _, step := range states.Steps {
		if step == nil || step.Name != eventstore.UniqueConstraintOwnersBackfillStep {
			continue
		}
		lastRun, _ := step.LastRun.(map[string]interface{})
		if lastRun == nil {
			return false, nil
		}
		version, _ := lastRun["version"].(string)
		return version != "" && version == build.Version(), nil
	}
	return false, nil
}
