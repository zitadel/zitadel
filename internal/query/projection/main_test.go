package projection

import (
	"database/sql"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	action_repo "github.com/zitadel/zitadel/internal/repository/action"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	key_repo "github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/org"
	proj_repo "github.com/zitadel/zitadel/internal/repository/project"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type expect func(mockRepository *mock.MockRepository)

func eventstoreExpect(t *testing.T, expects ...expect) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e(m)
	}
	es := eventstore.NewEventstore(m)
	iam_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	proj_repo.RegisterEventMappers(es)
	usergrant.RegisterEventMappers(es)
	key_repo.RegisterEventMappers(es)
	action_repo.RegisterEventMappers(es)
	return es
}

func eventPusherToEvents(eventsPushes ...eventstore.Command) []*repository.Event {
	events := make([]*repository.Event, len(eventsPushes))
	for i, event := range eventsPushes {
		data, err := eventstore.EventData(event)
		if err != nil {
			return nil
		}
		events[i] = &repository.Event{
			AggregateID:   event.Aggregate().ID,
			AggregateType: repository.AggregateType(event.Aggregate().Type),
			ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
			EditorService: event.EditorService(),
			EditorUser:    event.EditorUser(),
			Type:          repository.EventType(event.Type()),
			Version:       repository.Version(event.Aggregate().Version),
			Data:          data,
		}
	}
	return events
}

func expectFilter(events ...*repository.Event) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}

func eventFromEventPusher(event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:                            "",
		Sequence:                      0,
		PreviousAggregateSequence:     0,
		PreviousAggregateTypeSequence: 0,
		CreationDate:                  time.Time{},
		Type:                          repository.EventType(event.Type()),
		Data:                          data,
		EditorService:                 event.EditorService(),
		EditorUser:                    event.EditorUser(),
		Version:                       repository.Version(event.Aggregate().Version),
		AggregateID:                   event.Aggregate().ID,
		AggregateType:                 repository.AggregateType(event.Aggregate().Type),
		ResourceOwner:                 sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
	}
}

func eventFromEventPusherWithInstanceID(instanceID string, event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:                            "",
		Sequence:                      0,
		PreviousAggregateSequence:     0,
		PreviousAggregateTypeSequence: 0,
		CreationDate:                  time.Time{},
		Type:                          repository.EventType(event.Type()),
		Data:                          data,
		EditorService:                 event.EditorService(),
		EditorUser:                    event.EditorUser(),
		Version:                       repository.Version(event.Aggregate().Version),
		AggregateID:                   event.Aggregate().ID,
		AggregateType:                 repository.AggregateType(event.Aggregate().Type),
		ResourceOwner:                 sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		InstanceID:                    instanceID,
	}
}
