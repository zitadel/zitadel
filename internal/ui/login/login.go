package login

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/ui/login/handler"
)

type Config struct {
	Handler handler.Config
}

func Start(config Config, authRepo *eventsourcing.EsRepository, localDevMode bool) (*handler.Login, string) {
	return handler.CreateLogin(config.Handler, authRepo, localDevMode)
}
