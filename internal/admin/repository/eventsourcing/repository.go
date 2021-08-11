package eventsourcing

import (
	"context"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/static"
)

type Config struct {
	SearchLimit uint64
	Eventstore  v1.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	Domain      string
	APIDomain   string
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.OrgRepo
	eventstore.IAMRepository
	eventstore.AdministratorRepo
	eventstore.FeaturesRepo
	eventstore.UserRepo
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults, command *command.Commands, static static.Storage, roles []string, localDevMode bool) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
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

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults, command, static, localDevMode)
	assetsAPI := conf.APIDomain + "/assets/v1/"

	statikLoginFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start login statik dir")

	statikNotificationFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start notification statik dir")

	return &EsRepository{
		spooler: spool,
		OrgRepo: eventstore.OrgRepo{
			Eventstore:     es,
			View:           view,
			SearchLimit:    conf.SearchLimit,
			SystemDefaults: systemDefaults,
		},
		IAMRepository: eventstore.IAMRepository{
			Eventstore:                          es,
			View:                                view,
			SystemDefaults:                      systemDefaults,
			SearchLimit:                         conf.SearchLimit,
			Roles:                               roles,
			PrefixAvatarURL:                     assetsAPI,
			LoginDir:                            statikLoginFS,
			NotificationDir:                     statikNotificationFS,
			LoginTranslationFileContents:        make(map[string][]byte),
			NotificationTranslationFileContents: make(map[string][]byte),
		},
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
		FeaturesRepo: eventstore.FeaturesRepo{
			Eventstore:     es,
			View:           view,
			SearchLimit:    conf.SearchLimit,
			SystemDefaults: systemDefaults,
		},
		UserRepo: eventstore.UserRepo{
			Eventstore:      es,
			View:            view,
			SearchLimit:     conf.SearchLimit,
			SystemDefaults:  systemDefaults,
			PrefixAvatarURL: assetsAPI,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return nil
}
