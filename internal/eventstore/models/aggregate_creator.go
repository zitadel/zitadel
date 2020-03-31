package models

import (
	"context"

	"github.com/caos/logging"
)

type AggregateCreator struct {
	serviceName string
}

func NewAggregateCreator(serviceName string) *AggregateCreator {
	return &AggregateCreator{serviceName: serviceName}
}

func (c *AggregateCreator) NewAggregate(ctx context.Context, id string, typ AggregateType, v Version, latestSequence uint64) (*Aggregate, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	return &Aggregate{
		ID:             id,
		Type:           typ,
		latestSequence: latestSequence,
		Version:        v,
		Events:         make([]*Event, 0, 2),
		editorOrg:      editorOrg(ctx),
		editorService:  c.serviceName,
		editorUser:     editorUser(ctx),
		resourceOwner:  resourceOwner(ctx),
	}, nil
}

func MustNewAggregate(id string, typ AggregateType, v Version, latestSequence uint64, events ...*Event) *Aggregate {
	c := NewAggregateCreator("svc")
	aggregate, err := c.NewAggregate(context.TODO(), id, typ, v, latestSequence)
	logging.Log("MODEL-10XZW").OnError(err).Fatal("unable to create aggregate")
	for _, event := range events {
		aggregate, err = aggregate.AppendEvent(event.Type, event.Data)
		logging.Log("MODEL-ASLX5").OnError(err).Fatal("unable to append event")
	}

	return aggregate
}

func editorUser(ctx context.Context) string {
	return "userID"
}

func editorOrg(ctx context.Context) string {
	return "orgID"
}

func resourceOwner(ctx context.Context) string {
	return "resourceOwner"
}
