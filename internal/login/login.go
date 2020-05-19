package login

import (
	"context"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/login/handler"
)

type Config struct {
	handler.Config
}


func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults, authRepo *eventsourcing.EsRepository) {
	handler.StartLogin(ctx, config.Config, authRepo)
}
