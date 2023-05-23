package command

import (
	"context"
	"net/url"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

func (c *Commands) prepareCreateIntent(writeModel *IDPIntentWriteModel, idpID string, successURL, failureURL string) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if idpID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.Invalid")
		}
		successURL, err := url.Parse(successURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.Invalid")
		}
		failureURL, err := url.Parse(failureURL)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.Invalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			err = getIDPIntentWriteModel(ctx, writeModel, filter)
			if err != nil {
				return nil, err
			}
			exists, err := ExistsIDP(ctx, filter, idpID, writeModel.ResourceOwner)
			if !exists || err != nil {
				return nil, errors.ThrowPreconditionFailed(err, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting")
			}
			return []eventstore.Command{
				idpintent.NewStartedEvent(ctx, writeModel.aggregate, successURL, failureURL, idpID),
			}, nil
		}, nil
	}
}

func (c *Commands) CreateIntent(ctx context.Context, idpID, successURL, failureURL string) (string, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	resourceOwner := authz.GetCtxData(ctx).OrgID
	writeModel := NewIDPIntentWriteModel(id, resourceOwner)
	if err != nil {
		return "", nil, err
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareCreateIntent(writeModel, idpID, successURL, failureURL))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(writeModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) AuthURLFromProvider(ctx context.Context, idpID, state, callbackURL string) (string, error) {
	writeModel, err := IDPProviderWriteModel(ctx, c.eventstore.Filter, idpID)
	if err != nil {
		return "", err
	}
	provider, err := writeModel.ToProvider(callbackURL, c.idpConfigEncryption)
	if err != nil {
		return "", err
	}
	session, err := provider.BeginAuth(ctx, state)
	if err != nil {
		return "", err
	}
	return session.GetAuthURL(), nil
}

func getIDPIntentWriteModel(ctx context.Context, writeModel *IDPIntentWriteModel, filter preparation.FilterToQueryReducer) error {
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	writeModel.AppendEvents(events...)
	return writeModel.Reduce()
}
