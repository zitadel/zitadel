package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/setup"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig

}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	eventstore.OrgRepo
	eventstore.UserRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	//view, sql, err := mgmt_view.StartView(conf.View)
	//if err != nil {
	//	return nil, err
	//}

	//conf.Spooler.View = view
	//conf.Spooler.EsClient = es.Client
	//conf.Spooler.SQL = sql
	//spool := spooler.StartSpooler(conf.Spooler)

	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es})

	user, err := es_usr.StartUser(es_usr.UserConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}

	eventstoreRepos := setup.EventstoreRepos{OrgEvents: org, UserEvents: user}
	setup := setup.StartSetup(systemDefaults, eventstoreRepos)
	err = setup.Execute()
	logging.Log("SERVE-k280HZ").OnError(err).Panic("failed to execute setup")

	return &EsRepository{
		eventstore.OrgRepo{org},
		eventstore.UserRepo{user},
	}, nil
}

func (repo *EsRepository) Health() error {
	// return repo.ProjectEvents.Health(context.Background())
	return nil
}
