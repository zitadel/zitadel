package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
	PasswordSaltCost      int
	ClientSecretGenerator crypto.GeneratorConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	ProjectRepo
	OrgRepository
	OrgMemberRepository
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

	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore:            es,
		Cache:                 conf.Eventstore.Cache,
		PasswordSaltCost:      conf.PasswordSaltCost,
		ClientSecretGenerator: conf.ClientSecretGenerator,
	})
	if err != nil {
		return nil, err
	}
	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es})

	return &EsRepository{
		ProjectRepo{project},
		OrgRepository{org},
		OrgMemberRepository{org},
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.ProjectEvents.Health(context.Background())
}
