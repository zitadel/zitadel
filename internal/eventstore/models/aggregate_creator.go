package models

import (
	"context"
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

func editorUser(ctx context.Context) string {
	return "userID"
}

func editorOrg(ctx context.Context) string {
	return "orgID"
}

func resourceOwner(ctx context.Context) string {
	return "resourceOwner"
}
