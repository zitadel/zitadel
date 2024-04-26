package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeHumanEmail(ctx context.Context, email *domain.Email, emailCodeGenerator crypto.Generator) (*domain.Email, error) {
	if email.AggregateID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing")
	}
	if err := email.Validate(); err != nil {
		return nil, err
	}

	existingEmail, err := c.emailWriteModel(ctx, email.AggregateID, email.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	if existingEmail.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-J8dsk", "Errors.User.NotInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	changedEvent, hasChanged := existingEmail.NewChangedEvent(ctx, userAgg, email.EmailAddress)

	// only continue if there were changes or there were no changes and the email should be set to verified
	if !hasChanged && !(email.IsEmailVerified && existingEmail.IsEmailVerified != email.IsEmailVerified) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2b7fM", "Errors.User.Email.NotChanged")
	}

	events := make([]eventstore.Command, 0)
	if hasChanged {
		events = append(events, changedEvent)
	}
	if email.IsEmailVerified {
		events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
	} else {
		emailCode, _, err := domain.NewEmailCode(emailCodeGenerator)
		if err != nil {
			return nil, err
		}
		events = append(events, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry, ""))
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingEmail, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToEmail(existingEmail), nil
}

func (c *Commands) VerifyHumanEmail(ctx context.Context, userID, code, resourceowner string, emailCodeGenerator crypto.Generator) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-çm0ds", "Errors.User.Code.Empty")
	}

	existingCode, err := c.emailWriteModel(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if existingCode.Code == nil || existingCode.UserState == domain.UserStateUnspecified || existingCode.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3n8ud", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, emailCodeGenerator.Alg())
	if err == nil {
		pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingCode, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&existingCode.WriteModel), nil
	}

	_, err = c.eventstore.Push(ctx, user.NewHumanEmailVerificationFailedEvent(ctx, userAgg))
	logging.LogWithFields("COMMAND-Dg2z5", "userID", userAgg.ID).OnError(err).Error("NewHumanEmailVerificationFailedEvent push failed")
	return nil, zerrors.ThrowInvalidArgument(err, "COMMAND-Gdsgs", "Errors.User.Code.Invalid")
}

func (c *Commands) CreateHumanEmailVerificationCode(ctx context.Context, userID, resourceOwner string, emailCodeGenerator crypto.Generator, authRequestID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingEmail, err := c.emailWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-0Pe4r", "Errors.User.Email.NotFound")
	}
	if existingEmail.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-E3fbw", "Errors.User.NotInitialised")
	}
	if existingEmail.IsEmailVerified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3M9ds", "Errors.User.Email.AlreadyVerified")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	emailCode, _, err := domain.NewEmailCode(emailCodeGenerator)
	if err != nil {
		return nil, err
	}
	if authRequestID == "" {
		authRequestID = existingEmail.AuthRequestID
	}
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry, authRequestID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingEmail, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingEmail.WriteModel), nil
}

func (c *Commands) HumanEmailVerificationCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9fs", "Errors.IDMissing")
	}
	existingEmail, err := c.emailWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return zerrors.ThrowNotFound(nil, "COMMAND-6n8uH", "Errors.User.Email.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanEmailCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) emailWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanEmailWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanEmailWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
