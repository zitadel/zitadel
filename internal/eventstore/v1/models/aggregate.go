package models

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
)

type AggregateType string

func (at AggregateType) String() string {
	return string(at)
}

type Aggregates []*Aggregate

type Aggregate struct {
	ID               string
	typ              AggregateType
	PreviousSequence uint64
	version          Version

	editorService string
	editorUser    string
	resourceOwner string
	Events        []*Event
	Precondition  *precondition
}

func (a *Aggregate) Type() AggregateType {
	return a.typ
}

type precondition struct {
	Query      *SearchQuery
	Validation func(...*Event) error
}

func (a *Aggregate) AppendEvent(typ EventType, payload interface{}) (*Aggregate, error) {
	if string(typ) == "" {
		return a, errors.ThrowInvalidArgument(nil, "MODEL-TGoCb", "no event type")
	}
	data, err := eventData(payload)
	if err != nil {
		return a, err
	}

	e := &Event{
		CreationDate:     time.Now(),
		Data:             data,
		Type:             typ,
		AggregateID:      a.ID,
		AggregateType:    a.typ,
		AggregateVersion: a.version,
		EditorService:    a.editorService,
		EditorUser:       a.editorUser,
		ResourceOwner:    a.resourceOwner,
	}

	a.Events = append(a.Events, e)
	return a, nil
}

func (a *Aggregate) SetPrecondition(query *SearchQuery, validateFunc func(...*Event) error) *Aggregate {
	a.Precondition = &precondition{Query: query, Validation: validateFunc}
	return a
}

func (a *Aggregate) Validate() error {
	if a == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-yi5AC", "aggregate is nil")
	}
	if a.ID == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-FSjKV", "id not set")
	}
	if string(a.typ) == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-aj4t2", "type not set")
	}
	if err := a.version.Validate(); err != nil {
		return errors.ThrowPreconditionFailed(err, "MODEL-PupjX", "invalid version")
	}

	if a.editorService == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-clYbY", "editor service not set")
	}
	if a.editorUser == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-Xcssi", "editor user not set")
	}
	if a.resourceOwner == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-eBYUW", "resource owner not set")
	}
	if a.Precondition != nil && (a.Precondition.Query == nil || a.Precondition.Validation == nil) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-EEUvA", "invalid precondition")
	}

	return nil
}
