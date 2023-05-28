package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

type aggregateOpt func(*Aggregate)

// NewAggregate is the default constructor of an aggregate
// opts overwrite values calculated by given parameters
func NewAggregate(
	ctx context.Context,
	id string,
	typ eventstore.AggregateType,
	version eventstore.Version,
	opts ...eventstore.AggregateOpt,
) *Aggregate {
	return eventstore.NewAggregate(ctx, id, typ, version, opts...)
}

// WithResourceOwner overwrites the resource owner of the aggregate
// by default the resource owner is set by the context
func WithResourceOwner(resourceOwner string) eventstore.AggregateOpt {
	return eventstore.WithResourceOwner(resourceOwner)
}

// AggregateFromWriteModel maps the given WriteModel to an Aggregate
func AggregateFromWriteModel(
	wm *WriteModel,
	typ AggregateType,
	version eventstore.Version,
) *Aggregate {
	return eventstore.NewAggregate(
		authz.WithInstanceID(context.Background(), wm.InstanceID),
		wm.AggregateID,
		typ,
		version,
		eventstore.WithResourceOwner(wm.ResourceOwner),
	)
}

// Aggregate is the basic implementation of Aggregater
type Aggregate = eventstore.Aggregate

func isAggreagteTypes(a *eventstore.Aggregate, types ...AggregateType) bool {
	for _, typ := range types {
		if a.Type == typ {
			return true
		}
	}
	return false
}

func isAggregateIDs(a *eventstore.Aggregate, ids ...string) bool {
	for _, id := range ids {
		if a.ID == id {
			return true
		}
	}
	return false
}
