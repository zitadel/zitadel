package query

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
)

type expect func(mockRepository *mock.MockRepository)

func expectEventstore(expects ...expect) func(*testing.T) *eventstore.Eventstore {
	return func(t *testing.T) *eventstore.Eventstore {
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
		InstanceID:    event.Aggregate().InstanceID,
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
		Constraints:   event.UniqueConstraints(),
	}
}

func Test_cleanStaticQueries(t *testing.T) {
	query := `select
	foo,
	bar
from table;`
	want := "select foo, bar from table;"
	cleanStaticQueries(&query)
	assert.Equal(t, want, query)
}
