package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/spooler"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/id"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_user "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/v2/command"
	"github.com/caos/zitadel/internal/v2/query"
)

type Config struct {
	SearchLimit uint64
	Domain      string
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
	eventstore.KeyRepository
	eventstore.ApplicationRepo
	eventstore.UserSessionRepo
	eventstore.UserGrantRepo
	eventstore.OrgRepository
	eventstore.IAMRepository
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults, command *command.CommandSide, authZRepo *authz_repo.EsRepository) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	keyAlgorithm, err := crypto.NewAESCrypto(systemDefaults.KeyConfig.EncryptionConfig)
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

	iam, err := es_iam.StartIAM(
		es_iam.IAMConfig{
			Eventstore: es,
			Cache:      conf.Eventstore.Cache,
		},
		systemDefaults,
	)
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

	iamV2Query, err := query.StartQuerySide(&query.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
	if err != nil {
		return nil, err
	}

	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es, IAMDomain: conf.Domain}, systemDefaults)

	repos := handler.EventstoreRepos{UserEvents: user, ProjectEvents: project, OrgEvents: org, IamEvents: iam}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, repos, systemDefaults)

	return &EsRepository{
		spool,
		eventstore.UserRepo{
			SearchLimit:    conf.SearchLimit,
			Eventstore:     es,
			UserEvents:     user,
			OrgEvents:      org,
			View:           view,
			SystemDefaults: systemDefaults,
		},
		eventstore.AuthRequestRepo{
			Command:                    command,
			UserEvents:                 user,
			OrgEvents:                  org,
			AuthRequests:               authReq,
			View:                       view,
			UserSessionViewProvider:    view,
			UserViewProvider:           view,
			UserCommandProvider:        command,
			UserEventProvider:          user,
			OrgViewProvider:            view,
			IDPProviderViewProvider:    view,
			LoginPolicyViewProvider:    view,
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
			UserEvents:    user,
			ProjectEvents: project,
			View:          view,
		},
		eventstore.KeyRepository{
			View:               view,
			SigningKeyRotation: systemDefaults.KeyConfig.SigningKeyRotation.Duration,
		},
		eventstore.ApplicationRepo{
			Commands:      command,
			View:          view,
			ProjectEvents: project,
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
			OrgEventstore:  org,
			UserEventstore: user,
			SystemDefaults: systemDefaults,
		},
		eventstore.IAMRepository{
			IAMID:          systemDefaults.IamID,
			IAMV2QuerySide: iamV2Query,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}
