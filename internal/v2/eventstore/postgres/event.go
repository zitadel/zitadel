package postgres

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ eventstore.Event = (*event)(nil)

func intentToCommands(intent *intent) (commands []*command, err error) {
	commands = make([]*command, len(intent.Commands()))

	for i, cmd := range intent.Commands() {
		var payload []byte
		if cmd.Payload() != nil {
			payload, err = json.Marshal(cmd.Payload())
			if err != nil {
				logging.WithError(err).Warn("marshal payload failed")
				return nil, zerrors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
			}
		}

		commands[i] = &command{
			event: &event{
				aggregate: intent.Aggregate(),
				creator:   cmd.Creator(),
				revision:  cmd.Revision(),
				typ:       cmd.Type(),
				payload:   payload,
				// always add at least 1 to the currently stored sequence
				sequence: intent.sequence + uint32(i) + 1,
			},
			uniqueConstraints: cmd.UniqueConstraints(),
		}
	}

	return commands, nil
}

type command struct {
	*event

	uniqueConstraints []*eventstore.UniqueConstraint
}

type event struct {
	aggregate *eventstore.Aggregate
	creator   string
	revision  uint16
	typ       string
	createdAt time.Time
	sequence  uint32
	position  float64
	payload   []byte
}

// Aggregate implements [eventstore.Event].
func (e *event) Aggregate() *eventstore.Aggregate {
	return e.aggregate
}

// Creator implements [eventstore.Event].
func (e *event) Creator() string {
	return e.creator
}

// Revision implements [eventstore.Event].
func (e *event) Revision() uint16 {
	return e.revision
}

// Type implements [eventstore.Event].
func (e *event) Type() string {
	return e.typ
}

// CreatedAt implements [eventstore.Event].
func (e *event) CreatedAt() time.Time {
	return e.createdAt
}

// Sequence implements [eventstore.Event].
func (e *event) Sequence() uint32 {
	return e.sequence
}

// Sequence implements [eventstore.Event].
func (e *event) Position() float64 {
	return e.position
}

// Unmarshal implements [eventstore.Event].
func (e *event) Unmarshal(ptr any) error {
	if len(e.payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(e.payload, ptr); err != nil {
		return zerrors.ThrowInternal(err, "POSTG-u8qVo", "Errors.Internal")
	}

	return nil
}
