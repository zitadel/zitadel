package authz

import (
	"github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/query"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(config Config, systemDefaults sd.SystemDefaults, queries *query.Queries) (repository.Repository, error) {
	return eventsourcing.Start(config.Repository, systemDefaults, queries)
}
