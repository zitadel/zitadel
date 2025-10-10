package command

import (
	"context"

	"github.com/zitadel/logging"

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
	endpoints []string,
) error {
	_, err := c.eventstore.Push(ctx, milestone.NewPushedEvent(ctx, milestone.NewInstanceAggregate(instanceID), msType, endpoints, c.externalDomain))
	return err
}

func setupInstanceCreatedMilestone(validations *[]preparation.Validation, instanceID string) {
	*validations = append(*validations, func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				milestone.NewReachedEvent(ctx, milestone.NewInstanceAggregate(instanceID), milestone.InstanceCreated),
			}, nil
		}, nil
	})
}

func (s *OIDCSessionEvents) SetMilestones(ctx context.Context, clientID string, isHuman bool) (postCommit func(ctx context.Context), err error) {
	postCommit = func(ctx context.Context) {}
	milestones, err := s.commands.GetMilestonesReached(ctx)
	if err != nil {
		return postCommit, err
	}

	instance := authz.GetInstance(ctx)
	aggregate := milestone.NewAggregate(ctx)
	var invalidate bool
	if !milestones.AuthenticationSucceededOnInstance {
		s.events = append(s.events, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance))
		invalidate = true
	}
	if !milestones.AuthenticationSucceededOnApplication && isHuman && clientID != instance.ConsoleClientID() {
		s.events = append(s.events, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication))
		invalidate = true
	}
	if invalidate {
		postCommit = s.commands.invalidateMilestoneCachePostCommit(instance.InstanceID())
	}
	return postCommit, nil
}

func (s *SAMLSessionEvents) SetMilestones(ctx context.Context) (postCommit func(ctx context.Context), err error) {
	postCommit = func(ctx context.Context) {}
	milestones, err := s.commands.GetMilestonesReached(ctx)
	if err != nil {
		return postCommit, err
	}

	instance := authz.GetInstance(ctx)
	aggregate := milestone.NewAggregate(ctx)
	var invalidate bool
	if !milestones.AuthenticationSucceededOnInstance {
		s.events = append(s.events, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance))
		invalidate = true
	}
	if !milestones.AuthenticationSucceededOnApplication {
		s.events = append(s.events, milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication))
		invalidate = true
	}
	if invalidate {
		postCommit = s.commands.invalidateMilestoneCachePostCommit(instance.InstanceID())
	}
	return postCommit, nil
}

func (c *Commands) projectCreatedMilestone(ctx context.Context, cmds *[]eventstore.Command) (postCommit func(ctx context.Context), err error) {
	postCommit = func(ctx context.Context) {}
	if isSystemUser(ctx) {
		return postCommit, nil
	}
	milestones, err := c.GetMilestonesReached(ctx)
	if err != nil {
		return postCommit, err
	}
	if milestones.ProjectCreated {
		return postCommit, nil
	}
	aggregate := milestone.NewAggregate(ctx)
	*cmds = append(*cmds, milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated))
	return c.invalidateMilestoneCachePostCommit(aggregate.InstanceID), nil
}

func (c *Commands) applicationCreatedMilestone(ctx context.Context, cmds *[]eventstore.Command) (postCommit func(ctx context.Context), err error) {
	postCommit = func(ctx context.Context) {}
	if isSystemUser(ctx) {
		return postCommit, nil
	}
	milestones, err := c.GetMilestonesReached(ctx)
	if err != nil {
		return postCommit, err
	}
	if milestones.ApplicationCreated {
		return postCommit, nil
	}
	aggregate := milestone.NewAggregate(ctx)
	*cmds = append(*cmds, milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated))
	return c.invalidateMilestoneCachePostCommit(aggregate.InstanceID), nil
}

func (c *Commands) invalidateMilestoneCachePostCommit(instanceID string) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := c.caches.milestones.Invalidate(ctx, milestoneIndexInstanceID, instanceID)
		logging.WithFields("instance_id", instanceID).OnError(err).Error("failed to invalidate milestone cache")
	}
}

func isSystemUser(ctx context.Context) bool {
	return authz.GetCtxData(ctx).SystemMemberships != nil
}
