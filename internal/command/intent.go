package command

import (
	"context"
	"net/url"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

const (
	HandlerPrefix                 = ""
	EndpointLDAPLogin             = "/login/ldap"
	QueryAuthRequestID            = "authRequestID"
	EndpointExternalLoginCallback = "/login/externalidp/callback"
)

func (c *Commands) prepareCreateIntent(writeModel *IDPIntentWriteModel, idpID string, successURL, failureURL string) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.Invalid")
		}
		successURL, err := url.Parse(successURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.Invalid")
		}
		failureURL, err := url.Parse(failureURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			err = getIDPIntentWriteModel(ctx, writeModel, filter)
			if err != nil {
				return nil, err
			}
			exists, err := ExistsIDP(ctx, filter, idpID, writeModel.ResourceOwner)
			if !exists || err != nil {
				return nil, errors.ThrowPreconditionFailed(err, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting")
			}
			return []eventstore.Command{
				idpintent.NewStartedEvent(ctx, writeModel.aggregate, successURL, failureURL, idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) CreateIntent(ctx context.Context, idpID, successURL, failureURL string) (string, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	resourceOwner := authz.GetCtxData(ctx).OrgID
	writeModel := NewIDPIntentWriteModel(id, resourceOwner)
	if err != nil {
		return "", nil, err
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareCreateIntent(writeModel, idpID, successURL, failureURL))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
	//
	//identityProvider, err := s.query.IDPTemplateByID(ctx, false, req.IdpId, false)
	//if err != nil {
	//	return nil, err
	//}
	//baseURL := c.baseURL(ctx)
	//callbackURL := baseURL + EndpointExternalLoginCallback
	//
	//var provider idp.Provider
	//switch identityProvider.Type {
	//case domain.IDPTypeOAuth:
	//	provider, err = oauthProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeOIDC:
	//	provider, err = oidcProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeJWT:
	//	provider, err = jwtProvider(identityProvider, s.idpAlg)
	//case domain.IDPTypeAzureAD:
	//	provider, err = azureProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeGitHub:
	//	provider, err = githubProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeGitHubEnterprise:
	//	provider, err = githubEnterpriseProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeGitLab:
	//	provider, err = gitlabProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeGitLabSelfHosted:
	//	provider, err = gitlabSelfHostedProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeGoogle:
	//	provider, err = googleProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeLDAP:
	//	provider, err = ldapProvider(identityProvider, callbackURL, s.idpAlg)
	//case domain.IDPTypeUnspecified:
	//	fallthrough
	//default:
	//	return nil, errors.ThrowInvalidArgument(nil, "LOGIN-AShek", "Errors.ExternalIDP.IDPTypeNotImplemented")
	//}
	//if err != nil {
	//	return nil, err
	//}
	//intentID := gen()
	//session, err := provider.BeginAuth(ctx, intentID) //TODO generate state
	//if err != nil {
	//	return nil, err
	//}
}

func getIDPIntentWriteModel(ctx context.Context, writeModel *IDPIntentWriteModel, filter preparation.FilterToQueryReducer) error {
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	writeModel.AppendEvents(events...)
	return writeModel.Reduce()
}
