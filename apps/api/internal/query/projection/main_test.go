package projection

import (
	"database/sql"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
)

type expect func(mockRepository *mock.MockRepository)

func eventstoreExpect(t *testing.T, expects ...expect) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e(m)
	}
	es := eventstore.NewEventstore(
		&eventstore.Config{
			Querier: m.MockQuerier,
			Pusher:  m.MockPusher,
		},
	)
	return es
}

func expectFilter(events ...eventstore.Event) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}

func expectFilterError(err error) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEventsError(err)
	}
}

func eventFromEventPusher(event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:            "",
		Seq:           0,
		CreationDate:  time.Time{},
		Typ:           event.Type(),
		Data:          data,
		EditorUser:    event.Creator(),
		Version:       event.Aggregate().Version,
		AggregateID:   event.Aggregate().ID,
		AggregateType: event.Aggregate().Type,
		ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
	}
}
