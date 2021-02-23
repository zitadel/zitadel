package eventsourcing

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/v2/command"
	"net/http"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	"golang.org/x/text/language"
)

type Config struct {
	DefaultLanguage language.Tag
	Eventstore      v1.Config
	View            types.SQL
	Spooler         spooler.SpoolerConfig
	Domain          string
}

type EsRepository struct {
	spooler *es_spol.Spooler
}

func Start(conf Config, dir http.FileSystem, systemDefaults sd.SystemDefaults, command *command.CommandSide) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
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
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, command, systemDefaults, translator, dir)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
