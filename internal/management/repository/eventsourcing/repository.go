package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
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
	eventstore.IAMRepository
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

	statikNotificationFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start notification statik dir")

	return &EsRepository{
		spooler: spool,
		OrgRepository: eventstore.OrgRepository{
			SearchLimit:                         conf.SearchLimit,
			Eventstore:                          es,
			View:                                view,
			Roles:                               roles,
			SystemDefaults:                      systemDefaults,
			PrefixAvatarURL:                     assetsAPI,
			LoginDir:                            statikLoginFS,
			NotificationDir:                     statikNotificationFS,
			LoginTranslationFileContents:        make(map[string][]byte),
			NotificationTranslationFileContents: make(map[string][]byte),
			Query:                               queries,
		},
		ProjectRepo:   eventstore.ProjectRepo{es, conf.SearchLimit, view, roles, systemDefaults.IamID, assetsAPI, queries},
		UserRepo:      eventstore.UserRepo{es, conf.SearchLimit, view, queries, systemDefaults, assetsAPI},
		IAMRepository: eventstore.IAMRepository{IAMV2Query: queries},
		view:          view,
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.view.Health()
}
