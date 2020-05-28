package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/spooler"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/id"
	es_policy "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	es_user "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	Eventstore  es_int.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserRepo
	eventstore.AuthRequestRepo
	eventstore.TokenRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := auth_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}
	policy, err := es_policy.StartPolicy(
		es_policy.PolicyConfig{
			Eventstore: es,
			Cache:      conf.Eventstore.Cache,
		},
		systemDefaults,
	)
	if err != nil {
		return nil, err
	}
	user, err := es_user.StartUser(
		es_user.UserConfig{
			Eventstore: es,
			Cache:      conf.Eventstore.Cache,
		},
		systemDefaults,
	)
	if err != nil {
		return nil, err
	}
	authReq, err := cache.Start(conf.AuthRequest)
	if err != nil {
		return nil, err
	}

	repos := handler.EventstoreRepos{UserEvents: user}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, repos)

	return &EsRepository{
		spool,
		eventstore.UserRepo{
			UserEvents:   user,
			PolicyEvents: policy,
			View:         view,
		},
		eventstore.AuthRequestRepo{
			UserEvents:               user,
			AuthRequests:             authReq,
			View:                     view,
			UserSessionViewProvider:  view,
			UserViewProvider:         view,
			IdGenerator:              id.SonyFlakeGenerator,
			PasswordCheckLifeTime:    systemDefaults.VerificationLifetimes.PasswordCheck.Duration,
			MfaInitSkippedLifeTime:   systemDefaults.VerificationLifetimes.MfaInitSkip.Duration,
			MfaSoftwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaSoftwareCheck.Duration,
			MfaHardwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaHardwareCheck.Duration,
		},
		eventstore.TokenRepo{View: view},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}
