package eventstore

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

var _ eventstore.Command = (*mockCommand)(nil)

type mockCommand struct {
	aggregate   *eventstore.Aggregate
	payload     any
	constraints []*eventstore.UniqueConstraint
}

// Aggregate implements [eventstore.Command]
func (m *mockCommand) Aggregate() *eventstore.Aggregate {
	return m.aggregate
}

// Creator implements [eventstore.Command]
func (m *mockCommand) Creator() string {
	return "creator"
}

// Revision implements [eventstore.Command]
func (m *mockCommand) Revision() uint16 {
	return 1
}

// Type implements [eventstore.Command]
func (m *mockCommand) Type() eventstore.EventType {
	return "event.type"
}

// Payload implements [eventstore.Command]
func (m *mockCommand) Payload() any {
	return m.payload
}

// UniqueConstraints implements [eventstore.Command]
func (m *mockCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return m.constraints
}

func (e *mockCommand) Fields() []*eventstore.FieldOperation {
	return nil
}

func mockEvent(aggregate *eventstore.Aggregate, sequence uint64, payload Payload) eventstore.Event {
	return &event{
		command: &command{
			InstanceID:    aggregate.InstanceID,
			AggregateType: string(aggregate.Type),
			AggregateID:   aggregate.ID,
			Owner:         aggregate.ResourceOwner,
			Creator:       "creator",
			Revision:      1,
			CommandType:   "event.type",
			Payload:       payload,
		},
		sequence: sequence,
	}
}

func mockAggregate(id string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            id,
		Type:          "type",
		ResourceOwner: "ro",
		InstanceID:    "instance",
		Version:       "v1",
	}
}

func mockAggregateWithInstance(id, instance string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            id,
		InstanceID:    instance,
		Type:          "type",
		ResourceOwner: "ro",
		Version:       "v1",
	}
}
