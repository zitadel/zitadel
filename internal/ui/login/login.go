package login

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/ui/login/handler"
)

type Config struct {
	Handler handler.Config
}

func Start(config Config, command *command.Commands, query *query.Queries, authRepo *eventsourcing.EsRepository, staticStorage static.Storage, systemdefaults systemdefaults.SystemDefaults, localDevMode bool) (*handler.Login, string) {
	return handler.CreateLogin(config.Handler, command, query, authRepo, staticStorage, systemdefaults, localDevMode)
}
