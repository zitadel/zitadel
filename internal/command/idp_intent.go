package command

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/apple"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) prepareCreateIntent(writeModel *IDPIntentWriteModel, idpID, successURL, failureURL string, idpArguments map[string]any) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if idpID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.IDPMissing")
		}
		successURL, err := url.Parse(successURL)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.SuccessURLMissing")
		}
		failureURL, err := url.Parse(failureURL)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.FailureURLMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			err = getIDPIntentWriteModel(ctx, writeModel, filter)
			if err != nil {
				return nil, err
			}
			exists, err := ExistsIDP(ctx, filter, idpID)
			if !exists || err != nil {
				return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting")
			}
			return []eventstore.Command{
				idpintent.NewStartedEvent(ctx,
					IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
					successURL,
					failureURL,
					idpID,
					idpArguments,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) CreateIntent(ctx context.Context, intentID, idpID, successURL, failureURL, resourceOwner string, idpArguments map[string]any) (*IDPIntentWriteModel, *domain.ObjectDetails, error) {
	if intentID == "" {
		var err error
		intentID, err = c.idGenerator.Next()
		if err != nil {
			return nil, nil, err
		}
	}
	writeModel := NewIDPIntentWriteModel(intentID, resourceOwner, c.maxIdPIntentLifetime)

	//nolint: staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareCreateIntent(writeModel, idpID, successURL, failureURL, idpArguments))
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

func (c *Commands) GetProvider(ctx context.Context, idpID string, idpCallback string, samlRootURL string) (idp.Provider, error) {
	writeModel, err := IDPProviderWriteModel(ctx, c.eventstore.Filter, idpID)
	if err != nil {
		return nil, err
	}
	if writeModel.IDPType != domain.IDPTypeSAML {
		return writeModel.ToProvider(idpCallback, c.idpConfigEncryption)
	}
	return writeModel.ToSAMLProvider(
		samlRootURL,
		c.idpConfigEncryption,
		func(ctx context.Context, intentID string) (*samlsp.TrackedRequest, error) {
			intent, err := c.GetActiveIntent(ctx, intentID)
			if err != nil {
				return nil, err
			}
			return &samlsp.TrackedRequest{
				SAMLRequestID: intent.RequestID,
				Index:         intentID,
				URI:           intent.SuccessURL.String(),
			}, nil
		},
		func(ctx context.Context, intentID, samlRequestID string) error {
			intent, err := c.GetActiveIntent(ctx, intentID)
			if err != nil {
				return err
			}
			return c.RequestSAMLIDPIntent(ctx, intent, samlRequestID)
		},
	)
}

func (c *Commands) GetActiveIntent(ctx context.Context, intentID string) (*IDPIntentWriteModel, error) {
	intent, err := c.GetIntentWriteModel(ctx, intentID, "")
	if err != nil {
		return nil, err
	}
	if intent.State == domain.IDPIntentStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "IDP-gy3ctgkqe7", "Errors.Intent.NotStarted")
	}
	if intent.State != domain.IDPIntentStateStarted {
		// we still need to return the intent to be able to redirect to the failure url
		return intent, zerrors.ThrowInvalidArgument(nil, "IDP-Sfrgs", "Errors.Intent.NotStarted")
	}
	return intent, nil
}

func (c *Commands) AuthFromProvider(ctx context.Context, idpID, idpCallback, samlRootURL string) (state string, session idp.Session, err error) {
	state, err = c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	provider, err := c.GetProvider(ctx, idpID, idpCallback, samlRootURL)
	if err != nil {
		return "", nil, err
	}
	session, err = provider.BeginAuth(ctx, state)
	return state, session, err
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
		IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
		idpInfo,
		idpUser.GetID(),
		idpUser.GetPreferredUsername(),
		userID,
		accessToken,
		idToken,
		idpSession.ExpiresAt(),
	)
	err = c.pushAppendAndReduce(ctx, writeModel, cmd)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *Commands) SucceedSAMLIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, userID string, session *saml.Session) (string, error) {
	token, err := c.generateIntentToken(writeModel.AggregateID)
	if err != nil {
		return "", err
	}
	idpInfo, err := json.Marshal(idpUser)
	if err != nil {
		return "", err
	}
	assertionData, err := xml.Marshal(session.Assertion)
	if err != nil {
		return "", err
	}
	assertionEnc, err := crypto.Encrypt(assertionData, c.idpConfigEncryption)
	if err != nil {
		return "", err
	}
	cmd := idpintent.NewSAMLSucceededEvent(
		ctx,
		IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
		idpInfo,
		idpUser.GetID(),
		idpUser.GetPreferredUsername(),
		userID,
		assertionEnc,
		session.ExpiresAt(),
	)
	err = c.pushAppendAndReduce(ctx, writeModel, cmd)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *Commands) RequestSAMLIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, requestID string) error {
	return c.pushAppendAndReduce(ctx, writeModel, idpintent.NewSAMLRequestEvent(
		ctx,
		IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
		requestID,
	))
}

func (c *Commands) generateIntentToken(intentID string) (string, error) {
	token, err := c.idpConfigEncryption.Encrypt([]byte(intentID))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
}

func (c *Commands) SucceedLDAPIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, userID string, session *ldap.Session) (string, error) {
	token, err := c.generateIntentToken(writeModel.AggregateID)
	if err != nil {
		return "", err
	}
	idpInfo, err := json.Marshal(idpUser)
	if err != nil {
		return "", err
	}
	attributes := make(map[string][]string, len(session.Entry.Attributes))
	for _, item := range session.Entry.Attributes {
		attributes[item.Name] = item.Values
	}
	cmd := idpintent.NewLDAPSucceededEvent(
		ctx,
		IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
		idpInfo,
		idpUser.GetID(),
		idpUser.GetPreferredUsername(),
		userID,
		attributes,
		session.ExpiresAt(),
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
		IDPIntentAggregateFromWriteModel(&writeModel.WriteModel),
		reason,
	)
	_, err := c.eventstore.Push(ctx, cmd)
	return err
}

func (c *Commands) GetIntentWriteModel(ctx context.Context, id, resourceOwner string) (*IDPIntentWriteModel, error) {
	writeModel := NewIDPIntentWriteModel(id, resourceOwner, c.maxIdPIntentLifetime)
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
		tokens = s.Tokens()
	case *apple.Session:
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
