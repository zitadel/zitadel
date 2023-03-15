package command

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

// ResendInitialMail resend initial mail and changes email if provided
func (c *Commands) ResendInitialMail(ctx context.Context, userID string, email domain.EmailAddress, resourceOwner string, initCodeGenerator crypto.Generator) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2n8vs", "Errors.User.UserIDMissing")
	}

	existingCode, err := c.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2M9df", "Errors.User.NotFound")
	}
	if existingCode.UserState != domain.UserStateInitial {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.AlreadyInitialised")
	}
	var events []eventstore.Command
	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	if email != "" && existingCode.Email != email {
		changedEvent, _ := existingCode.NewChangedEvent(ctx, userAgg, email)
		events = append(events, changedEvent)
	}
	initCode, err := domain.NewInitUserCode(initCodeGenerator)
	if err != nil {
		return nil, err
	}
	events = append(events, user.NewHumanInitialCodeAddedEvent(ctx, userAgg, initCode.Code, initCode.Expiry))
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingCode, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingCode.WriteModel), nil
}

func (c *Commands) HumanVerifyInitCode(ctx context.Context, userID, resourceOwner, code, passwordString string, initCodeGenerator crypto.Generator) error {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-mkM9f", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-44G8s", "Errors.User.Code.Empty")
	}

	existingCode, err := c.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-mmn5f", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, initCodeGenerator)
	if err != nil {
		_, err = c.eventstore.Push(ctx, user.NewHumanInitializedCheckFailedEvent(ctx, userAgg))
		logging.WithFields("userID", userAgg.ID).OnError(err).Error("NewHumanInitializedCheckFailedEvent push failed")
		return caos_errs.ThrowInvalidArgument(err, "COMMAND-11v6G", "Errors.User.Code.Invalid")
	}
	events := []eventstore.Command{
		user.NewHumanInitializedCheckSucceededEvent(ctx, userAgg),
	}
	if !existingCode.IsEmailVerified {
		events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
	}
	if passwordString != "" {
		passwordWriteModel := NewHumanPasswordWriteModel(userID, existingCode.ResourceOwner)
		passwordWriteModel.UserState = domain.UserStateActive
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
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) HumanInitCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M9fs", "Errors.IDMissing")
	}
	existingInitCode, err := c.getHumanInitWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingInitCode.UserState == domain.UserStateUnspecified || existingInitCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-556zg", "Errors.User.Code.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingInitCode.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanInitialCodeSentEvent(ctx, userAgg))
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
