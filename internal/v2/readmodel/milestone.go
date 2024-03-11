package readmodel

import (
	"context"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/projection"
)

type Milestone struct {
	projection.InstanceCreatedMilestone
	projection.InstanceRemovedMilestone
	projection.AuthOnInstanceMilestone
	projection.AuthOnAppMilestone
	projection.ProjectCreatedMilestone
	projection.AppCreatedMilestone

	position eventstore.GlobalPosition
}

func NewMilestone(instance string) *Milestone {
	return &Milestone{
		InstanceCreatedMilestone: *projection.NewInstanceCreatedMilestone(),
		InstanceRemovedMilestone: *projection.NewInstanceRemovedMilestone(),
		AuthOnInstanceMilestone:  *projection.NewAuthOnInstanceMilestone(),
		AuthOnAppMilestone:       *projection.NewAuthOnAppMilestone(),
		ProjectCreatedMilestone:  *projection.NewProjectCreatedMilestone(),
		AppCreatedMilestone:      *projection.NewAppCreatedMilestone(),
	}
}

func (rm *Milestone) Filter(ctx context.Context) []*eventstore.Filter {
	return eventstore.MergeFilters(
		rm.InstanceCreatedMilestone.Filter,
		rm.InstanceRemovedMilestone.Filter,
		rm.AuthOnInstanceMilestone.Filter,
		rm.AuthOnAppMilestone.Filter,
		rm.ProjectCreatedMilestone.Filter,
		rm.AppCreatedMilestone.Filter,
	)
}

func (rm *Milestone) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		if err := rm.InstanceCreatedMilestone.Reduce(event); err != nil {
			return err
		}

		if err := rm.InstanceRemovedMilestone.Reduce(event); err != nil {
			return err
		}

		if err := rm.AuthOnInstanceMilestone.Reduce(event); err != nil {
			return err
		}

		if err := rm.AuthOnAppMilestone.Reduce(event); err != nil {
			return err
		}

		if err := rm.ProjectCreatedMilestone.Reduce(event); err != nil {
			return err
		}

		if err := rm.AppCreatedMilestone.Reduce(event); err != nil {
			return err
		}

		rm.position = event.Position()
	}

	return nil
}
