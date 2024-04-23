package eventsourcing

import (
	"context"

	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/eventstore"
	auth_handler "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/handler"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/auth_request/repository/cache"
	"github.com/zitadel/zitadel/internal/command"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	eventstore2 "github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	SearchLimit                uint64
	Spooler                    auth_handler.Config
	AmountOfCachedAuthRequests uint16
}

type EsRepository struct {
	eventstore.UserRepo
	eventstore.AuthRequestRepo
	eventstore.TokenRepo
	eventstore.RefreshTokenRepo
	eventstore.UserSessionRepo
	eventstore.OrgRepository
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults, command *command.Commands, queries *query.Queries, dbClient *database.DB, esV2 *eventstore2.Eventstore, oidcEncryption crypto.EncryptionAlgorithm, userEncryption crypto.EncryptionAlgorithm) (*EsRepository, error) {
	view, err := auth_view.StartView(dbClient, oidcEncryption, queries, esV2)
	if err != nil {
		return nil, err
	}

	auth_handler.Register(ctx, conf.Spooler, view, queries)
	auth_handler.Start(ctx)

	authReq := cache.Start(dbClient, conf.AmountOfCachedAuthRequests)

	userRepo := eventstore.UserRepo{
		SearchLimit:    conf.SearchLimit,
		Eventstore:     esV2,
		View:           view,
		Query:          queries,
		SystemDefaults: systemDefaults,
	}
	//TODO: remove as soon as possible
	queryView := queryViewWrapper{
		queries,
		view,
	}
	return &EsRepository{
		userRepo,
		eventstore.AuthRequestRepo{
			PrivacyPolicyProvider:     queries,
			LabelPolicyProvider:       queries,
			Command:                   command,
			Query:                     queries,
			OrgViewProvider:           queries,
			AuthRequests:              authReq,
			View:                      view,
			UserCodeAlg:               userEncryption,
			UserSessionViewProvider:   view,
			UserViewProvider:          view,
			UserCommandProvider:       command,
			UserEventProvider:         &userRepo,
			IDPProviderViewProvider:   queries,
			IDPUserLinksProvider:      queries,
			LockoutPolicyViewProvider: queries,
			LoginPolicyViewProvider:   queries,
			UserGrantProvider:         queryView,
			ProjectProvider:           queryView,
			ApplicationProvider:       queries,
			CustomTextProvider:        queries,
			IdGenerator:               id.SonyFlakeGenerator(),
		},
		eventstore.TokenRepo{
			View:       view,
			Eventstore: esV2,
		},
		eventstore.RefreshTokenRepo{
			View:         view,
			Eventstore:   esV2,
			SearchLimit:  conf.SearchLimit,
			KeyAlgorithm: oidcEncryption,
		},
		eventstore.UserSessionRepo{
			View: view,
		},
		eventstore.OrgRepository{
			SearchLimit:    conf.SearchLimit,
			View:           view,
			SystemDefaults: systemDefaults,
			Eventstore:     esV2,
			Query:          queries,
		},
	}, nil
}

type queryViewWrapper struct {
	*query.Queries
	*auth_view.View
}

func (q queryViewWrapper) UserGrantsByProjectAndUserID(ctx context.Context, projectID, userID string) ([]*query.UserGrant, error) {
	userGrantProjectID, err := query.NewUserGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	queries := &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantUserID, userGrantProjectID}}
	grants, err := q.Queries.UserGrants(ctx, queries, true)
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
