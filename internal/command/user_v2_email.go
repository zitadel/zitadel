package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserEmail struct {
	eventstore *eventstore.Eventstore
	aggregate  *eventstore.Aggregate
	events     []eventstore.Command
	model      *HumanEmailWriteModel
}

func (c *Commands) UserEmail(ctx context.Context, userID, resourceOwner string) (*UserEmail, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing")
	}

	model, err := c.emailWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if model.UserState == domain.UserStateUnspecified || model.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-J8dsk", "Errors.User.NotInitialised")
	}
	return &UserEmail{
		eventstore: c.eventstore,
		aggregate:  UserAggregateFromWriteModel(&model.WriteModel),
		model:      model,
	}, nil
}

func (c *UserEmail) Change(ctx context.Context, email domain.EmailAddress) error {
	if err := email.Validate(); err != nil {
		return err
	}
	event, hasChanged := c.model.NewChangedEvent(ctx, c.aggregate, email)
	if !hasChanged {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2b7fM", "Errors.User.Email.NotChanged")
	}
	c.events = append(c.events, event)
	return nil
}

func (c *UserEmail) SetVerified(ctx context.Context) {
	c.events = append(c.events, user.NewHumanEmailVerifiedEvent(ctx, c.aggregate))
}

func (c *UserEmail) AddCode(ctx context.Context, code *domain.EmailCode, urlTmpl *string) {
	c.events = append(c.events, user.NewHumanEmailCodeAddedEventV2(ctx, c.aggregate, code.Code, code.Expiry, urlTmpl))
}

func (c *UserEmail) Push(ctx context.Context) (*domain.Email, error) {
	pushedEvents, err := c.eventstore.Push(ctx, c.events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(c.model, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToEmail(c.model), nil
}
