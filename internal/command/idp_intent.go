package command

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

func (c *Commands) prepareCreateIntent(writeModel *IDPIntentWriteModel, idpID string, successURL, failureURL string) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.IDPMissing")
		}
		successURL, err := url.Parse(successURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.SuccessURLMissing")
		}
		failureURL, err := url.Parse(failureURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.FailureURLMissing")
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

func (c *Commands) CreateIntent(ctx context.Context, idpID, successURL, failureURL, resourceOwner string) (*IDPIntentWriteModel, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	writeModel := NewIDPIntentWriteModel(id, resourceOwner)
	if err != nil {
		return nil, nil, err
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareCreateIntent(writeModel, idpID, successURL, failureURL))
	if err != nil {
		return nil, nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return nil, nil, err
	}
	return writeModel, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) GetProvider(ctx context.Context, idpID string, callbackURL string) (idp.Provider, error) {
	writeModel, err := IDPProviderWriteModel(ctx, c.eventstore.Filter, idpID)
	if err != nil {
		return nil, err
	}
	return writeModel.ToProvider(callbackURL, c.idpConfigEncryption)
}

func (c *Commands) GetActiveIntent(ctx context.Context, intentID string) (*IDPIntentWriteModel, error) {
	intent, err := c.GetIntentWriteModel(ctx, intentID, "")
	if err != nil {
		return nil, err
	}
	if intent.State == domain.IDPIntentStateUnspecified {
		return nil, errors.ThrowNotFound(nil, "IDP-Hk38e", "Errors.Intent.NotStarted")
	}
	if intent.State != domain.IDPIntentStateStarted {
		return nil, errors.ThrowInvalidArgument(nil, "IDP-Sfrgs", "Errors.Intent.NotStarted")
	}
	return intent, nil
}

func (c *Commands) AuthURLFromProvider(ctx context.Context, idpID, state string, callbackURL string) (string, error) {
	provider, err := c.GetProvider(ctx, idpID, callbackURL)
	if err != nil {
		return "", err
	}
	session, err := provider.BeginAuth(ctx, state)
	if err != nil {
		return "", err
	}
	return session.GetAuthURL(), nil
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

func (c *Commands) SucceedIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, idpSession idp.Session, userID string) (string, error) {
	token, err := c.generateIntentToken(writeModel.AggregateID)
	if err != nil {
		return "", err
	}
	accessToken, idToken, err := tokensForSucceededIDPIntent(idpSession, c.idpConfigEncryption)
	if err != nil {
		return "", err
	}
	idpInfo, err := json.Marshal(idpUser)
	if err != nil {
		return "", err
	}
	cmd := idpintent.NewSucceededEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		idpInfo,
		idpUser.GetID(),
		idpUser.GetPreferredUsername(),
		userID,
		accessToken,
		idToken,
	)
	err = c.pushAppendAndReduce(ctx, writeModel, cmd)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *Commands) generateIntentToken(intentID string) (string, error) {
	token, err := c.idpConfigEncryption.Encrypt([]byte(intentID))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
}

func (c *Commands) SucceedLDAPIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, userID string, attributes map[string][]string) (string, error) {
	token, err := c.generateIntentToken(writeModel.AggregateID)
	if err != nil {
		return "", err
	}
	idpInfo, err := json.Marshal(idpUser)
	if err != nil {
		return "", err
	}
	cmd := idpintent.NewLDAPSucceededEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		idpInfo,
		idpUser.GetID(),
		idpUser.GetPreferredUsername(),
		userID,
		attributes,
	)
	err = c.pushAppendAndReduce(ctx, writeModel, cmd)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *Commands) FailIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, reason string) error {
	cmd := idpintent.NewFailedEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		reason,
	)
	_, err := c.eventstore.Push(ctx, cmd)
	return err
}

func (c *Commands) GetIntentWriteModel(ctx context.Context, id, resourceOwner string) (*IDPIntentWriteModel, error) {
	writeModel := NewIDPIntentWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, err
}

// tokensForSucceededIDPIntent extracts the oidc.Tokens if available (and encrypts the access_token) for the succeeded event payload
func tokensForSucceededIDPIntent(session idp.Session, encryptionAlg crypto.EncryptionAlgorithm) (*crypto.CryptoValue, string, error) {
	var tokens *oidc.Tokens[*oidc.IDTokenClaims]
	switch s := session.(type) {
	case *oauth.Session:
		tokens = s.Tokens
	case *openid.Session:
		tokens = s.Tokens
	case *jwt.Session:
		tokens = s.Tokens
	case *azuread.Session:
		tokens = s.Tokens
	default:
		return nil, "", nil
	}
	if tokens.Token == nil || tokens.AccessToken == "" {
		return nil, tokens.IDToken, nil
	}
	accessToken, err := crypto.Encrypt([]byte(tokens.AccessToken), encryptionAlg)
	return accessToken, tokens.IDToken, err
}
