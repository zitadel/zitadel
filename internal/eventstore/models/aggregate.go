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
	ID             string
	Type           AggregateType
	latestSequence uint64
	Version        Version

	editorService string
	editorUser    string
	editorOrg     string
	resourceOwner string
	Events        []*Event
	Appender      appender
}

type appender func(...*Event)

func (a *Aggregate) AppendEvent(typ EventType, payload interface{}) (*Aggregate, error) {
	data, err := eventData(payload)
	if err != nil {
		return a, nil
	}

	e := &Event{
		CreationDate:     time.Now(),
		Data:             data,
		Type:             typ,
		PreviousSequence: a.latestSequence,
		AggregateID:      a.ID,
		AggregateType:    a.Type,
		AggregateVersion: a.Version,
		EditorOrg:        a.editorOrg,
		EditorService:    a.editorService,
		EditorUser:       a.editorUser,
		ResourceOwner:    a.resourceOwner,
	}

	a.Events = append(a.Events, e)
	return a, nil
}

func (a *Aggregate) Validate() error {
	if a == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-yi5AC", "aggregate is nil")
	}
	if a.ID == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-FSjKV", "id not set")
	}
	if string(a.Type) == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-aj4t2", "type not set")
	}
	if len(a.Events) < 1 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-PupjX", "no events set")
	}
	for _, event := range a.Events {
		if err := event.Validate(); err != nil {
			return err
		}
	}
	return a.Version.Validate()
}

func (a *Aggregate) Appender(appendFn appender) *Aggregate {
	a.appender = appendFn
	return a
}
