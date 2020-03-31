package eventstore

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/eventstore/models"
)

type AggregateCreator struct {
	serviceName string
}

func NewAggregateCreator(serviceName string) *AggregateCreator {
	return &AggregateCreator{serviceName: serviceName}
}

type AggregateType = models.AggregateType

type Aggregate struct {
	id             string
	typ            AggregateType
	events         []*Event
	latestSequence uint64
	version        Version
	editorService  string
	editorUser     string
	editorOrg      string
	resourceOwner  string
}

func (c *AggregateCreator) NewAggregate(ctx context.Context, id string, typ AggregateType, v Version, latestSequence uint64) (*Aggregate, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	return &Aggregate{
		id:             id,
		typ:            typ,
		latestSequence: latestSequence,
		version:        v,
		events:         make([]*Event, 0, 2),
		editorOrg:      editorOrg(ctx),
		editorService:  c.serviceName,
		editorUser:     editorUser(ctx),
		resourceOwner:  resourceOwner(ctx),
	}, nil
}

func (a *Aggregate) AppendEvent(typ EventType, payload interface{}) (*Aggregate, error) {
	data, err := eventData(payload)
	if err != nil {
		return a, nil
	}
	e := &Event{
		modifierService:  a.editorService,
		creationDate:     time.Now(),
		data:             data,
		modifierTenant:   a.editorOrg,
		modifierUser:     a.editorUser,
		resourceOwner:    a.resourceOwner,
		typ:              typ,
		aggregateID:      a.id,
		aggregateType:    a.typ,
		aggregateVersion: a.version,
		previousSequence: a.latestSequence,
	}

	a.events = append(a.events, e)
	return a, nil
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
