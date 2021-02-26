package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

//ResendInitialMail resend inital mail and changes email if provided
func (c *Commands) ResendInitialMail(ctx context.Context, userID, email, resourceOwner string) (err error) {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.UserIDMissing")
	}

	existingCode, err := c.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9df", "Errors.User.NotFound")
	}
	if existingCode.UserState != domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.AlreadyInitialised")
	}
	var events []eventstore.EventPusher
	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	if email != "" && existingCode.Email != email {
		changedEvent, _ := existingCode.NewChangedEvent(ctx, userAgg, email)
		events = append(events, changedEvent)
	}
	initCode, err := domain.NewInitUserCode(c.initializeUserCode)
	if err != nil {
		return err
	}
	events = append(events, user.NewHumanInitialCodeAddedEvent(ctx, userAgg, initCode.Code, initCode.Expiry))
	_, err = c.eventstore.PushEvents(ctx, events...)
	return err
}

func (c *Commands) HumanVerifyInitCode(ctx context.Context, userID, resourceOwner, code, passwordString string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-mkM9f", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-44G8s", "Errors.User.Code.Empty")
	}

	existingCode, err := c.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-mmn5f", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, c.initializeUserCode)
	if err != nil {
		_, err = c.eventstore.PushEvents(ctx, user.NewHumanInitializedCheckFailedEvent(ctx, userAgg))
		logging.LogWithFields("COMMAND-Dg2z5", "userID", userAgg.ID).OnError(err).Error("NewHumanInitializedCheckFailedEvent push failed")
		return caos_errs.ThrowInvalidArgument(err, "COMMAND-11v6G", "Errors.User.Code.Invalid")
	}
	events := []eventstore.EventPusher{
		user.NewHumanInitializedCheckSucceededEvent(ctx, userAgg),
	}
	if !existingCode.IsEmailVerified {
		events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
	}
	if passwordString != "" {
		passwordWriteModel := NewHumanPasswordWriteModel(userID, existingCode.ResourceOwner)
		password := &domain.Password{
			SecretString:   passwordString,
			ChangeRequired: false,
		}
		passwordEvent, err := c.changePassword(ctx, "", password, userAgg, passwordWriteModel)
		if err != nil {
			return err
		}
		events = append(events, passwordEvent)
	}
	events = append(events, user.NewHumanInitialCodeSentEvent(ctx, userAgg))
	_, err = c.eventstore.PushEvents(ctx, events...)
	return err
}

func (c *Commands) HumanInitCodeSent(ctx context.Context, orgID, userID string) (err error) {
	existingInitCode, err := c.getHumanInitWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingInitCode.UserState == domain.UserStateUnspecified || existingInitCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-556zg", "Errors.User.Code.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingInitCode.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, user.NewHumanInitialCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) getHumanInitWriteModelByID(ctx context.Context, userID, resourceowner string) (*HumanInitCodeWriteModel, error) {
	initWriteModel := NewHumanInitCodeWriteModel(userID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, initWriteModel)
	if err != nil {
		return nil, err
	}
	return initWriteModel, nil
}
