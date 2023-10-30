package eventstore

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	_ eventstore.Event = (*event)(nil)
)

type event struct {
	aggregate *eventstore.Aggregate
	creator   string
	revision  uint16
	typ       eventstore.EventType
	createdAt time.Time
	sequence  uint64
	position  float64
	payload   Payload
}

func commandToEvent(sequence *latestSequence, command eventstore.Command) (_ *event, err error) {
	var payload Payload
	if command.Payload() != nil {
		payload, err = json.Marshal(command.Payload())
		if err != nil {
			logging.WithError(err).Warn("marshal payload failed")
			return nil, errors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
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

// Aggregate implements [eventstore.Event]
func (e *event) Aggregate() *eventstore.Aggregate {
	return e.aggregate
}

// Creator implements [eventstore.Event]
func (e *event) Creator() string {
	return e.creator
}

// Revision implements [eventstore.Event]
func (e *event) Revision() uint16 {
	return e.revision
}

// Type implements [eventstore.Event]
func (e *event) Type() eventstore.EventType {
	return e.typ
}

// CreatedAt implements [eventstore.Event]
func (e *event) CreatedAt() time.Time {
	return e.createdAt
}

// Sequence implements [eventstore.Event]
func (e *event) Sequence() uint64 {
	return e.sequence
}

// Sequence implements [eventstore.Event]
func (e *event) Position() float64 {
	return e.position
}

// Unmarshal implements [eventstore.Event]
func (e *event) Unmarshal(ptr any) error {
	if len(e.payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(e.payload, ptr); err != nil {
		return errors.ThrowInternal(err, "V3-u8qVo", "Errors.Internal")
	}

	return nil
}

// DataAsBytes implements [eventstore.Event]
func (e *event) DataAsBytes() []byte {
	return e.payload
}
