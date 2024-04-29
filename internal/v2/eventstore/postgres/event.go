package postgres

import (
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func intentToCommands(intent *intent) (commands []*command, err error) {
	commands = make([]*command, len(intent.Commands()))

	for i, cmd := range intent.Commands() {
		var payload unmarshalPayload
		if cmd.Payload() != nil {
			payload, err = json.Marshal(cmd.Payload())
			if err != nil {
				logging.WithError(err).Warn("marshal payload failed")
				return nil, zerrors.ThrowInternal(err, "POSTG-MInPK", "Errors.Internal")
			}
		}

		commands[i] = &command{
			Event: &eventstore.Event[eventstore.StoragePayload]{
				Aggregate: *intent.Aggregate(),
				Creator:   cmd.Creator(),
				Revision:  cmd.Revision(),
				Type:      cmd.Type(),
				// always add at least 1 to the currently stored sequence
				Sequence: intent.sequence + uint32(i) + 1,
				Payload:  payload,
			},
			intent:            intent,
			uniqueConstraints: cmd.UniqueConstraints(),
		}
	}

	return commands, nil
}

type command struct {
	*eventstore.Event[eventstore.StoragePayload]

	intent            *intent
	uniqueConstraints []*eventstore.UniqueConstraint
}

var _ eventstore.StoragePayload = (unmarshalPayload)(nil)

type unmarshalPayload []byte

// Unmarshal implements eventstore.StoragePayload.
func (p unmarshalPayload) Unmarshal(ptr any) error {
	if len(p) == 0 {
		return nil
	}
	if err := json.Unmarshal(p, ptr); err != nil {
		return zerrors.ThrowInternal(err, "POSTG-u8qVo", "Errors.Internal")
	}

	return nil
}
