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
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/repository/org"
	proj_repo "github.com/zitadel/zitadel/internal/repository/project"
	quota_repo "github.com/zitadel/zitadel/internal/repository/quota"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
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
	iam_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	proj_repo.RegisterEventMappers(es)
	quota_repo.RegisterEventMappers(es)
	limits.RegisterEventMappers(es)
	usergrant.RegisterEventMappers(es)
	key_repo.RegisterEventMappers(es)
	action_repo.RegisterEventMappers(es)
	return es
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
