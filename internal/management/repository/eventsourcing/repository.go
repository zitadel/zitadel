package eventsourcing

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/v2/query"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
)

type Config struct {
	SearchLimit uint64
	Domain      string
	Eventstore  v1.Config
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
	view *mgmt_view.View
}

func Start(conf Config, systemDefaults sd.SystemDefaults, roles []string) (*EsRepository, error) {

	es, err := v1.Start(conf.Eventstore)
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
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults)

	return &EsRepository{
		spooler:       spool,
		OrgRepository: eventstore.OrgRepository{conf.SearchLimit, es, view, roles, systemDefaults},
		ProjectRepo:   eventstore.ProjectRepo{es, conf.SearchLimit, view, roles, systemDefaults.IamID},
		UserRepo:      eventstore.UserRepo{es, conf.SearchLimit, view, systemDefaults},
		UserGrantRepo: eventstore.UserGrantRepo{conf.SearchLimit, view},
		IAMRepository: eventstore.IAMRepository{
			IAMV2Query: iamV2Query,
		},
		view: view,
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.view.Health()
}
