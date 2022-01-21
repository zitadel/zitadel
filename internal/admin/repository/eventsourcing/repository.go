package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/static"
)

type Config struct {
	SearchLimit uint64
	Eventstore  v1.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	Domain      string
	APIDomain   string
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.AdministratorRepo
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults, command *command.Commands, static static.Storage, localDevMode bool) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := admin_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults, command, static, localDevMode)

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
