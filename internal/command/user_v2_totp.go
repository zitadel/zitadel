package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func (c *Commands) AddUserTOTP(ctx context.Context, userID, resourceowner string) (*domain.TOTP, error) {
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	prep, err := c.createHumanTOTP(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if err = c.pushAppendAndReduce(ctx, prep.wm, prep.cmds...); err != nil {
		return nil, err
	}
	return &domain.TOTP{
		ObjectDetails: writeModelToObjectDetails(&prep.wm.WriteModel),
		Secret:        prep.key.Secret(),
		URI:           prep.key.URL(),
	}, nil
}

func (c *Commands) CheckUserTOTP(ctx context.Context, userID, code, resourceOwner string) (*domain.ObjectDetails, error) {
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	return c.HumanCheckMFATOTPSetup(ctx, userID, code, "", resourceOwner)
}
