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
	id             string
	typ            AggregateType
	latestSequence uint64
	version        Version

	editorService string
	editorUser    string
	editorOrg     string
	resourceOwner string
	Events        []*Event
	Appender      appender
}

type appender func(...*Event)

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
		PreviousSequence: a.latestSequence,
		AggregateID:      a.id,
		AggregateType:    a.typ,
		AggregateVersion: a.version,
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
	if a.id == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-FSjKV", "id not set")
	}
	if string(a.typ) == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-aj4t2", "type not set")
	}
	if len(a.Events) < 1 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-PupjX", "no events set")
	}
	if err := a.version.Validate(); err != nil {
		return errors.ThrowPreconditionFailed(err, "MODEL-PupjX", "invalid version")
	}
	for _, event := range a.Events {
		if err := event.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Aggregate) SetAppender(appendFn appender) *Aggregate {
	a.Appender = appendFn
	return a
}

func (a *Aggregate) OverwriteEditorOrg(orgID string) *Aggregate {
	a.editorOrg = orgID
	return a
}

func (a *Aggregate) OverwriteEditorUser(userID string) *Aggregate {
	a.editorUser = userID
	return a
}

func (a *Aggregate) OverwriteResourceOwner(resourceOwner string) *Aggregate {
	a.resourceOwner = resourceOwner
	return a
}
