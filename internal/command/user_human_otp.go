package command

import (
	"context"
	"time"

	"github.com/pquerna/otp"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ImportHumanTOTP(ctx context.Context, userID, userAgentID, resourceOwner string, key string) error {
	encryptedSecret, err := crypto.Encrypt([]byte(key), c.multifactors.OTP.CryptoMFA)
	if err != nil {
		return err
	}
	if err = c.checkUserExists(ctx, userID, resourceOwner); err != nil {
		return err
	}

	otpWriteModel, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if otpWriteModel.State == domain.MFAStateReady {
		return zerrors.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)

	_, err = c.eventstore.Push(ctx,
		user.NewHumanOTPAddedEvent(ctx, userAgg, encryptedSecret),
		user.NewHumanOTPVerifiedEvent(ctx, userAgg, userAgentID),
	)
	return err
}

func (c *Commands) AddHumanTOTP(ctx context.Context, userID, resourceOwner string) (*domain.TOTP, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}
	prep, err := c.createHumanTOTP(ctx, userID, resourceOwner)
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
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-SqyJz", "Errors.User.NotFound")
	}
	org, err := c.getOrg(ctx, human.ResourceOwner)
	if err != nil {
		logging.WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org for loginname")
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-55M9f", "Errors.Org.NotFound")
	}
	orgPolicy, err := c.domainPolicyWriteModel(ctx, org.AggregateID)
	if err != nil {
		logging.WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to get org policy for loginname")
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-8ugTs", "Errors.Org.DomainPolicy.NotFound")
	}

	otpWriteModel, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.State == domain.MFAStateReady {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
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

func (c *Commands) HumanCheckMFATOTPSetup(ctx context.Context, userID, code, userAgentID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotExisting")
	}
	if existingOTP.State == domain.MFAStateReady {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady")
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

func (c *Commands) HumanCheckMFATOTP(ctx context.Context, userID, code, resourceOwner string, authRequest *domain.AuthRequest) error {
	if userID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if existingOTP.State != domain.MFAStateReady {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotReady")
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.totpWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOTP.State == domain.MFAStateUnspecified || existingOTP.State == domain.MFAStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Hd9sd", "Errors.User.MFA.OTP.NotExisting")
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

// AddHumanOTPSMS adds the OTP SMS factor to a user.
// It can only be added if it not already is and the phone has to be verified.
func (c *Commands) AddHumanOTPSMS(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	return c.addHumanOTPSMS(ctx, userID, resourceOwner)
}

// AddHumanOTPSMSWithCheckSucceeded adds the OTP SMS factor to a user.
// It can only be added if it's not already and the phone has to be verified.
// An OTPSMSCheckSucceededEvent will be added to the passed AuthRequest, if not nil.
func (c *Commands) AddHumanOTPSMSWithCheckSucceeded(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) (*domain.ObjectDetails, error) {
	if authRequest == nil {
		return c.addHumanOTPSMS(ctx, userID, resourceOwner)
	}
	event := func(ctx context.Context, userAgg *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanOTPSMSCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest))
	}
	return c.addHumanOTPSMS(ctx, userID, resourceOwner, event)
}

func (c *Commands) addHumanOTPSMS(ctx context.Context, userID, resourceOwner string, events ...eventCallback) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-QSF2s", "Errors.User.UserIDMissing")
	}
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	otpWriteModel, err := c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.otpAdded {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-Ad3g2", "Errors.User.MFA.OTP.AlreadyReady")
	}
	if !otpWriteModel.phoneVerified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Q54j2", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	cmds := make([]eventstore.Command, len(events)+1)
	cmds[0] = user.NewHumanOTPSMSAddedEvent(ctx, userAgg)
	for i, event := range events {
		cmds[i+1] = event(ctx, userAgg)
	}
	if err = c.pushAppendAndReduce(ctx, otpWriteModel, cmds...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&otpWriteModel.WriteModel), nil
}

