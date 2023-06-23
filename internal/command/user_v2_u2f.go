package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) RegisterUserU2F(ctx context.Context, userID, resourceOwner, rpID string) (*domain.WebAuthNRegistrationDetails, error) {
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	return c.registerUserU2F(ctx, userID, resourceOwner, rpID)
}

func (c *Commands) registerUserU2F(ctx context.Context, userID, resourceOwner, rpID string) (*domain.WebAuthNRegistrationDetails, error) {
	wm, userAgg, webAuthN, err := c.createUserU2F(ctx, userID, resourceOwner, rpID)
	if err != nil {
		return nil, err
	}
	return c.pushUserU2F(ctx, wm, userAgg, webAuthN)
}

func (c *Commands) createUserU2F(ctx context.Context, userID, resourceOwner, rpID string) (*HumanWebAuthNWriteModel, *eventstore.Aggregate, *domain.WebAuthNToken, error) {
	tokens, err := c.getHumanU2FTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, nil, err
	}
	return c.addHumanWebAuthN(ctx, userID, resourceOwner, rpID, tokens, domain.AuthenticatorAttachmentUnspecified, domain.UserVerificationRequirementRequired)
}

func (c *Commands) pushUserU2F(ctx context.Context, wm *HumanWebAuthNWriteModel, userAgg *eventstore.Aggregate, webAuthN *domain.WebAuthNToken) (*domain.WebAuthNRegistrationDetails, error) {
	cmd := user.NewHumanU2FAddedEvent(ctx, userAgg, wm.WebauthNTokenID, webAuthN.Challenge, webAuthN.RPID)
	err := c.pushAppendAndReduce(ctx, wm, cmd)
	if err != nil {
		return nil, err
	}
	return &domain.WebAuthNRegistrationDetails{
		ObjectDetails:                      writeModelToObjectDetails(&wm.WriteModel),
		ID:                                 wm.WebauthNTokenID,
		PublicKeyCredentialCreationOptions: webAuthN.CredentialCreationData,
	}, nil
}
