package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

type milestoneIndex int

const (
	milestoneIndexInstanceID milestoneIndex = iota
)

type MilestonesReached struct {
	InstanceID                           string
	InstanceCreated                      bool
	AuthenticationSucceededOnInstance    bool
	ProjectCreated                       bool
	ApplicationCreated                   bool
	AuthenticationSucceededOnApplication bool
	InstanceDeleted                      bool
}

// complete returns true if all milestones except InstanceDeleted are reached.
func (m *MilestonesReached) complete() bool {
	return m.InstanceCreated &&
		m.AuthenticationSucceededOnInstance &&
		m.ProjectCreated &&
		m.ApplicationCreated &&
		m.AuthenticationSucceededOnApplication
}

// GetMilestonesReached finds the milestone state for the current instance.
func (c *Commands) GetMilestonesReached(ctx context.Context) (*MilestonesReached, error) {
	milestones, ok := c.getCachedMilestonesReached(ctx)
	if ok {
		return milestones, nil
	}
	model := NewMilestonesReachedWriteModel(authz.GetInstance(ctx).InstanceID())
	if err := c.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	milestones = &model.MilestonesReached
	c.setCachedMilestonesReached(ctx, milestones)
	return milestones, nil
}

// getCachedMilestonesReached checks for milestone completeness on an instance and returns a filled
// [MilestonesReached] object.
// Otherwise it looks for the object in the milestone cache.
func (c *Commands) getCachedMilestonesReached(ctx context.Context) (*MilestonesReached, bool) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	if _, ok := c.milestonesCompleted.Load(instanceID); ok {
		return &MilestonesReached{
			InstanceID:                           instanceID,
			InstanceCreated:                      true,
			AuthenticationSucceededOnInstance:    true,
			ProjectCreated:                       true,
			ApplicationCreated:                   true,
			AuthenticationSucceededOnApplication: true,
			InstanceDeleted:                      false,
		}, ok
	}
	return c.caches.milestones.Get(ctx, milestoneIndexInstanceID, instanceID)
}

// setCachedMilestonesReached stores the current milestones state in the milestones cache.
// If the milestones are complete, the instance ID is stored in milestonesCompleted instead.
func (c *Commands) setCachedMilestonesReached(ctx context.Context, milestones *MilestonesReached) {
	if milestones.complete() {
		c.milestonesCompleted.Store(milestones.InstanceID, struct{}{})
		return
	}
	c.caches.milestones.Set(ctx, milestones)
}

// Keys implements cache.Entry
func (c *MilestonesReached) Keys(i milestoneIndex) []string {
	if i == milestoneIndexInstanceID {
		return []string{c.InstanceID}
	}
	return nil
}

// MilestonePushed writes a new milestone.PushedEvent with the milestone.Aggregate to the eventstore
func (c *Commands) MilestonePushed(
	ctx context.Context,
	instanceID string,
	msType milestone.Type,
	pushedDate time.Time,
	endpoints []string,
	primaryDomain string,
) error {
	_, err := c.eventstore.Push(ctx, milestone.NewPushedEvent(ctx, milestone.NewInstanceAggregate(instanceID), msType, pushedDate, endpoints, c.externalDomain, primaryDomain))
	return err
}

func setupInstanceCreatedMilestone(validations *[]preparation.Validation, instanceID string, reachedDate time.Time) {
	*validations = append(*validations, func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				milestone.NewReachedEvent(ctx, milestone.NewInstanceAggregate(instanceID), milestone.InstanceCreated, reachedDate),
			}, nil
		}, nil
	})
}

func (c *Commands) oidcSessionMilestones(ctx context.Context, clientID string, isHuman bool, reachedDate time.Time) error {
	milestones, err := c.GetMilestonesReached(ctx)
	if err != nil {
		return err
	}

	instance := authz.GetInstance(ctx)
	var cmds []eventstore.Command
	aggregate := milestone.NewAggregate(ctx)
	if !milestones.AuthenticationSucceededOnInstance {
		cmds = append(cmds, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, reachedDate))
	}
	if !milestones.AuthenticationSucceededOnApplication && isHuman && clientID != instance.ConsoleClientID() {
		cmds = append(cmds, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication, reachedDate))
	}
	if len(cmds) == 0 {
		return nil
	}
	if _, err = c.eventstore.Push(ctx, cmds...); err != nil {
		return err
	}
	return c.caches.milestones.Invalidate(ctx, milestoneIndexInstanceID, instance.InstanceID())
}

func (c *Commands) projectCreatedMilestone(ctx context.Context, reachedDate time.Time) error {
	if isSystemUser(ctx) {
		return nil
	}
	milestones, err := c.GetMilestonesReached(ctx)
	if err != nil {
		return err
	}
	if milestones.ProjectCreated {
		return nil
	}
	aggregate := milestone.NewAggregate(ctx)
	_, err = c.eventstore.Push(ctx, milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated, reachedDate))
	if err != nil {
		return err
	}
	return c.caches.milestones.Invalidate(ctx, milestoneIndexInstanceID, authz.GetInstance(ctx).InstanceID())
}

func (c *Commands) applicationCreatedMilestone(ctx context.Context, reachedDate time.Time) error {
	if isSystemUser(ctx) {
		return nil
	}
	milestones, err := c.GetMilestonesReached(ctx)
	if err != nil {
		return err
	}
	if milestones.ApplicationCreated {
		return nil
	}
	aggregate := milestone.NewAggregate(ctx)
	_, err = c.eventstore.Push(ctx, milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated, reachedDate))
	if err != nil {
		return err
	}
	return c.caches.milestones.Invalidate(ctx, milestoneIndexInstanceID, authz.GetInstance(ctx).InstanceID())
}

func (c *Commands) instanceRemovedMilestone(ctx context.Context, instanceID string, reachedDate time.Time) error {
	aggregate := milestone.NewInstanceAggregate(instanceID)
	_, err := c.eventstore.Push(ctx, milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceDeleted, reachedDate))
	if err != nil {
		return err
	}
	return c.caches.milestones.Invalidate(ctx, milestoneIndexInstanceID, instanceID)
}

func isSystemUser(ctx context.Context) bool {
	return authz.GetCtxData(ctx).SystemMemberships != nil
}
