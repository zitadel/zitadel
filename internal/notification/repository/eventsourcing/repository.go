package eventsourcing

import (
	"net/http"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/eventstore/v1"

	"golang.org/x/text/language"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
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

func Start(conf Config, dir http.FileSystem, systemDefaults sd.SystemDefaults, command *command.Commands, apiDomain string) (*EsRepository, error) {
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

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, command, systemDefaults, dir, apiDomain)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
