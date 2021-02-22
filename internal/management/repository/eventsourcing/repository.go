package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/v2/query"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
)

type Config struct {
	SearchLimit uint64
	Domain      string
	Eventstore  es_int.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.OrgRepository
	eventstore.ProjectRepo
	eventstore.UserRepo
	eventstore.UserGrantRepo
	eventstore.IAMRepository
}

func Start(conf Config, systemDefaults sd.SystemDefaults, roles []string) (*EsRepository, error) {

	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := mgmt_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}

	iamV2Query, err := query.StartQuerySide(&query.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
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
	eventstoreRepos := handler.EventstoreRepos{IamEvents: iam}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, eventstoreRepos, systemDefaults)

	return &EsRepository{
		spooler:       spool,
		OrgRepository: eventstore.OrgRepository{conf.SearchLimit, es, iam, view, roles, systemDefaults},
		ProjectRepo:   eventstore.ProjectRepo{es, conf.SearchLimit, iam, view, roles, systemDefaults.IamID},
		UserRepo:      eventstore.UserRepo{es, conf.SearchLimit, view, systemDefaults},
		UserGrantRepo: eventstore.UserGrantRepo{conf.SearchLimit, view},
		IAMRepository: eventstore.IAMRepository{
			IAMV2Query: iamV2Query,
		},
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}

func (repo *EsRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	return repo.IAMRepository.IAMByID(ctx, id)
}
