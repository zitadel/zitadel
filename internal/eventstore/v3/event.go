package eventstore

import (
	"encoding/json"
	"time"
)

type action interface {
	Aggregate() *Aggregate

	// Creator is the userid of the user which created the action
	Creator() string
	// Type describes the action
	Type() EventType
	// Revision of the action
	Revision() uint16
}

type Command interface {
	action
	// Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Payload() any
	// UniqueConstraints should be added for unique attributes of an event, if nil constraints will not be checked
	UniqueConstraints() []*UniqueConstraint
}

type Event interface {
	action

	// Sequence of the event in the aggregate
	Sequence() uint64
	CreatedAt() time.Time

	// Unmarshal parses the payload and stores the result
	// in the value pointed to by ptr. If ptr is nil or not a pointer,
	// Unmarshal returns an error
	Unmarshal(ptr any) error

	// Deprecated: only use for migration
	DataAsBytes() []byte
}

// AggregateType is the object name
type AggregateType string

// EventType is the description of the change
type EventType string

var (
	_ Event = (*event)(nil)
)

type event struct {
	aggregate *Aggregate
	creator   string
	revision  uint16
	typ       EventType
	createdAt time.Time
	sequence  uint64
	payload   []byte
}

func commandToEvent(sequence *latestSequence, command Command) (_ *event, err error) {
	var payload Payload
	if command.Payload() != nil {
		payload, err = json.Marshal(command.Payload())
		if err != nil {
			return nil, err
		}
	}
	return &event{
		aggregate: sequence.aggregate,
		creator:   command.Creator(),
		revision:  command.Revision(),
		typ:       command.Type(),
		payload:   payload,
		sequence:  sequence.sequence,
	}, nil
}

// CreationDate implements [eventstore.Event]
func (e *event) CreationDate() time.Time {
	return e.CreatedAt()
}

// EditorUser implements [eventstore.Event]
func (e *event) EditorUser() string {
	return e.Creator()
}

// Aggregate implements [Event]
func (e *event) Aggregate() *Aggregate {
	return e.aggregate
}

// Creator implements [Event]
func (e *event) Creator() string {
	return e.creator
}

// Revision implements [Event]
func (e *event) Revision() uint16 {
	return e.revision
}

// Type implements [Event]
func (e *event) Type() EventType {
	return e.typ
}

// CreatedAt implements [Event]
func (e *event) CreatedAt() time.Time {
	return e.createdAt
}

// Sequence implements [Event]
func (e *event) Sequence() uint64 {
	return e.sequence
}

// Unmarshal implements [Event]
func (e *event) Unmarshal(ptr any) error {
	return json.Unmarshal(e.payload, ptr)
}

// DataAsBytes implements [Event]
func (e *event) DataAsBytes() []byte {
	return e.payload
}
