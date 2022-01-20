package command

import (
	"context"
	"time"

	"github.com/caos/oidc/pkg/oidc"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddPersonalAccessToken(ctx context.Context, userID, resourceOwner string, expirationDate time.Time, allowedUserType domain.UserType) (*domain.Token, string, error) {
	userWriteModel, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, "", err
	}
	if !isUserStateExists(userWriteModel.UserState) {
		return nil, "", errors.ThrowPreconditionFailed(nil, "COMMAND-Dggw2", "Errors.User.NotFound")
	}
	if allowedUserType != domain.UserTypeUnspecified && userWriteModel.UserType != allowedUserType {
		return nil, "", errors.ThrowPreconditionFailed(nil, "COMMAND-Df2f1", "Errors.User.WrongType")
	}
	tokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}
	tokenWriteModel := NewPersonalAccessTokenWriteModel(userID, tokenID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, tokenWriteModel)
	if err != nil {
		return nil, "", err
	}

	expirationDate, err = domain.ValidateExpirationDate(expirationDate)
	if err != nil {
		return nil, "", err
	}

	events, err := c.eventstore.Push(ctx,
		user.NewPersonalAccessTokenAddedEvent(
			ctx,
			UserAggregateFromWriteModel(&tokenWriteModel.WriteModel),
			tokenID,
			expirationDate,
			[]string{oidc.ScopeOpenID},
		),
	)
	if err != nil {
		return nil, "", err
	}
	err = AppendAndReduce(tokenWriteModel, events...)
	if err != nil {
		return nil, "", err
	}
	return personalTokenWriteModelToToken(tokenWriteModel, c.keyAlgorithm)
}

func (c *Commands) RemovePersonalAccessToken(ctx context.Context, userID, tokenID, resourceOwner string) (*domain.ObjectDetails, error) {
	tokenWriteModel, err := c.personalAccessTokenWriteModelByID(ctx, userID, tokenID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !tokenWriteModel.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.PAT.NotFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewPersonalAccessTokenRemovedEvent(ctx, UserAggregateFromWriteModel(&tokenWriteModel.WriteModel), tokenID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(tokenWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&tokenWriteModel.WriteModel), nil
}

func (c *Commands) personalAccessTokenWriteModelByID(ctx context.Context, userID, tokenID, resourceOwner string) (writeModel *PersonalAccessTokenWriteModel, err error) {
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-4n8vs", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewPersonalAccessTokenWriteModel(userID, tokenID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
