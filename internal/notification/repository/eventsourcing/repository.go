package eventsourcing

import (
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/v2/command"
	"net/http"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	"golang.org/x/text/language"
)

type Config struct {
	DefaultLanguage language.Tag
	Eventstore      es_int.Config
	View            types.SQL
	Spooler         spooler.SpoolerConfig
	Domain          string
}

type EsRepository struct {
	spooler *es_spol.Spooler
}

func Start(conf Config, dir http.FileSystem, systemDefaults sd.SystemDefaults, command *command.CommandSide) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := noti_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}

	translator, err := i18n.NewTranslator(dir, i18n.TranslatorConfig{DefaultLanguage: conf.DefaultLanguage})
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
	eventstoreRepos := handler.EventstoreRepos{IAMEvents: iam}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, command, eventstoreRepos, systemDefaults, translator, dir)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
