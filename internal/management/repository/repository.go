package repository

import (
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type Repository struct {
	//spooler *es_spooler.Spooler
	eventsourcing.ProjectRepo
}

func Start(conf Config) (*Repository, error) {
	es := es_int.Start(conf.Eventstore)

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

	return &Repository{
		eventsourcing.ProjectRepo{project},
	}, nil
}

func (repo *Repository) Health() error {
	return nil
}
