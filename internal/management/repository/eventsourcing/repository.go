package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"

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
	APIDomain   string
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
	eventstore.FeaturesRepo
	view *mgmt_view.View
}

func Start(conf Config, systemDefaults sd.SystemDefaults, roles []string, queries *query.Queries, staticStorage static.Storage) (*EsRepository, error) {

	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := mgmt_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults, staticStorage)
	assetsAPI := conf.APIDomain + "/assets/v1/"

	statikLoginFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start login statik dir")

	return &EsRepository{
		spooler: spool,
		OrgRepository: eventstore.OrgRepository{
			conf.SearchLimit,
			es,
			view,
			roles,
			systemDefaults,
			assetsAPI,
			statikLoginFS},
		ProjectRepo:   eventstore.ProjectRepo{es, conf.SearchLimit, view, roles, systemDefaults.IamID, assetsAPI},
		UserRepo:      eventstore.UserRepo{es, conf.SearchLimit, view, systemDefaults, assetsAPI},
		UserGrantRepo: eventstore.UserGrantRepo{conf.SearchLimit, view, assetsAPI},
		IAMRepository: eventstore.IAMRepository{IAMV2Query: queries},
		FeaturesRepo:  eventstore.FeaturesRepo{es, view, conf.SearchLimit, systemDefaults},
		view:          view,
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.view.Health()
}
