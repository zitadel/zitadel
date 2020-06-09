package authz

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, authZ auth.Config, systemDefaults sd.SystemDefaults) (*eventsourcing.EsRepository, error) {
	return eventsourcing.Start(config.Repository, authZ, systemDefaults)
}
