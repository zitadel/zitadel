package models

import (
	"time"

	"github.com/caos/eventstore-lib/pkg/models"
)

var _ models.Event = (*Event)(nil)

type Event struct {
	id               string
	creationDate     time.Time
	typ              string
	Sequence         uint64
	PreviousSequence uint64
	data             []byte
	ModifierService  string
	ModifierTenant   string
	ModifierUser     string
	ResourceOwner    string
	AggregateType    string
	AggregateID      string
}

func NewEvent(id string) *Event {
	return &Event{id: id}
}

func (e *Event) ID() string {
	return e.id
}

func (e *Event) CreationDate() time.Time {
	return e.creationDate
}

func (e *Event) SetCreationDate(creationDate time.Time) {
	e.creationDate = creationDate
}

func (e *Event) SetID(id string) {
	e.id = id
}

func (e *Event) Type() string {
	return e.typ
}

func (e *Event) SetType(typ string) {
	e.typ = typ
}

func (e *Event) Data() []byte {
	return e.data
}

func (e *Event) SetData(data []byte) {
	e.data = data
}
