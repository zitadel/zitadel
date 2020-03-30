package models

import (
	"regexp"

	"github.com/caos/eventstore-lib/pkg/models"
	"github.com/caos/zitadel/internal/errors"
)

var _ models.Aggregate = (*Aggregate)(nil)

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

type Version string

type Aggregate struct {
	id             string
	typ            string
	events         []*Event
	latestSequence uint64
	version        Version
}

func NewAggregate(id, typ, version string, latestSequence uint64, events ...*Event) *Aggregate {
	return &Aggregate{id: id, typ: typ, events: events, latestSequence: latestSequence}
}

func (a *Aggregate) Type() string {
	return a.typ
}

func (a *Aggregate) ID() string {
	return a.id
}

func (a *Aggregate) Events() models.Events {
	events := make(Events, len(a.events))
	for idx, event := range a.events {
		events[idx] = event
	}

	return &events
}

func (a *Aggregate) LatestSequence() uint64 {
	return a.latestSequence
}

func (a *Aggregate) Validate() error {
	if a.id == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-FSjKV", "id not set")
	}
	if a.typ == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-aj4t2", "type not set")
	}
	if len(a.events) < 1 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-PupjX", "no events set")
	}
	return a.version.Validate()
}

func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}
