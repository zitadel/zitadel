package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	iamEventTypePrefix = eventstore.EventType("iam.")
)

const (
	AggregateType    = "iam"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(
	id,
	resourceOwner string,
	previousSequence uint64,
) *Aggregate {

	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			id,
			AggregateType,
			resourceOwner,
			AggregateVersion,
			previousSequence,
		),
	}
}

func AggregateFromWriteModel(wm *eventstore.WriteModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, AggregateType, AggregateVersion),
	}
}

func AggregateFromReadModel(rm *ReadModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			rm.AggregateID,
			AggregateType,
			rm.ResourceOwner,
			AggregateVersion,
			rm.ProcessedSequence,
		),
	}
}

func (a *Aggregate) PushMemberAdded(ctx context.Context, userID string, roles ...string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberAddedEvent(ctx, userID, roles...))
	return a
}

func (a *Aggregate) PushMemberChanged(ctx context.Context, changed *MemberWriteModel) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberChangedEvent(ctx, changed.UserID, changed.Roles...))
	return a
}

func (a *Aggregate) PushMemberChangedFromExisting(ctx context.Context, current *MemberWriteModel, roles ...string) *Aggregate {
	e, err := MemberChangedEventFromExisting(ctx, current, roles...)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushMemberRemoved(ctx context.Context, userID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberRemovedEvent(ctx, userID))
	return a
}

func (a *Aggregate) PushStepStarted(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepStartedEvent(ctx, step))
	return a
}

func (a *Aggregate) PushStepDone(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepDoneEvent(ctx, step))
	return a
}
