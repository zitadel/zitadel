package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddHumanOTP(ctx context.Context, userID, resourceowner string) (*domain.OTP, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}
	human, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	org, err := r.getOrg(ctx, human.ResourceOwner)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := r.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		return nil, err
	}
	otpWriteModel, err := r.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.MFAState == domain.MFAStateReady {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	accountName := domain.GenerateLoginName(human.GetUsername(), org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = human.EmailAddress
	}
	key, secret, err := domain.NewOTPKey(r.multifactors.OTP.Issuer, accountName, r.multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	userAgg.PushEvents(
		user.NewHumanOTPAddedEvent(ctx, secret),
	)

	err = r.eventstore.PushAggregate(ctx, otpWriteModel, userAgg)
	if err != nil {
		return nil, err
	}
	return &domain.OTP{
		ObjectRoot: models.ObjectRoot{
			AggregateID: human.AggregateID,
		},
		SecretString: key.Secret(),
		Url:          key.URL(),
	}, nil
}

func (r *CommandSide) CheckMFAOTPSetup(ctx context.Context, userID, code, userAgentID, resourceowner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingOTP.OTPState == domain.OTPStateUnspecified || existingOTP.OTPState == domain.OTPStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.MFA.OTP.NotExisting")
	}
	if existingOTP.MFAState == domain.MFAStateReady {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady")
	}
	if err := domain.VerifyMFAOTP(code, existingOTP.Secret, r.multifactors.OTP.CryptoMFA); err != nil {
		return err
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	userAgg.PushEvents(
		user.NewHumanOTPVerifiedEvent(ctx, userAgentID),
	)

	return r.eventstore.PushAggregate(ctx, existingOTP, userAgg)
}

func (r *CommandSide) RemoveHumanOTP(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingOTP.OTPState == domain.OTPStateUnspecified || existingOTP.OTPState == domain.OTPStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	userAgg.PushEvents(
		user.NewHumanOTPRemovedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, existingOTP, userAgg)
}

func (r *CommandSide) otpWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
