package login

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/ui/login/handler"
)

type Config struct {
	Handler handler.Config
}

func Start(config Config, command *command.CommandSide, query *query.QuerySide, authRepo *eventsourcing.EsRepository, systemdefaults systemdefaults.SystemDefaults, localDevMode bool) (*handler.Login, string) {
	return handler.CreateLogin(config.Handler, command, query, authRepo, systemdefaults, localDevMode)
}
