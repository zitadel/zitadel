package eventsourcing

import (
	"context"
	es_policy "github.com/caos/zitadel/internal/policy/repository/eventsourcing"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/setup"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	SearchLimit uint64
	Eventstore  es_int.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	//TODO: should this be located in system-defaults?
	Domain string
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.OrgRepo
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	iam, err := es_iam.StartIam(es_iam.IamConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}

	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es, IAMDomain: conf.Domain}, systemDefaults)

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
	policy, err := es_policy.StartPolicy(es_policy.PolicyConfig{
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

	eventstoreRepos := setup.EventstoreRepos{OrgEvents: org, UserEvents: user, ProjectEvents: project, IamEvents: iam, PolicyEvents: policy}
	err = setup.StartSetup(systemDefaults, eventstoreRepos).Execute(ctx)
	logging.Log("SERVE-k280HZ").OnError(err).Panic("failed to execute setup")

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient)

	return &EsRepository{
		spooler: spool,
		OrgRepo: eventstore.OrgRepo{
			Eventstore:     es,
			OrgEventstore:  org,
			UserEventstore: user,
			View:           view,
			SearchLimit:    conf.SearchLimit,
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
