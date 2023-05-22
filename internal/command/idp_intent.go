package command

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

func (c *Commands) SucceedIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, idpSession idp.Session, userID string) (string, error) {
	token, err := c.idpConfigEncryption.Encrypt([]byte(writeModel.AggregateID))
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
	cmd, err := idpintent.NewSucceededEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		idpInfo,
		userID,
		accessToken,
		idToken,
	)
	if err != nil {
		return "", err
	}
	_, err = c.eventstore.Push(ctx, cmd)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
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
	default:
		return nil, "", nil
	}
	if tokens.Token == nil || tokens.AccessToken == "" {
		return nil, tokens.IDToken, nil
	}
	accessToken, err := crypto.Encrypt([]byte(tokens.AccessToken), encryptionAlg)
	return accessToken, tokens.IDToken, err
}
