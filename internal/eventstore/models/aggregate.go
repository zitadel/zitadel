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
	resourceOwner string
	Events        []*Event
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
		PreviousSequence: a.latestSequence,
		AggregateID:      a.id,
		AggregateType:    a.typ,
		AggregateVersion: a.version,
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
