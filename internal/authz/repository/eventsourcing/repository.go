package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v1"

	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
)

type Config struct {
	Eventstore v1.Config
	View       types.SQL
	Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserGrantRepo
	eventstore.IamRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults, queries *query.Queries) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(sqlClient, idGenerator)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults)

	return &EsRepository{
		spool,
		eventstore.UserGrantRepo{
			View:       view,
			IamID:      systemDefaults.IamID,
			Auth:       authZ,
			Eventstore: es,
		},
		eventstore.IamRepo{
			IAMID:      systemDefaults.IamID,
			IAMV2Query: queries,
		},
		eventstore.TokenVerifierRepo{
			//TODO: Add Token Verification Key
			Eventstore: es,
			IAMID:      systemDefaults.IamID,
			View:       view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserGrantRepo.Health(); err != nil {
		return err
	}
	return nil
}
