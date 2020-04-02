package models

import (
	"database/sql"
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
		PreviousSequence: sql.NullInt64{Int64: int64(a.latestSequence), Valid: true},
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
	if err := a.version.Validate(); err != nil {
		return errors.ThrowPreconditionFailed(err, "MODEL-PupjX", "invalid version")
	}

	if a.editorOrg == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-di3x5", "editor org not set")
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

	return nil
}

func (a *Aggregate) SetAppender(appendFn appender) *Aggregate {
	a.Appender = appendFn
	return a
}
