package eventsourcing

import (
	"database/sql"
	"net/http"

	"github.com/caos/zitadel/internal/command"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/query"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
)

type Config struct {
	Spooler spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
}

func Start(conf Config, dir http.FileSystem, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, dbClient *sql.DB, assetsPrefix string) (*EsRepository, error) {
	es, err := v1.Start(dbClient)
	if err != nil {
		return nil, err
	}

	view, err := noti_view.StartView(dbClient)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, dbClient, command, queries, systemDefaults, dir, assetsPrefix)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
