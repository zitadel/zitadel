package login

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/login/handler"
)

type Config struct {
	Handler handler.Config
}

func Start(config Config, authRepo *eventsourcing.EsRepository, pathPrefix string) *handler.Login {
	return handler.CreateLogin(config.Handler, authRepo, pathPrefix)
}