func (c *Commands) RemoveHumanOTPSMS(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-S3br2", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if userID != authz.GetCtxData(ctx).UserID {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, existingOTP.WriteModel.ResourceOwner, userID); err != nil {
			return nil, err
		}
	}
	if !existingOTP.otpAdded {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Sr3h3", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	if err = c.pushAppendAndReduce(ctx, existingOTP, user.NewHumanOTPSMSRemovedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) HumanSendOTPSMS(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) error {
	smsWriteModel := func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error) {
		return c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	}
	codeAddedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPSMSCodeAddedEvent(ctx, aggregate, code, expiry, info)
	}
	return c.sendHumanOTP(
		ctx,
		userID,
		resourceOwner,
		authRequest,
		smsWriteModel,
		domain.SecretGeneratorTypeOTPSMS,
		c.defaultSecretGenerators.OTPSMS,
		codeAddedEvent,
	)
}

func (c *Commands) HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string) (err error) {
	smsWriteModel := func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error) {
		return c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	}
	codeSentEvent := func(ctx context.Context, aggregate *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanOTPSMSCodeSentEvent(ctx, aggregate)
	}
	return c.humanOTPSent(ctx, userID, resourceOwner, smsWriteModel, codeSentEvent)
}

func (c *Commands) HumanCheckOTPSMS(ctx context.Context, userID, code, resourceOwner string, authRequest *domain.AuthRequest) error {
	writeModel := func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error) {
		return c.otpSMSCodeWriteModelByID(ctx, userID, resourceOwner)
	}
	succeededEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPSMSCheckSucceededEvent(ctx, aggregate, authRequestDomainToAuthRequestInfo(authRequest))
	}
	failedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPSMSCheckFailedEvent(ctx, aggregate, authRequestDomainToAuthRequestInfo(authRequest))
	}
	return c.humanCheckOTP(
		ctx,
		userID,
		code,
		resourceOwner,
		authRequest,
		writeModel,
		succeededEvent,
		failedEvent,
	)
}

// AddHumanOTPEmail adds the OTP Email factor to a user.
// It can only be added if it not already is and the phone has to be verified.
func (c *Commands) AddHumanOTPEmail(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	return c.addHumanOTPEmail(ctx, userID, resourceOwner)
}

// AddHumanOTPEmailWithCheckSucceeded adds the OTP Email factor to a user.
// It can only be added if it's not already and the email has to be verified.
// An OTPEmailCheckSucceededEvent will be added to the passed AuthRequest, if not nil.
func (c *Commands) AddHumanOTPEmailWithCheckSucceeded(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) (*domain.ObjectDetails, error) {
	if authRequest == nil {
		return c.addHumanOTPEmail(ctx, userID, resourceOwner)
	}
	event := func(ctx context.Context, userAgg *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanOTPEmailCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest))
	}
	return c.addHumanOTPEmail(ctx, userID, resourceOwner, event)
}

func (c *Commands) addHumanOTPEmail(ctx context.Context, userID, resourceOwner string, events ...eventCallback) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Sg1hz", "Errors.User.UserIDMissing")
	}
	otpWriteModel, err := c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.otpAdded {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-MKL2s", "Errors.User.MFA.OTP.AlreadyReady")
	}
	if !otpWriteModel.emailVerified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-KLJ2d", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	cmds := make([]eventstore.Command, len(events)+1)
	cmds[0] = user.NewHumanOTPEmailAddedEvent(ctx, userAgg)
	for i, event := range events {
		cmds[i+1] = event(ctx, userAgg)
	}
	if err = c.pushAppendAndReduce(ctx, otpWriteModel, cmds...); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&otpWriteModel.WriteModel), nil
}

func (c *Commands) RemoveHumanOTPEmail(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-S2h11", "Errors.User.UserIDMissing")
	}

	existingOTP, err := c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if userID != authz.GetCtxData(ctx).UserID {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, existingOTP.WriteModel.ResourceOwner, userID); err != nil {
			return nil, err
		}
	}
	if !existingOTP.otpAdded {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-b312D", "Errors.User.MFA.OTP.NotExisting")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	if err = c.pushAppendAndReduce(ctx, existingOTP, user.NewHumanOTPEmailRemovedEvent(ctx, userAgg)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingOTP.WriteModel), nil
}

func (c *Commands) HumanSendOTPEmail(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) error {
	smsWriteModel := func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error) {
		return c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	}
	codeAddedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPEmailCodeAddedEvent(ctx, aggregate, code, expiry, info)
	}
	return c.sendHumanOTP(
		ctx,
		userID,
		resourceOwner,
		authRequest,
		smsWriteModel,
		domain.SecretGeneratorTypeOTPEmail,
		c.defaultSecretGenerators.OTPEmail,
		codeAddedEvent,
	)
}

