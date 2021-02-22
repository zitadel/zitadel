package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type Config struct {
	SearchLimit uint64
	Eventstore  es_int.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	Domain      string
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.OrgRepo
	eventstore.IAMRepository
	eventstore.AdministratorRepo
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults, roles []string) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	iam, err := es_iam.StartIAM(es_iam.IAMConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
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

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, handler.EventstoreRepos{IamEvents: iam}, systemDefaults)

	return &EsRepository{
		spooler: spool,
		OrgRepo: eventstore.OrgRepo{
			Eventstore:     es,
			View:           view,
			SearchLimit:    conf.SearchLimit,
			SystemDefaults: systemDefaults,
		},
		IAMRepository: eventstore.IAMRepository{
			IAMEventstore:  iam,
			View:           view,
			SystemDefaults: systemDefaults,
			SearchLimit:    conf.SearchLimit,
			Roles:          roles,
		},
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
}
