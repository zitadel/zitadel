package eventsourcing

import (
	"context"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"

	es_int "github.com/caos/zitadel/internal/eventstore"
	es_pol "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
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
	PolicyRepo
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

	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	policy, err := es_pol.StartPolicy(es_pol.PolicyConfig{Eventstore: es, Cache: conf.Eventstore.Cache})
	if err != nil {
		return nil, err
	}
	policy, err := es_pol.StartPolicy(es_pol.PolicyConfig{Eventstore: es, Cache: conf.Eventstore.Cache})
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		ProjectRepo{project},
		PolicyRepo{policy},
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.ProjectEvents.Health(context.Background())
}
