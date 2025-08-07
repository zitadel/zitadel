package mirror

import (
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type succeededPayload struct {
	// Source is the name of the database data are mirrored from
	Source string `json:"source"`
	// Position until data will be mirrored
	Position decimal.Decimal `json:"position"`
}

const SucceededType = eventTypePrefix + "succeeded"

type SucceededEvent eventstore.Event[succeededPayload]

var _ eventstore.TypeChecker = (*SucceededEvent)(nil)

func (e *SucceededEvent) ActionType() string {
	return SucceededType
}

func SucceededEventFromStorage(event *eventstore.StorageEvent) (e *SucceededEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "MIRRO-xh5IW", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[succeededPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &SucceededEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

func NewSucceededCommand(source string, position decimal.Decimal) *eventstore.Command {
	return &eventstore.Command{
		Action: eventstore.Action[any]{
			Creator:  Creator,
			Type:     SucceededType,
			Revision: 1,
			Payload: succeededPayload{
				Source:   source,
				Position: position,
			},
		},
	}
}
