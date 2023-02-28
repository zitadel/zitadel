package authz

import (
	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query"
)

func Start(queries *query.Queries, dbClient *database.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, externalSecure bool) (repository.Repository, error) {
	return eventsourcing.Start(queries, dbClient, keyEncryptionAlgorithm, externalSecure)
}
