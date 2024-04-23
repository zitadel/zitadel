package projection

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

type Machine struct {
	projection

	id              string
	Name            string
	Description     string
	AccessTokenType domain.OIDCTokenType

	// TODO: separate projection?
	Secret *string
}

var _ eventstore.Reducer = (*Machine)(nil)

func NewMachineProjection(id string) *Machine {
	return &Machine{
		id: id,
	}
}

// Reduce implements eventstore.Reducer.
func (m *Machine) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !m.projection.shouldReduce(event) {
			continue
		}
		switch event.Type {
		case "user.machine.added":
			e, err := user.MachineAddedEventFromStorage(event)
			if err != nil {
				return err
			}

			m.Name = e.Payload.Name
			m.Description = e.Payload.Description
			m.AccessTokenType = e.Payload.AccessTokenType
		case "user.machine.changed":
			e, err := user.MachineChangedEventFromStorage(event)
			if err != nil {
				return err
			}

			if e.Payload.Name != nil {
				m.Name = *e.Payload.Name
			}
			if e.Payload.Description != nil {
				m.Description = *e.Payload.Description
			}
			if e.Payload.AccessTokenType != nil {
				m.AccessTokenType = *e.Payload.AccessTokenType
			}
		case "user.machine.secret.set":
			e, err := user.MachineSecretSetEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		case "user.machine.secret.updated":
			e, err := user.MachineSecretHashUpdatedEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		case "user.machine.secret.removed":
			e, err := user.MachineSecretHashUpdatedEventFromStorage(event)
			if err != nil {
				return err
			}
			m.Secret = &e.Payload.HashedSecret
		default:
			continue
		}
		m.projection.reduce(event)
	}
	return nil
}

func (m *Machine) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&m.position),
			),
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(m.id),
				eventstore.AppendEvent(
					eventstore.EventTypes(
						"user.machine.added",
						"user.machine.changed",
						"user.machine.secret.set",
						"user.machine.secret.updated",
						"user.machine.secret.removed",
					),
				),
			),
		),
	}
}
