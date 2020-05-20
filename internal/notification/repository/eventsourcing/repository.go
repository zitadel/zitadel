package eventsourcing

import (
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	View       types.SQL
	Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
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
	view, err := noti_view.StartView(sqlClient)
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
	eventstoreRepos := handler.EventstoreRepos{UserEvents: user}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, eventstoreRepos, systemDefaults)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
