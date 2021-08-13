package eventsourcing

import (
	"context"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/spooler"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	es2 "github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/query"
)

type Config struct {
	SearchLimit uint64
	Domain      string
	APIDomain   string
	Eventstore  v1.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
}

type EsRepository struct {
	spooler    *es_spol.Spooler
	Eventstore v1.Eventstore
	eventstore.UserRepo
	eventstore.AuthRequestRepo
	eventstore.TokenRepo
	eventstore.RefreshTokenRepo
	eventstore.KeyRepository
	eventstore.ApplicationRepo
	eventstore.UserSessionRepo
	eventstore.UserGrantRepo
	eventstore.OrgRepository
	eventstore.IAMRepository
	eventstore.FeaturesRepo
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, authZRepo *authz_repo.EsRepository, esV2 *es2.Eventstore) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	keyAlgorithm, err := crypto.NewAESCrypto(systemDefaults.KeyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	idGenerator := id.SonyFlakeGenerator

	assetsAPI := conf.APIDomain + "/assets/v1/"

	view, err := auth_view.StartView(sqlClient, keyAlgorithm, idGenerator, assetsAPI)
	if err != nil {
		return nil, err
	}

	authReq, err := cache.Start(conf.AuthRequest)
	if err != nil {
		return nil, err
	}

	statikLoginFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-20opp").OnError(err).Panic("unable to start login statik dir")

	keyChan := make(chan *key_model.KeyView)
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults, keyChan)
	locker := spooler.NewLocker(sqlClient)

	userRepo := eventstore.UserRepo{
		SearchLimit:     conf.SearchLimit,
		Eventstore:      es,
		View:            view,
		SystemDefaults:  systemDefaults,
		PrefixAvatarURL: assetsAPI,
	}
	return &EsRepository{
		spool,
		es,
		userRepo,
		eventstore.AuthRequestRepo{
			Command:                    command,
			AuthRequests:               authReq,
			View:                       view,
			Eventstore:                 es,
			UserSessionViewProvider:    view,
			UserViewProvider:           view,
			UserCommandProvider:        command,
			UserEventProvider:          &userRepo,
			OrgViewProvider:            view,
			IDPProviderViewProvider:    view,
			LoginPolicyViewProvider:    view,
			LockoutPolicyViewProvider:  view,
			UserGrantProvider:          view,
			IdGenerator:                idGenerator,
			PasswordCheckLifeTime:      systemDefaults.VerificationLifetimes.PasswordCheck.Duration,
			ExternalLoginCheckLifeTime: systemDefaults.VerificationLifetimes.PasswordCheck.Duration,
			MFAInitSkippedLifeTime:     systemDefaults.VerificationLifetimes.MFAInitSkip.Duration,
			SecondFactorCheckLifeTime:  systemDefaults.VerificationLifetimes.SecondFactorCheck.Duration,
			MultiFactorCheckLifeTime:   systemDefaults.VerificationLifetimes.MultiFactorCheck.Duration,
			IAMID:                      systemDefaults.IamID,
		},
		eventstore.TokenRepo{
			View:       view,
			Eventstore: es,
		},
		eventstore.RefreshTokenRepo{
			View:         view,
			Eventstore:   es,
			SearchLimit:  conf.SearchLimit,
			KeyAlgorithm: keyAlgorithm,
		},
		eventstore.KeyRepository{
			View:                     view,
			Commands:                 command,
			Eventstore:               esV2,
			SigningKeyRotationCheck:  systemDefaults.KeyConfig.SigningKeyRotationCheck.Duration,
			SigningKeyGracefulPeriod: systemDefaults.KeyConfig.SigningKeyGracefulPeriod.Duration,
			KeyAlgorithm:             keyAlgorithm,
			Locker:                   locker,
			KeyChan:                  keyChan,
		},
		eventstore.ApplicationRepo{
			Commands: command,
			View:     view,
		},

		eventstore.UserSessionRepo{
			View: view,
		},
		eventstore.UserGrantRepo{
			SearchLimit: conf.SearchLimit,
			View:        view,
			IamID:       systemDefaults.IamID,
			Auth:        authZ,
			AuthZRepo:   authZRepo,
		},
		eventstore.OrgRepository{
			SearchLimit:    conf.SearchLimit,
			View:           view,
			SystemDefaults: systemDefaults,
			Eventstore:     es,
		},
		eventstore.IAMRepository{
			IAMID:          systemDefaults.IamID,
			LoginDir:       statikLoginFS,
			IAMV2QuerySide: queries,
		},
		eventstore.FeaturesRepo{
			Eventstore: es,
			View:       view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}
