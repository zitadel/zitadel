package project

import (
	"context"
	"strings"
	"time"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func AddApp(a *project.Aggregate, id, name string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-0wTYg", "Errors.Invalid.Argument")
		}
		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-P7gKR", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if exists, err := ExistsProject(ctx, filter, a.ID, a.ResourceOwner); !exists || err != nil {
				return nil, errors.ThrowAlreadyExists(err, "PROJE-5LQ0U", "Errors.Project.NotFound")
			}
			return []eventstore.Command{
				project.NewApplicationAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					name,
				),
			}, nil
		}, nil
	}
}

func AddOIDCConfig(
	a project.Aggregate,
	version domain.OIDCVersion,
	appID string,
	clientID string,
	clientSecret *crypto.CryptoValue,
	redirectUris []string,
	responseTypes []domain.OIDCResponseType,
	grantTypes []domain.OIDCGrantType,
	applicationType domain.OIDCApplicationType,
	authMethodType domain.OIDCAuthMethodType,
	postLogoutRedirectUris []string,
	devMode bool,
	accessTokenType domain.OIDCTokenType,
	accessTokenRoleAssertion bool,
	idTokenRoleAssertion bool,
	idTokenUserinfoAssertion bool,
	clockSkew time.Duration,
	additionalOrigins []string,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if appID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument")
		}
		if clientID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-ghTsJ", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if exists, err := ExistsApp(ctx, filter, a.ID, appID, a.ResourceOwner); err != nil || !exists {
				return nil, errors.ThrowNotFound(err, "PROJE-sLDbG", "Errors.Project.Application.NotFound")
			}
			return []eventstore.Command{
				project.NewOIDCConfigAddedEvent(
					ctx,
					&a.Aggregate,
					version,
					appID,
					clientID,
					clientSecret,
					redirectUris,
					responseTypes,
					grantTypes,
					applicationType,
					authMethodType,
					postLogoutRedirectUris,
					devMode,
					accessTokenType,
					accessTokenRoleAssertion,
					idTokenRoleAssertion,
					idTokenUserinfoAssertion,
					clockSkew,
					additionalOrigins,
				),
			}, nil
		}, nil
	}
}

func AddAPIConfig(
	a project.Aggregate,
	appID string,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authMethodType domain.APIAuthMethodType,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if appID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument")
		}
		if clientID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-XXED5", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if exists, err := ExistsApp(ctx, filter, a.ID, appID, a.ResourceOwner); err != nil || !exists {
				return nil, errors.ThrowNotFound(err, "PROJE-EpG1p", "Errors.Project.Application.NotFound")
			}
			return []eventstore.Command{
				project.NewAPIConfigAddedEvent(
					ctx,
					&a.Aggregate,
					appID,
					clientID,
					clientSecret,
					authMethodType,
				),
			}, nil
		}, nil
	}
}

func ExistsApp(ctx context.Context, filter preparation.FilterToQueryReducer, projectID, appID, resourceOwner string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(resourceOwner).
		OrderAsc().
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(projectID).
		EventTypes(
			project.ApplicationAddedType,
			project.ApplicationRemovedType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationAddedEvent:
			if e.AppID == appID {
				exists = true
			}
		case *project.ApplicationRemovedEvent:
			if e.AppID == appID {
				exists = false
			}
		}
	}

	return exists, nil
}
