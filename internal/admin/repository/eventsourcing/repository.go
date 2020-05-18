package eventsourcing

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/setup"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
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
	iam, err := es_iam.StartIam(es_iam.IamConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}

	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es})

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

	eventstoreRepos := setup.EventstoreRepos{OrgEvents: org, UserEvents: user, ProjectEvents: project, IamEvents: iam}
	err = setup.StartSetup(systemDefaults, eventstoreRepos).Execute()
	logging.Log("SERVE-k280HZ").OnError(err).Panic("failed to execute setup")

	return &EsRepository{
		OrgRepo: eventstore.OrgRepo{
			Eventstore:     es,
			OrgEventstore:  org,
			UserEventstore: user,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	err := repo.Eventstore.Health(ctx)
	if err != nil {
		return err
	}
	err = repo.UserEventstore.Health(ctx)
	if err != nil {
		return err
	}
	return repo.OrgEventstore.Health(ctx)
}
