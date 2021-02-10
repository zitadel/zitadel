package usergrant

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	AggregateType    = "usergrant"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func (a *Aggregate) UserGrantAdded(
	ctx context.Context,
	userID,
	projectID,
	projectGrantID string,
	roleKeys []string,
) *UserGrantAddedEvent {
	return NewUserGrantAddedEvent(ctx, &a.Aggregate, userID, projectID, projectGrantID, roleKeys)
}

func (a *Aggregate) UserGrantChanged(
	ctx context.Context,
	roleKeys []string,
) *UserGrantChangedEvent {
	return NewUserGrantChangedEvent(ctx, &a.Aggregate, roleKeys)
}

func (a *Aggregate) UserGrantCascadeChanged(
	ctx context.Context,
	roleKeys []string,
) *UserGrantCascadeChangedEvent {
	return NewUserGrantCascadeChangedEvent(ctx, &a.Aggregate, roleKeys)
}

func (a *Aggregate) UserGrantRemoved(
	ctx context.Context,
	userID,
	projectID string,
) *UserGrantRemovedEvent {
	return NewUserGrantRemovedEvent(ctx, &a.Aggregate, userID, projectID)
}

func (a *Aggregate) UserGrantCascadeRemoved(
	ctx context.Context,
	userID,
	projectID string,
) *UserGrantCascadeRemovedEvent {
	return NewUserGrantCascadeRemovedEvent(ctx, &a.Aggregate, userID, projectID)
}

func (a *Aggregate) UserGrantDeactivatedEvent(ctx context.Context) *UserGrantDeactivatedEvent {
	return NewUserGrantDeactivatedEvent(ctx, &a.Aggregate)
}

func (a *Aggregate) UserGrantReactivatedEvent(ctx context.Context) *UserGrantReactivatedEvent {
	return NewUserGrantReactivatedEvent(ctx, &a.Aggregate)
}
