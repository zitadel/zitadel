package authz

import (
	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(config Config, systemDefaults sd.SystemDefaults, queries *query.Queries) (repository.Repository, error) {
	return eventsourcing.Start(config.Repository, systemDefaults, queries)
}
