package domain

import (
	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/event"
	"github.com/zitadel/zitadel/backend/repository/sql"
	"github.com/zitadel/zitadel/backend/repository/telemetry/logged"
	"github.com/zitadel/zitadel/backend/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

func (b *Instance) userCommandRepo(tx database.Transaction) repository.UserRepository {
	return logged.NewUser(
		b.logger,
		traced.NewUser(
			b.tracer,
			event.NewUser(
				eventstore.New(tx),
				sql.NewUser(tx),
			),
		),
	)
}

func (b *Instance) userQueryRepo(tx database.QueryExecutor) repository.UserRepository {
	return logged.NewUser(
		b.logger,
		traced.NewUser(
			b.tracer,
			sql.NewUser(tx),
		),
	)
}

type User struct {
	ID       string
	Username string
}
