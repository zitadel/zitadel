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
	"github.com/caos/zitadel/internal/crypto"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/id"
	es_key "github.com/caos/zitadel/internal/key/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_user "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	Eventstore  es_int.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	KeyConfig   es_key.KeyConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserRepo
	eventstore.AuthRequestRepo
	eventstore.TokenRepo
	eventstore.KeyRepository
	eventstore.ApplicationRepo
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

	keyAlgorithm, err := crypto.NewAESCrypto(conf.KeyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	idGenerator := id.SonyFlakeGenerator

	view, err := auth_view.StartView(sqlClient, keyAlgorithm, idGenerator)
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

	key, err := es_key.StartKey(es, conf.KeyConfig, keyAlgorithm, idGenerator)
	if err != nil {
		return nil, err
	}

	project, err := es_proj.StartProject(
		es_proj.ProjectConfig{
			Cache:      conf.Eventstore.Cache,
			Eventstore: es,
		},
		systemDefaults,
	)
	if err != nil {
		return nil, err
	}

	repos := handler.EventstoreRepos{UserEvents: user}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, repos)

	return &EsRepository{
		spool,
		eventstore.UserRepo{
			UserEvents: user,
			View:       view,
		},
		eventstore.AuthRequestRepo{
			UserEvents:               user,
			AuthRequests:             authReq,
			View:                     view,
			UserSessionViewProvider:  view,
			UserViewProvider:         view,
			UserEventProvider:        user,
			IdGenerator:              idGenerator,
			PasswordCheckLifeTime:    systemDefaults.VerificationLifetimes.PasswordCheck.Duration,
			MfaInitSkippedLifeTime:   systemDefaults.VerificationLifetimes.MfaInitSkip.Duration,
			MfaSoftwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaSoftwareCheck.Duration,
			MfaHardwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaHardwareCheck.Duration,
		},
		eventstore.TokenRepo{View: view},
		eventstore.KeyRepository{
			KeyEvents:          key,
			View:               view,
			SigningKeyRotation: conf.KeyConfig.SigningKeyRotation.Duration,
		},
		eventstore.ApplicationRepo{
			View:          view,
			ProjectEvents: project,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}
