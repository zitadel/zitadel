package notification

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing"
	"github.com/caos/zitadel/internal/query"

	_ "github.com/caos/zitadel/internal/notification/statik"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(config Config, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, dbClient *sql.DB, assetsPrefix string) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	_, err = eventsourcing.Start(config.Repository, statikFS, systemDefaults, command, queries, dbClient, assetsPrefix)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
}
