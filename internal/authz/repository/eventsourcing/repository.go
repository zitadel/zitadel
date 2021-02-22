package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/v2/query"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/id"
)

type Config struct {
	Domain      string
	Eventstore  es_int.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserGrantRepo
	eventstore.IamRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(sqlClient, idGenerator)
	if err != nil {
		return nil, err
	}

	iamV2, err := query.StartQuerySide(&query.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
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
			IAMV2Query: iamV2,
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
