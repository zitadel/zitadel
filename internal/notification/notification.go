package notification

import (
	"context"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"
	"github.com/zitadel/zitadel/internal/command"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/notification/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/query"

	_ "github.com/zitadel/zitadel/internal/notification/statik"
)

type Config struct {
	APIDomain  string
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, hasStatics bool) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	apiDomain := config.APIDomain
	if !hasStatics {
		apiDomain = ""
	}
	_, err = eventsourcing.Start(config.Repository, statikFS, systemDefaults, command, queries, apiDomain)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
}
