package eventsourcing

import (
	"context"

	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	authz_view "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type EsRepository struct {
	eventstore.UserMembershipRepo
	eventstore.TokenVerifierRepo
}

func Start(queries *query.Queries, dbClient *database.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, externalSecure bool) (repository.Repository, error) {
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
			ExternalSecure:       externalSecure,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.TokenVerifierRepo.Health(); err != nil {
		return err
	}
	return nil
}
