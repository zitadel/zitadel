package postgres

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func intentToCommands(intent *intent) (commands []*command, err error) {
	commands = make([]*command, len(intent.Commands()))

	for i, cmd := range intent.Commands() {
		commands[i] = &command{
			Command:  cmd,
			intent:   intent,
			sequence: intent.nextSequence(),
		}

		if reflect.ValueOf(cmd.Payload).IsZero() {
			continue
		}

		if commands[i].payload, err = json.Marshal(cmd.Payload); err != nil {
			logging.WithFields("type", cmd.Type).WithError(err).Debug("unable to marshal event payload")
			return nil, zerrors.ThrowInternal(err, "POSTG-MInPK", "Errors.Internal")
		}
	}

	return commands, nil
}

type command struct {
	eventstore.Command

	intent *intent

	payload   []byte
	position  eventstore.GlobalPosition
	createdAt time.Time
	sequence  uint32
}

func (cmd *command) toEvent() *eventstore.StorageEvent {
	return &eventstore.StorageEvent{
		Action: eventstore.Action[eventstore.Unmarshal]{
			Creator:  cmd.Creator,
			Type:     cmd.Type,
			Revision: cmd.Revision,
			Payload: func(ptr any) error {
				return json.Unmarshal(cmd.payload, ptr)
			},
		},
		Aggregate: *cmd.intent.Aggregate(),
		Sequence:  cmd.intent.sequence,
		Position:  cmd.position,
		CreatedAt: cmd.createdAt,
	}
}
