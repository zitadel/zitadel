package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
)

type Config struct {
	SearchLimit uint64
	Domain      string
	APIDomain   string
	Eventstore  v1.Config
	View        types.SQL
}

type EsRepository struct {
	eventstore.OrgRepository
	eventstore.ProjectRepo
	eventstore.UserRepo
	eventstore.IAMRepository
}

func Start(conf Config, systemDefaults sd.SystemDefaults, roles []string, queries *query.Queries, staticStorage static.Storage) (*EsRepository, error) {

	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	assetsAPI := conf.APIDomain + "/assets/v1/"

	statikLoginFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start login statik dir")

	statikNotificationFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start notification statik dir")

	return &EsRepository{
		OrgRepository: eventstore.OrgRepository{
			SearchLimit:                         conf.SearchLimit,
			Eventstore:                          es,
			Roles:                               roles,
			SystemDefaults:                      systemDefaults,
			PrefixAvatarURL:                     assetsAPI,
			LoginDir:                            statikLoginFS,
			NotificationDir:                     statikNotificationFS,
			LoginTranslationFileContents:        make(map[string][]byte),
			NotificationTranslationFileContents: make(map[string][]byte),
			Query:                               queries,
		},
		ProjectRepo:   eventstore.ProjectRepo{es, roles, systemDefaults.IamID, assetsAPI, queries},
		UserRepo:      eventstore.UserRepo{es, queries, systemDefaults, assetsAPI},
		IAMRepository: eventstore.IAMRepository{IAMV2Query: queries},
	}, nil
}
