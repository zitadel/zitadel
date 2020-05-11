package eventsourcing

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_grant "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
)

type Config struct {
	SearchLimit uint64
	Eventstore  es_int.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.ProjectRepo
	eventstore.UserRepo
	eventstore.UserGrantRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := mgmt_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}

	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	user, err := es_usr.StartUser(es_usr.UserConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	usergrant, err := es_grant.StartUserGrant(es_grant.UserGrantConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	})
	if err != nil {
		return nil, err
	}
	eventstoreRepos := handler.EventstoreRepos{ProjectEvents: project}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, eventstoreRepos)

	return &EsRepository{
		spool,
		eventstore.ProjectRepo{conf.SearchLimit, project, view},
		eventstore.UserRepo{conf.SearchLimit, user, view},
		eventstore.UserGrantRepo{conf.SearchLimit, usergrant, view},
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.ProjectEvents.Health(context.Background())
}
