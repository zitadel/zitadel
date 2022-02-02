package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
)

type Config struct {
	Eventstore v1.Config
	View       types.SQL
	Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserMembershipRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults, queries *query.Queries, keyConfig *crypto.KeyConfig) (repository.Repository, error) {
	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(sqlClient, idGenerator, queries)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults)

	keyAlgorithm, err := crypto.NewAESCrypto(keyConfig)
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		spool,
		eventstore.UserMembershipRepo{
			View: view,
		},
		eventstore.TokenVerifierRepo{
			TokenVerificationKey: keyAlgorithm,
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
