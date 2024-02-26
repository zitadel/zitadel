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

	position  float64
	inTxOrder uint32
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
		rm.position = event.Position()
		//TODO: rm.inTxOrder = event.InTxOrder()
	}

	if err := rm.InstanceCreatedMilestone.Reduce(events...); err != nil {
		return err
	}

	if err := rm.InstanceRemovedMilestone.Reduce(events...); err != nil {
		return err
	}

	if err := rm.AuthOnInstanceMilestone.Reduce(events...); err != nil {
		return err
	}

	if err := rm.AuthOnAppMilestone.Reduce(events...); err != nil {
		return err
	}

	if err := rm.ProjectCreatedMilestone.Reduce(events...); err != nil {
		return err
	}

	if err := rm.AppCreatedMilestone.Reduce(events...); err != nil {
		return err
	}

	return nil
}
