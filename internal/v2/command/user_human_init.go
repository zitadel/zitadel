package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

//ResendInitialMail resend inital mail and changes email if provided
func (r *CommandSide) ResendInitialMail(ctx context.Context, userID, email, resourceOwner string) (err error) {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2n8vs", "Errors.User.UserIDMissing")
	}

	existingCode, err := r.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9df", "Errors.User.NotFound")
	}
	if existingCode.UserState != domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.AlreadyInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	if email != "" && existingCode.Email != email {
		changedEvent, _ := existingCode.NewChangedEvent(ctx, email)
		userAgg.PushEvents(changedEvent)
	}
	initCode, err := domain.NewInitUserCode(r.initializeUserCode)
	if err != nil {
		return err
	}
	userAgg.PushEvents(user.NewHumanInitialCodeAddedEvent(ctx, initCode.Code, initCode.Expiry))
	return r.eventstore.PushAggregate(ctx, existingCode, userAgg)
}

func (r *CommandSide) HumanVerifyInitCode(ctx context.Context, userID, resourceOwner, code, passwordString string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-mkM9f", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-44G8s", "Errors.User.Code.Empty")
	}

	existingCode, err := r.getHumanInitWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-mmn5f", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, r.initializeUserCode)
	if err != nil {
		userAgg.PushEvents(user.NewHumanInitializedCheckFailedEvent(ctx))
		err = r.eventstore.PushAggregate(ctx, existingCode, userAgg)
		logging.LogWithFields("COMMAND-Dg2z5", "userID", userAgg.ID()).OnError(err).Error("NewHumanInitializedCheckFailedEvent push failed")
		return caos_errs.ThrowInvalidArgument(err, "COMMAND-11v6G", "Errors.User.Code.Invalid")
	}

	userAgg.PushEvents(user.NewHumanInitializedCheckSucceededEvent(ctx))
	if !existingCode.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	}
	if passwordString != "" {
		passwordWriteModel := NewHumanPasswordWriteModel(userID, existingCode.ResourceOwner)
		password := &domain.Password{
			SecretString:   passwordString,
			ChangeRequired: true,
		}
		err = r.changePassword(ctx, existingCode.ResourceOwner, userID, "", password, userAgg, passwordWriteModel)
		if err != nil {
			return err
		}
	}
	return r.eventstore.PushAggregate(ctx, existingCode, userAgg)
}

func (r *CommandSide) HumanInitCodeSent(ctx context.Context, orgID, userID string) (err error) {
	existingInitCode, err := r.getHumanInitWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingInitCode.UserState == domain.UserStateUnspecified || existingInitCode.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-556zg", "Errors.User.Code.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingInitCode.WriteModel)
	userAgg.PushEvents(user.NewHumanInitialCodeSentEvent(ctx))
	return r.eventstore.PushAggregate(ctx, existingInitCode, userAgg)
}

func (r *CommandSide) getHumanInitWriteModelByID(ctx context.Context, userID, resourceowner string) (*HumanInitCodeWriteModel, error) {
	initWriteModel := NewHumanInitCodeWriteModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, initWriteModel)
	if err != nil {
		return nil, err
	}
	return initWriteModel, nil
}