func (c *Commands) HumanOTPEmailCodeSent(ctx context.Context, userID, resourceOwner string) (err error) {
	smsWriteModel := func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error) {
		return c.otpEmailWriteModelByID(ctx, userID, resourceOwner)
	}
	codeSentEvent := func(ctx context.Context, aggregate *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanOTPEmailCodeSentEvent(ctx, aggregate)
	}
	return c.humanOTPSent(ctx, userID, resourceOwner, smsWriteModel, codeSentEvent)
}

func (c *Commands) HumanCheckOTPEmail(ctx context.Context, userID, code, resourceOwner string, authRequest *domain.AuthRequest) error {
	writeModel := func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error) {
		return c.otpEmailCodeWriteModelByID(ctx, userID, resourceOwner)
	}
	succeededEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPEmailCheckSucceededEvent(ctx, aggregate, authRequestDomainToAuthRequestInfo(authRequest))
	}
	failedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
		return user.NewHumanOTPEmailCheckFailedEvent(ctx, aggregate, authRequestDomainToAuthRequestInfo(authRequest))
	}
	return c.humanCheckOTP(
		ctx,
		userID,
		code,
		resourceOwner,
		authRequest,
		writeModel,
		succeededEvent,
		failedEvent,
	)
}

// sendHumanOTP creates a code for a registered mechanism (sms / email), which is used for a check (during login)
func (c *Commands) sendHumanOTP(
	ctx context.Context,
	userID, resourceOwner string,
	authRequest *domain.AuthRequest,
	writeModelByID func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error),
	secretGeneratorType domain.SecretGeneratorType,
	defaultSecretGenerator *crypto.GeneratorConfig,
	codeAddedEvent func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo) eventstore.Command,
) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-S3SF1", "Errors.User.UserIDMissing")
	}
	existingOTP, err := writeModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !existingOTP.OTPAdded() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFD52", "Errors.User.MFA.OTP.NotReady")
	}
	config, err := secretGeneratorConfigWithDefault(ctx, c.eventstore.Filter, secretGeneratorType, defaultSecretGenerator)
	if err != nil {
		return err
	}
	gen := crypto.NewEncryptionGenerator(*config, c.userEncryption)
	value, _, err := crypto.NewCode(gen)
	if err != nil {
		return err
	}
	userAgg := &user.NewAggregate(userID, resourceOwner).Aggregate
	_, err = c.eventstore.Push(ctx, codeAddedEvent(ctx, userAgg, value, gen.Expiry(), authRequestDomainToAuthRequestInfo(authRequest)))
	return err
}

func (c *Commands) humanOTPSent(
	ctx context.Context,
	userID, resourceOwner string,
	writeModelByID func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error),
	codeSentEvent func(ctx context.Context, aggregate *eventstore.Aggregate) eventstore.Command,
) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-AE2h2", "Errors.User.UserIDMissing")
	}
	existingOTP, err := writeModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !existingOTP.OTPAdded() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SD3gh", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := &user.NewAggregate(userID, resourceOwner).Aggregate
	_, err = c.eventstore.Push(ctx, codeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) humanCheckOTP(
	ctx context.Context,
	userID, code, resourceOwner string,
	authRequest *domain.AuthRequest,
	writeModelByID func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error),
	checkSucceededEvent func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command,
	checkFailedEvent func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command,
) error {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty")
	}
	existingOTP, err := writeModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !existingOTP.OTPAdded() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady")
	}
	if existingOTP.Code() == nil {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound")
	}
	userAgg := &user.NewAggregate(userID, existingOTP.ResourceOwner()).Aggregate
	err = crypto.VerifyCodeWithAlgorithm(existingOTP.CodeCreationDate(), existingOTP.CodeExpiry(), existingOTP.Code(), code, c.userEncryption)
	if err == nil {
		_, err = c.eventstore.Push(ctx, checkSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
		return err
	}
	_, pushErr := c.eventstore.Push(ctx, checkFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	logging.WithFields("userID", userID).OnError(pushErr).Error("otp failure check push failed")
	return err
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

func (c *Commands) otpSMSCodeWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPSMSCodeWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPSMSCodeWriteModel(userID, resourceOwner)
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

func (c *Commands) otpEmailCodeWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPEmailCodeWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPEmailCodeWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
