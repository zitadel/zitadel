package authz

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
)

func Start(queries *query.Queries, dbClient *sql.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm) (repository.Repository, error) {
	return eventsourcing.Start(queries, dbClient, keyEncryptionAlgorithm)
}
