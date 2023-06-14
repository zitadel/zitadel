package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func (c *Commands) AddUserOTP(ctx context.Context, userID, resourceowner string) (*domain.OTPv2, error) {
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	prep, err := c.createHumanOTP(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if err = c.pushAppendAndReduce(ctx, prep.wm, prep.cmds...); err != nil {
		return nil, err
	}
	return &domain.OTPv2{
		ObjectDetails: writeModelToObjectDetails(&prep.wm.WriteModel),
		Secret:        prep.key.Secret(),
		URI:           prep.key.URL(),
	}, nil
}
