package eventsourcing

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/authz/repository"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_spol "github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Spooler spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserMembershipRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults, queries *query.Queries, dbClient *sql.DB, keyEncryptionAlgorithm crypto.EncryptionAlgorithm) (repository.Repository, error) {
	es, err := v1.Start(dbClient)
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(dbClient, idGenerator, queries)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, dbClient, systemDefaults)

	return &EsRepository{
		spool,
		eventstore.UserMembershipRepo{
			View: view,
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
	if err := repo.UserMembershipRepo.Health(); err != nil {
		return err
	}
	return nil
}
