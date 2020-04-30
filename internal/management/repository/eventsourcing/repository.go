package eventsourcing

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"

	es_int "github.com/caos/zitadel/internal/eventstore"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_grant "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	ProjectRepo
	UserRepo
	UserGrantRepo
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
	return &EsRepository{
		ProjectRepo{project},
		UserRepo{user},
		UserGrantRepo{usergrant},
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.ProjectEvents.Health(context.Background())
}
