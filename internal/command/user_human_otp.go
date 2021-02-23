package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddHumanOTP(ctx context.Context, userID, resourceowner string) (*domain.OTP, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}
	human, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		logging.Log("COMMAND-DAqe1").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get human for loginname")
		return nil, err
	}
	org, err := r.getOrg(ctx, human.ResourceOwner)
	if err != nil {
		logging.Log("COMMAND-Cm0ds").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org for loginname")
		return nil, err
	}
	orgPolicy, err := r.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		logging.Log("COMMAND-y5zv9").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org policy for loginname")
		return nil, err
	}
	otpWriteModel, err := r.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.State == domain.MFAStateReady {
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
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanOTPAddedEvent(ctx, userAgg, secret))

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

func (r *CommandSide) HumanCheckMFAOTPSetup(ctx context.Context, userID, code, userAgentID, resourceowner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotExisting")
	}
	if existingOTP.State == domain.MFAStateReady {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady")
	}
	if err := domain.VerifyMFAOTP(code, existingOTP.Secret, r.multifactors.OTP.CryptoMFA); err != nil {
		return err
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)

	_, err = r.eventstore.PushEvents(ctx, user.NewHumanOTPVerifiedEvent(ctx, userAgg, userAgentID))
	return err
}

func (r *CommandSide) HumanCheckMFAOTP(ctx context.Context, userID, code, resourceowner string, authRequest *domain.AuthRequest) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingOTP.State != domain.MFAStateReady {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	err = domain.VerifyMFAOTP(code, existingOTP.Secret, r.multifactors.OTP.CryptoMFA)
	if err == nil {
		_, err = r.eventstore.PushEvents(ctx, user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		return err
	}
	_, pushErr := r.eventstore.PushEvents(ctx, user.NewHumanOTPCheckFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	logging.Log("COMMAND-9fj7s").OnError(pushErr).Error("error create password check failed event")
	return err
}

func (r *CommandSide) HumanRemoveOTP(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-Hd9sd", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, user.NewHumanOTPRemovedEvent(ctx, userAgg))
	return err
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
