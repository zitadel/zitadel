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

	SetUpStarted Step
	SetUpDone    Step
}

func NewAggregate(
	id,
	resourceOwner string,
	previousSequence uint64,
) *Aggregate {

	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(id, AggregateType, resourceOwner, AggregateVersion, previousSequence),
	}
}

func AggregateFromReadModel(rm *ReadModel) *Aggregate {
	return &Aggregate{
		Aggregate:    *eventstore.NewAggregate(rm.AggregateID, AggregateType, rm.ResourceOwner, AggregateVersion, rm.ProcessedSequence),
		SetUpDone:    rm.SetUpDone,
		SetUpStarted: rm.SetUpStarted,
	}
}

func (a *Aggregate) PushMemberAdded(ctx context.Context, userID string, roles ...string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberAddedEvent(ctx, userID, roles...))
	return a
}

func (a *Aggregate) PushMemberChanged(ctx context.Context, userID string, roles ...string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberChangedEvent(ctx, userID, roles...))
	return a
}

func (a *Aggregate) PushMemberRemoved(ctx context.Context, userID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberRemovedEvent(ctx, userID))
	return a
}
