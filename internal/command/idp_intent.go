package command

import (
	"context"
	"encoding/base64"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

func (c *Commands) SucceedIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, userID string, tokens *oidc.Tokens[*oidc.IDTokenClaims]) (string, error) {
	token, err := c.idpConfigEncryption.Encrypt([]byte(writeModel.AggregateID))
	if err != nil {
		return "", err
	}
	var idToken string
	var accessToken *crypto.CryptoValue
	if tokens != nil {
		accessToken, err = crypto.Encrypt([]byte(tokens.AccessToken), c.idpConfigEncryption)
		if err != nil {
			return "", err
		}
		idToken = tokens.IDToken
	}
	cmd, err := idpintent.NewSucceededEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		idpUser,
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

func (c *Commands) FailIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel) error {
	cmd := idpintent.NewFailedEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
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
