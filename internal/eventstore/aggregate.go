package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
)

type aggregateOpt func(*Aggregate)

//NewAggregate is the default constructor of an aggregate
// opts overwrite values calculated by given parameters
func NewAggregate(
	ctx context.Context,
	id string,
	typ AggregateType,
	version Version,
	opts ...aggregateOpt,
) *Aggregate {
	a := &Aggregate{
		ID:            id,
		Type:          typ,
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		Tenant:        authz.GetCtxData(ctx).TenantID,
		Version:       version,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

//WithResourceOwner overwrites the resource owner of the aggregate
// by default the resource owner is set by the context
func WithResourceOwner(resourceOwner string) aggregateOpt {
	return func(aggregate *Aggregate) {
		aggregate.ResourceOwner = resourceOwner
	}
}

//AggregateFromWriteModel maps the given WriteModel to an Aggregate
func AggregateFromWriteModel(
	wm *WriteModel,
	typ AggregateType,
	version Version,
) *Aggregate {
	return &Aggregate{
		ID:            wm.AggregateID,
		Type:          typ,
		ResourceOwner: wm.ResourceOwner,
		Tenant:        wm.Tenant,
		Version:       version,
	}
}

//Aggregate is the basic implementation of Aggregater
type Aggregate struct {
	//ID is the unique identitfier of this aggregate
	ID string `json:"-"`
	//Type is the name of the aggregate.
	Type AggregateType `json:"-"`
	//ResourceOwner is the org this aggregates belongs to
	ResourceOwner string `json:"-"`
	//Tenant is the system this aggregate belongs to
	Tenant string `json:"-"`
	//Version is the semver this aggregate represents
	Version Version `json:"-"`
}

func isAggreagteTypes(a Aggregate, types ...AggregateType) bool {
	for _, typ := range types {
		if a.Type == typ {
			return true
		}
	}
	return false
}

func isAggregateIDs(a Aggregate, ids ...string) bool {
	for _, id := range ids {
		if a.ID == id {
			return true
		}
	}
	return false
}
