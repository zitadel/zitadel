package models

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

type Aggregate struct {
	ID             string
	Typ            string
	Events         []*Event
	LatestSequence uint64
	Version        version
}

func NewAggregate(id, typ string, v version, latestSequence uint64, events ...*Event) (*Aggregate, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := event.Validate(); err != nil {
			return nil, err
		}
	}
	return &Aggregate{
		ID:             id,
		Typ:            typ,
		Events:         events,
		LatestSequence: latestSequence,
		Version:        v,
	}, nil
}

func MustNewAggregate(id, typ string, v version, latestSequence uint64, events ...*Event) *Aggregate {
	aggregate, err := NewAggregate(id, typ, v, latestSequence, events...)
	logging.Log("MODEL-10XZW").OnError(err).Fatal("unable to create aggregate")
	return aggregate
}

func (a *Aggregate) Validate() error {
	if a.ID == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-FSjKV", "id not set")
	}
	if a.Typ == "" {
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
