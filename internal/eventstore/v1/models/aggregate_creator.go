package models

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
)

type AggregateCreator struct {
	serviceName string
}

func NewAggregateCreator(serviceName string) *AggregateCreator {
	return &AggregateCreator{serviceName: serviceName}
}

type option func(*Aggregate)

func (c *AggregateCreator) NewAggregate(ctx context.Context, id string, typ AggregateType, version Version, previousSequence uint64, opts ...option) (*Aggregate, error) {
	ctxData := authz.GetCtxData(ctx)
	editorUser := ctxData.UserID
	resourceOwner := ctxData.OrgID

	aggregate := &Aggregate{
		ID:               id,
		typ:              typ,
		PreviousSequence: previousSequence,
		version:          version,
		Events:           make([]*Event, 0, 2),
		editorService:    c.serviceName,
		editorUser:       editorUser,
		resourceOwner:    resourceOwner,
	}

	for _, opt := range opts {
		opt(aggregate)
	}

	if err := aggregate.Validate(); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func OverwriteEditorUser(userID string) func(*Aggregate) {
	return func(a *Aggregate) {
		a.editorUser = userID
	}
}

func OverwriteResourceOwner(resourceOwner string) func(*Aggregate) {
	return func(a *Aggregate) {
		a.resourceOwner = resourceOwner
	}
}
