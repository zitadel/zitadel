package command

import (
	"context"
	"encoding/base64"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

func (c *Commands) SucceedIDPIntent(ctx context.Context, writeModel *IDPIntentWriteModel, idpUser idp.User, userID string) (string, error) {
	token, err := c.idpConfigEncryption.Encrypt([]byte(writeModel.AggregateID))
	if err != nil {
		return "", err
	}
	cmd, err := idpintent.NewSucceededEvent(
		ctx,
		&idpintent.NewAggregate(writeModel.AggregateID, writeModel.ResourceOwner).Aggregate,
		idpUser,
		userID,
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

func (c *Commands) GetIntentWriteModel(ctx context.Context, id, resourceOwner string) (*IDPIntentWriteModel, error) {
	writeModel := NewIDPIntentWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, err
}
