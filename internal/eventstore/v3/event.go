package eventstore

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	_ eventstore.Event = (*event)(nil)
)

type command struct {
	InstanceID    string
	AggregateType string
	AggregateID   string
	CommandType   string
	Revision      uint16
	Payload       []byte
	Creator       string
	Owner         string
}

type event struct {
	aggregate *eventstore.Aggregate
	command   *command
	createdAt time.Time
	sequence  uint64
	position  float64
}

// func commandToEvent(sequence *latestSequence, cmd eventstore.Command) (_ *event, err error) {
// 	var payload Payload
// 	if cmd.Payload() != nil {
// 		payload, err = json.Marshal(cmd.Payload())
// 		if err != nil {
// 			logging.WithError(err).Warn("marshal payload failed")
// 			return nil, zerrors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
// 		}
// 	}
// 	return &event{
// 		aggregate: sequence.aggregate,
// 		creator:   cmd.Creator(),
// 		revision:  cmd.Revision(),
// 		typ:       cmd.Type(),
// 		payload:   payload,
// 		sequence:  sequence.sequence,
// 	}, nil
// }

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
	return e.command.Creator
}

// Revision implements [eventstore.Event]
func (e *event) Revision() uint16 {
	return e.command.Revision
}

// Type implements [eventstore.Event]
func (e *event) Type() eventstore.EventType {
	return eventstore.EventType(e.command.CommandType)
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
	if len(e.command.Payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(Payload(e.command.Payload), ptr); err != nil {
		return zerrors.ThrowInternal(err, "V3-u8qVo", "Errors.Internal")
	}

	return nil
}

// DataAsBytes implements [eventstore.Event]
func (e *event) DataAsBytes() []byte {
	return []byte(e.command.Payload)
}
