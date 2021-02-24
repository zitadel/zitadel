package authz

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/query"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, authZ authz.Config, systemDefaults sd.SystemDefaults, queries *query.Queries) (*eventsourcing.EsRepository, error) {
	return eventsourcing.Start(config.Repository, authZ, systemDefaults, queries)
}
