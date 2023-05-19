package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/repository/intent"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
)

const (
	HandlerPrefix                 = ""
	EndpointLDAPLogin             = "/login/ldap"
	QueryAuthRequestID            = "authRequestID"
	EndpointExternalLoginCallback = "/login/externalidp/callback"
)

type AddIntent struct {
	IDPID      string
	SuccessURL string
	FailureURL string
}

func prepareCreateIntent(a *intent.Aggregate, addIntent *AddIntent) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if a.ResourceOwner == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x901m1n", "Errors.ResourceOwnerMissing")
		}
		if a.ID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-ap0mbs", "Errors.Intent.IDMissing")
		}
		if addIntent.IDPID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.Invalid")
		}
		if addIntent.SuccessURL == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.Invalid")
		}
		if addIntent.FailureURL == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getIntentWriteModel(ctx, a.ID, a.ResourceOwner, filter)
			if err != nil {
				return nil, err
			}

			exists, err := ExistsIDP(ctx, filter, writeModel.IDPID, writeModel.ResourceOwner)
			if !exists || err != nil {
				return nil, errors.ThrowPreconditionFailed(err, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting")
			}
			return []eventstore.Command{
				intent.NewIntentAddedEvent(ctx, &a.Aggregate, addIntent.IDPID, addIntent.SuccessURL, addIntent.FailureURL),
			}, nil
		}, nil
	}
}

func (c *Commands) CreateIntent(ctx context.Context) (*object.Details, error) {
	identityProvider, err := s.query.IDPTemplateByID(ctx, false, req.IdpId, false)
	if err != nil {
		return nil, err
	}
	baseURL := c.baseURL(ctx)
	callbackURL := baseURL + EndpointExternalLoginCallback

	var provider idp.Provider
	switch identityProvider.Type {
	case domain.IDPTypeOAuth:
		provider, err = oauthProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeOIDC:
		provider, err = oidcProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeJWT:
		provider, err = jwtProvider(identityProvider, s.idpAlg)
	case domain.IDPTypeAzureAD:
		provider, err = azureProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeGitHub:
		provider, err = githubProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeGitHubEnterprise:
		provider, err = githubEnterpriseProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeGitLab:
		provider, err = gitlabProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeGitLabSelfHosted:
		provider, err = gitlabSelfHostedProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeGoogle:
		provider, err = googleProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeLDAP:
		provider, err = ldapProvider(identityProvider, callbackURL, s.idpAlg)
	case domain.IDPTypeUnspecified:
		fallthrough
	default:
		return nil, errors.ThrowInvalidArgument(nil, "LOGIN-AShek", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}
	if err != nil {
		return nil, err
	}
	intentID := gen()
	session, err := provider.BeginAuth(ctx, intentID) //TODO generate state
	if err != nil {
		return nil, err
	}
}

func getIntentWriteModel(ctx context.Context, id, resourceOwner string, filter preparation.FilterToQueryReducer) (*IntentWriteModel, error) {
	writeModel := NewIntentWriteModel(id, resourceOwner)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}

func (c *Commands) baseURL(ctx context.Context) string {
	return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), c.externalSecure) + HandlerPrefix
}
