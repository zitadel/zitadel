package command

import (
	"context"

	"github.com/pquerna/otp"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ImportHumanTOTP(ctx context.Context, userID, userAgentID, resourceowner string, key string) error {
	encryptedSecret, err := crypto.Encrypt([]byte(key), c.multifactors.OTP.CryptoMFA)
	if err != nil {
		return err
	}
	if err = c.checkUserExists(ctx, userID, resourceowner); err != nil {
		return err
	}

	otpWriteModel, err := c.totpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if otpWriteModel.State == domain.MFAStateReady {
		return caos_errs.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)

	_, err = c.eventstore.Push(ctx,
		user.NewHumanOTPAddedEvent(ctx, userAgg, encryptedSecret),
		user.NewHumanOTPVerifiedEvent(ctx, userAgg, userAgentID),
	)
	return err
}

func (c *Commands) AddHumanTOTP(ctx context.Context, userID, resourceowner string) (*domain.TOTP, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}
	prep, err := c.createHumanTOTP(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, prep.wm, prep.cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.TOTP{
		ObjectDetails: writeModelToObjectDetails(&prep.wm.WriteModel),
		Secret:        prep.key.Secret(),
		URI:           prep.key.URL(),
	}, nil
}

type preparedTOTP struct {
	wm      *HumanTOTPWriteModel
	userAgg *eventstore.Aggregate
	key     *otp.Key
	cmds    []eventstore.Command
}

func (c *Commands) createHumanTOTP(ctx context.Context, userID, resourceOwner string) (*preparedTOTP, error) {
	human, err := c.getHuman(ctx, userID, resourceOwner)
	if err != nil {
		logging.WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get human for loginname")
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-MM9fs", "Errors.User.NotFound")
	}
	org, err := c.getOrg(ctx, human.ResourceOwner)
	if err != nil {
		logging.WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org for loginname")
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-55M9f", "Errors.Org.NotFound")
	}
	orgPolicy, err := c.getOrgDomainPolicy(ctx, org.AggregateID)
	if err != nil {
		logging.WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org policy for loginname")
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-8ugTs", "Errors.Org.DomainPolicy.NotFound")
	}

	otpWriteModel, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.State == domain.MFAStateReady {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)

	accountName := domain.GenerateLoginName(human.GetUsername(), org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = string(human.EmailAddress)
	}
	issuer := c.multifactors.OTP.Issuer
	if issuer == "" {
		issuer = authz.GetInstance(ctx).RequestedDomain()
	}
	key, secret, err := domain.NewTOTPKey(issuer, accountName, c.multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	return &preparedTOTP{
		wm:      otpWriteModel,
		userAgg: userAgg,
		key:     key,
		cmds: []eventstore.Command{
			user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
		},
	}, nil
}

func (c *Commands) HumanCheckMFATOTPSetup(ctx context.Context, userID, code, userAgentID, resourceowner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotExisting")
	}
	if existingOTP.State == domain.MFAStateReady {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady")
	}
	if err := domain.VerifyTOTP(code, existingOTP.Secret, c.multifactors.OTP.CryptoMFA); err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanOTPVerifiedEvent(ctx, userAgg, userAgentID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingOTP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) HumanCheckMFATOTP(ctx context.Context, userID, code, resourceowner string, authRequest *domain.AuthRequest) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingOTP.State != domain.MFAStateReady {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	err = domain.VerifyTOTP(code, existingOTP.Secret, c.multifactors.OTP.CryptoMFA)
	if err == nil {
		_, err = c.eventstore.Push(ctx, user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		return err
	}
	_, pushErr := c.eventstore.Push(ctx, user.NewHumanOTPCheckFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	logging.OnError(pushErr).Error("error create password check failed event")
	return err
}

func (c *Commands) HumanRemoveTOTP(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Hd9sd", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanOTPRemovedEvent(ctx, userAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingOTP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) AddHumanOTPSMS(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-QSF2s", "Errors.User.UserIDMissing")
	}
	otpWriteModel, err := c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.otpAdded {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-Ad3g2", "Errors.User.MFA.OTP.AlreadyReady")
	}
	if !otpWriteModel.phoneVerified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Q54j2", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	if err = c.pushAppendAndReduce(ctx, otpWriteModel, user.NewHumanOTPSMSAddedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&otpWriteModel.WriteModel), nil
}

func (c *Commands) RemoveHumanOTPSMS(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-S3br2", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingOTP.otpAdded {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Sr3h3", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	if err = c.pushAppendAndReduce(ctx, existingOTP, user.NewHumanOTPSMSRemovedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) AddHumanOTPEmail(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Sg1hz", "Errors.User.UserIDMissing")
	}
	otpWriteModel, err := c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.otpAdded {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-MKL2s", "Errors.User.MFA.OTP.AlreadyReady")
	}
	if !otpWriteModel.emailVerified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-KLJ2d", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	if err = c.pushAppendAndReduce(ctx, otpWriteModel, user.NewHumanOTPEmailAddedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&otpWriteModel.WriteModel), nil
}

func (c *Commands) RemoveHumanOTPEmail(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-S2h11", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingOTP.otpAdded {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-b312D", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	if err = c.pushAppendAndReduce(ctx, existingOTP, user.NewHumanOTPEmailRemovedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) totpWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanTOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanTOTPWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) otpSMSWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPSMSWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPSMSWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) otpEmailWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPEmailWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPEmailWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
