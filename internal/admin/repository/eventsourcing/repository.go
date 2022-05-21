package eventsourcing

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_spol "github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	SearchLimit uint64
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.AdministratorRepo
}

func Start(conf Config, static static.Storage, dbClient *sql.DB) (*EsRepository, error) {
	es, err := v1.Start(dbClient)
	if err != nil {
		return nil, err
	}
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, dbClient, static)

	return &EsRepository{
		spooler: spool,
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return nil
}
