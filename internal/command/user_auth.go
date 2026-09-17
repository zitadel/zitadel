package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) userStateForAuthentication(ctx context.Context, userID, resourceOwner, userErrorID, orgErrorID string) (_ *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userStateModel, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !userStateModel.UserState.IsEnabled() {
		return nil, zerrors.ThrowPreconditionFailed(nil, userErrorID, "Errors.User.NotActive")
	}
	if err := c.checkOrgActiveForAuthentication(ctx, resourceOwner, orgErrorID); err != nil {
		return nil, err
	}
	return userStateModel, nil
}

func (c *Commands) checkOrgActiveForAuthentication(ctx context.Context, orgID, errorID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	orgWriteModel, err := c.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if !orgWriteModel.State.IsActive() {
		return zerrors.ThrowPreconditionFailed(nil, errorID, "Errors.User.NotActive")
	}
	return nil
}
