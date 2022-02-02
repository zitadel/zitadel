package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/spooler"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_spol "github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
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
	eventstore.UserSessionRepo
	eventstore.OrgRepository
}

func Start(conf Config, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, keyConfig *crypto.KeyConfig) (*EsRepository, error) {
	es, err := v1.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	keyAlgorithm, err := crypto.NewAESCrypto(keyConfig)
	if err != nil {
		return nil, err
	}
	idGenerator := id.SonyFlakeGenerator

	assetsAPI := conf.APIDomain + "/assets/v1/"

	view, err := auth_view.StartView(sqlClient, keyAlgorithm, queries, idGenerator, assetsAPI)
	if err != nil {
		return nil, err
	}

	authReq, err := cache.Start(conf.AuthRequest)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, systemDefaults, queries)

	userRepo := eventstore.UserRepo{
		SearchLimit:     conf.SearchLimit,
		Eventstore:      es,
		View:            view,
		Query:           queries,
		SystemDefaults:  systemDefaults,
		PrefixAvatarURL: assetsAPI,
	}
	//TODO: remove as soon as possible
	queryView := queryViewWrapper{
		queries,
		view,
	}
	return &EsRepository{
		spool,
		es,
		userRepo,
		eventstore.AuthRequestRepo{
			PrivacyPolicyProvider:      queries,
			LabelPolicyProvider:        queries,
			Command:                    command,
			Query:                      queries,
			OrgViewProvider:            queries,
			AuthRequests:               authReq,
			View:                       view,
			Eventstore:                 es,
			UserSessionViewProvider:    view,
			UserViewProvider:           view,
			UserCommandProvider:        command,
			UserEventProvider:          &userRepo,
			IDPProviderViewProvider:    view,
			LockoutPolicyViewProvider:  queries,
			LoginPolicyViewProvider:    queries,
			UserGrantProvider:          queryView,
			ProjectProvider:            queryView,
			ApplicationProvider:        queries,
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
		eventstore.UserSessionRepo{
			View: view,
		},
		eventstore.OrgRepository{
			SearchLimit:    conf.SearchLimit,
			View:           view,
			SystemDefaults: systemDefaults,
			Eventstore:     es,
			Query:          queries,
		},
	}, nil
}

type queryViewWrapper struct {
	*query.Queries
	*auth_view.View
}

func (q queryViewWrapper) UserGrantsByProjectAndUserID(projectID, userID string) ([]*query.UserGrant, error) {
	userGrantProjectID, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	queries := &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantUserID, userGrantProjectID}}
	grants, err := q.Queries.UserGrants(context.TODO(), queries)
	if err != nil {
		return nil, err
	}
	return grants.UserGrants, nil
}
func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}
