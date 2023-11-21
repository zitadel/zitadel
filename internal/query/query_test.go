package query

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	action_repo "github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	key_repo "github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	proj_repo "github.com/zitadel/zitadel/internal/repository/project"
	quota_repo "github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
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
		iam_repo.RegisterEventMappers(es)
		org.RegisterEventMappers(es)
		usr_repo.RegisterEventMappers(es)
		proj_repo.RegisterEventMappers(es)
		usergrant.RegisterEventMappers(es)
		key_repo.RegisterEventMappers(es)
		action_repo.RegisterEventMappers(es)
		session.RegisterEventMappers(es)
		idpintent.RegisterEventMappers(es)
		authrequest.RegisterEventMappers(es)
		oidcsession.RegisterEventMappers(es)
		quota_repo.RegisterEventMappers(es)
		limits.RegisterEventMappers(es)
		feature.RegisterEventMappers(es)
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
