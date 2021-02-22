package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/v2/query"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/id"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type Config struct {
	Domain      string
	Eventstore  es_int.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserGrantRepo
	eventstore.IamRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(sqlClient, idGenerator)
	if err != nil {
		return nil, err
	}

	iam, err := es_iam.StartIAM(es_iam.IAMConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	iamV2, err := query.StartQuerySide(&query.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
	if err != nil {
		return nil, err
	}

	repos := handler.EventstoreRepos{IAMEvents: iam, ProjectEvents: project}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, repos, systemDefaults)

	return &EsRepository{
		spool,
		eventstore.UserGrantRepo{
			View:      view,
			IamID:     systemDefaults.IamID,
			Auth:      authZ,
			IamEvents: iam,
		},
		eventstore.IamRepo{
			IAMID:      systemDefaults.IamID,
			IAMEvents:  iam,
			IAMV2Query: iamV2,
		},
		eventstore.TokenVerifierRepo{
			//TODO: Add Token Verification Key
			Eventstore:    es,
			IAMID:         systemDefaults.IamID,
			IAMEvents:     iam,
			ProjectEvents: project,
			View:          view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserGrantRepo.Health(); err != nil {
		return err
	}
	return nil
}
