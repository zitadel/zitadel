package migrate

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	aggregateType     = eventstore.AggregateType("system")
	aggregateOwner    = "SYSTEM"
	aggregateInstance = ""
	eventCreator      = "COPY"

	startedType   = eventstore.EventType("system.copy.started")
	succeededType = eventstore.EventType("system.copy.succeeded")
	failedType    = eventstore.EventType("system.copy.failed")
)

func queryLastSuccessfulMigration(ctx context.Context, es *eventstore.Eventstore, destination string) (*lastSuccessfulMigration, error) {
	lastSuccess := &lastSuccessfulMigration{
		destination: destination,
	}
	if shouldIgnorePrevious {
		return lastSuccess, nil
	}
	err := es.FilterToReducer(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			EditorUser(eventCreator).
			InstanceID(aggregateInstance).
			OrderDesc().
			ResourceOwner(aggregateOwner).
			AddQuery().
			AggregateTypes(aggregateType).
			EventTypes(startedType, succeededType).
			Builder(),
		lastSuccess,
	)
	if err != nil {
		return nil, err
	}

	return lastSuccess, nil
}

type lastSuccessfulMigration struct {
	eventstore.ReadModel
	successID   string
	position    float64
	destination string
}

// Reduce implements eventstore.reducer.
func (m *lastSuccessfulMigration) Reduce() error {
	for _, event := range m.Events {
		if err := m.reduceEvent(event); err != nil {
			return err
		}
		if m.position > 0 {
			break
		}
	}
	return m.ReadModel.Reduce()
}

type destinationPayload struct {
	Destination string `json:"destination"`
}

func (m *lastSuccessfulMigration) reduceEvent(event eventstore.Event) error {
	payload := new(destinationPayload)
	if err := event.Unmarshal(payload); err != nil {
		return err
	}

	if m.destination != payload.Destination {
		return nil
	}

	switch event.Type() {
	case succeededType:
		if m.successID != "" {
			// there is a migration which succeeded later
			return nil
		}
		m.successID = event.Aggregate().ID
	case startedType:
		if event.Aggregate().ID != m.successID {
			return nil
		}
		m.position = event.Position()
	}
	return nil
}

func writeMigrationStart(ctx context.Context, es *eventstore.Eventstore, id string, destination string) (position float64, err error) {
	events, err := es.Push(ctx, &migrationStarted{
		id:          id,
		Destination: destination,
		Instances:   instanceIDs,
		System:      system,
	})
	if err != nil {
		return 0, err
	}
	return events[0].Position(), nil
}

var _ eventstore.Command = (*migrationStarted)(nil)

type migrationStarted struct {
	id          string
	Destination string   `json:"destination"`
	Instances   []string `json:"instances,omitempty"`
	System      bool     `json:"system,omitempty"`
}

// Aggregate implements eventstore.Command.
func (m *migrationStarted) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            m.id,
		Type:          aggregateType,
		ResourceOwner: aggregateOwner,
		InstanceID:    aggregateInstance,
		Version:       "v1",
	}
}

// Creator implements eventstore.Command.
func (*migrationStarted) Creator() string {
	return eventCreator
}

// Payload implements eventstore.Command.
func (m *migrationStarted) Payload() any {
	return m
}

// Revision implements eventstore.Command.
func (*migrationStarted) Revision() uint16 {
	return 1
}

// Type implements eventstore.Command.
func (*migrationStarted) Type() eventstore.EventType {
	return startedType
}

// UniqueConstraints implements eventstore.Command.
func (*migrationStarted) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func writeMigrationDone(ctx context.Context, es *eventstore.Eventstore, id string, err error, destination string) error {
	_, err = es.Push(ctx, &migrationDone{id: id, err: err, Destination: destination})
	return err
}

var _ eventstore.Command = (*migrationDone)(nil)

type migrationDone struct {
	err         error
	id          string
	Destination string `json:"destination"`
}

// Aggregate implements eventstore.Command.
func (m *migrationDone) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            m.id,
		Type:          aggregateType,
		ResourceOwner: aggregateOwner,
		InstanceID:    aggregateInstance,
		Version:       "v1",
	}
}

// Creator implements eventstore.Command.
func (*migrationDone) Creator() string {
	return eventCreator
}

// Payload implements eventstore.Command.
func (m *migrationDone) Payload() any {
	return m
}

// Revision implements eventstore.Command.
func (*migrationDone) Revision() uint16 {
	return 1
}

// Type implements eventstore.Command.
func (m *migrationDone) Type() eventstore.EventType {
	if m.err != nil {
		return failedType
	}
	return succeededType
}

// UniqueConstraints implements eventstore.Command.
func (*migrationDone) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
