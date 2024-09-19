package authz

import (
	"github.com/zitadel/zitadel/v2/internal/authz/repository"
	"github.com/zitadel/zitadel/v2/internal/authz/repository/eventsourcing"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/database"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/query"
)

func Start(queries *query.Queries, es *eventstore.Eventstore, dbClient *database.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, externalSecure bool) (repository.Repository, error) {
	return eventsourcing.Start(queries, es, dbClient, keyEncryptionAlgorithm, externalSecure)
}
