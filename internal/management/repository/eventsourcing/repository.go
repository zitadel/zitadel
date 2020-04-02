package eventsourcing

import (
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	ProjectRepo
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

	project, err := es_proj.StartProject(es_proj.ProjectConfig{Eventstore: es})
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		ProjectRepo{project},
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
