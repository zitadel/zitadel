package eventsourcing

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	authz_view "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/crypto"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type EsRepository struct {
	eventstore.UserMembershipRepo
	eventstore.TokenVerifierRepo
}

func Start(queries *query.Queries, dbClient *sql.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm) (repository.Repository, error) {
	es, err := v1.Start(dbClient)
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator()
	view, err := authz_view.StartView(dbClient, idGenerator, queries)
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		eventstore.UserMembershipRepo{
			Queries: queries,
		},
		eventstore.TokenVerifierRepo{
			TokenVerificationKey: keyEncryptionAlgorithm,
			Eventstore:           es,
			View:                 view,
			Query:                queries,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.TokenVerifierRepo.Health(); err != nil {
		return err
	}
	return nil
}
