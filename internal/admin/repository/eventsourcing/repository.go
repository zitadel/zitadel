package eventsourcing

import (
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	OrgRepo
}

func Start(conf Config) (*EsRepository, error) {
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

	org, err := es_org.StartOrg(es_org.OrgConfig{Eventstore: es})
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		OrgRepo{org},
	}, nil
}

func (repo *EsRepository) Health() error {
	// return repo.ProjectEvents.Health(context.Background())
	return nil
}
