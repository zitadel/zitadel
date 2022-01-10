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

func (c *Commands) AddMachineToken(ctx context.Context, userID, resourceOwner string, expirationDate time.Time) (*domain.Token, string, error) {
	err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, "", err
	}
	tokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}
	tokenWriteModel := NewMachineTokenWriteModel(userID, tokenID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, tokenWriteModel)
	if err != nil {
		return nil, "", err
	}

	expirationDate, err = domain.ValidateExpirationDate(expirationDate)
	if err != nil {
		return nil, "", err
	}

	events, err := c.eventstore.Push(ctx,
		user.NewMachineTokenAddedEvent(
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
	return machineTokenWriteModelToToken(tokenWriteModel, c.keyAlgorithm)
}

func (c *Commands) RemoveMachineToken(ctx context.Context, userID, tokenID, resourceOwner string) (*domain.ObjectDetails, error) {
	keyWriteModel, err := c.machineTokenWriteModelByID(ctx, userID, tokenID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !keyWriteModel.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-4m77G", "Errors.User.Machine.Token.NotFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		user.NewMachineTokenRemovedEvent(ctx, UserAggregateFromWriteModel(&keyWriteModel.WriteModel), tokenID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(keyWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&keyWriteModel.WriteModel), nil
}

func (c *Commands) machineTokenWriteModelByID(ctx context.Context, userID, tokenID, resourceOwner string) (writeModel *MachineTokenWriteModel, err error) {
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-4n8vs", "Errors.User.UserIDMissing")
	}
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewMachineTokenWriteModel(userID, tokenID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
