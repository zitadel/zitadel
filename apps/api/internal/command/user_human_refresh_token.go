package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) RevokeRefreshToken(ctx context.Context, userID, orgID, tokenID string) (*domain.ObjectDetails, error) {
	removeEvent, refreshTokenWriteModel, err := c.removeRefreshToken(ctx, userID, orgID, tokenID)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(refreshTokenWriteModel, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&refreshTokenWriteModel.WriteModel), nil
}

func (c *Commands) RevokeRefreshTokens(ctx context.Context, userID, orgID string, tokenIDs []string) (err error) {
	if len(tokenIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Gfj42", "Errors.IDMissing")
	}
	events := make([]eventstore.Command, len(tokenIDs))
	for i, tokenID := range tokenIDs {
		event, _, err := c.removeRefreshToken(ctx, userID, orgID, tokenID)
		if err != nil {
			return err
		}
		events[i] = event
	}
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) removeRefreshToken(ctx context.Context, userID, orgID, tokenID string) (*user.HumanRefreshTokenRemovedEvent, *HumanRefreshTokenWriteModel, error) {
	if userID == "" || orgID == "" || tokenID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-GVDgf", "Errors.IDMissing")
	}
	refreshTokenWriteModel := NewHumanRefreshTokenWriteModel(userID, orgID, tokenID)
	err := c.eventstore.FilterToQueryReducer(ctx, refreshTokenWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenWriteModel.UserState != domain.UserStateActive {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-BHt2w", "Errors.User.RefreshToken.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return user.NewHumanRefreshTokenRemovedEvent(ctx, userAgg, tokenID), refreshTokenWriteModel, nil
}
