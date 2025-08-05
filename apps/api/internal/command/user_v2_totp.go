package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
)

func (c *Commands) AddUserTOTP(ctx context.Context, userID, resourceOwner string) (*domain.TOTP, error) {
	prep, err := c.createHumanTOTP(ctx, userID, resourceOwner)
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
	return c.HumanCheckMFATOTPSetup(ctx, userID, code, "", resourceOwner)
}
