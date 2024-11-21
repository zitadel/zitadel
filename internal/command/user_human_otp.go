package command

import (
	"context"
	"time"

	"github.com/pquerna/otp"
	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ImportHumanTOTP(ctx context.Context, userID, userAgentID, resourceOwner string, key string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
	if err := c.checkPermissionUpdateUserCredentials(ctx, human.ResourceOwner, userID); err != nil {
		return nil, err
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
		issuer = http_util.DomainContext(ctx).RequestedDomain()
	}
	key, err := domain.NewTOTPKey(issuer, accountName)
	if err != nil {
		return nil, err
	}
	encryptedSecret, err := crypto.Encrypt([]byte(key.Secret()), c.multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	return &preparedTOTP{
		wm:      otpWriteModel,
		userAgg: userAgg,
		key:     key,
		cmds: []eventstore.Command{
			user.NewHumanOTPAddedEvent(ctx, userAgg, encryptedSecret),
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
	if err := c.checkPermissionUpdateUserCredentials(ctx, existingOTP.ResourceOwner, userID); err != nil {
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
	commands, err := checkTOTP(
		ctx,
		userID,
		resourceOwner,
		code,
		c.eventstore.FilterToQueryReducer,
		c.multifactors.OTP.CryptoMFA,
		authRequestDomainToAuthRequestInfo(authRequest),
	)

	_, pushErr := c.eventstore.Push(ctx, commands...)
	logging.OnError(pushErr).Error("error create password check failed event")
	return err
}

func checkTOTP(
	ctx context.Context,
	userID, resourceOwner, code string,
	queryReducer func(ctx context.Context, r eventstore.QueryReducer) error,
	alg crypto.EncryptionAlgorithm,
	optionalAuthRequestInfo *user.AuthRequestInfo,
) ([]eventstore.Command, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP := NewHumanTOTPWriteModel(userID, resourceOwner)
	err := queryReducer(ctx, existingOTP)
	if err != nil {
		return nil, err
	}
	if existingOTP.State != domain.MFAStateReady {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotReady")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	verifyErr := domain.VerifyTOTP(code, existingOTP.Secret, alg)

	// recheck for additional events (failed OTP checks or locks)
	recheckErr := queryReducer(ctx, existingOTP)
	if recheckErr != nil {
		return nil, recheckErr
	}
	if existingOTP.UserLocked {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SF3fg", "Errors.User.Locked")
	}

	// the OTP check succeeded and the user was not locked in the meantime
	if verifyErr == nil {
		return []eventstore.Command{user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, optionalAuthRequestInfo)}, nil
	}

	// the OTP check failed, therefore check if the limit was reached and the user must additionally be locked
	commands := make([]eventstore.Command, 0, 2)
	commands = append(commands, user.NewHumanOTPCheckFailedEvent(ctx, userAgg, optionalAuthRequestInfo))
	lockoutPolicy, err := getLockoutPolicy(ctx, existingOTP.ResourceOwner, queryReducer)
	if err != nil {
		return nil, err
	}
	if lockoutPolicy.MaxOTPAttempts > 0 && existingOTP.CheckFailedCount+1 >= lockoutPolicy.MaxOTPAttempts {
		commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
	}
	return commands, verifyErr
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
	if err := c.checkPermissionUpdateUser(ctx, existingOTP.ResourceOwner, userID); err != nil {
		return nil, err
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
	otpWriteModel, err := c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateUserCredentials(ctx, otpWriteModel.ResourceOwner(), userID); err != nil {
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
	if err := c.checkPermissionUpdateUser(ctx, existingOTP.WriteModel.ResourceOwner, userID); err != nil {
		return nil, err
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
	codeAddedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo, generatorID string) eventstore.Command {
		return user.NewHumanOTPSMSCodeAddedEvent(ctx, aggregate, code, expiry, info, generatorID)
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
		c.newPhoneCode,
	)
}

func (c *Commands) HumanOTPSMSCodeSent(ctx context.Context, userID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) (err error) {
	smsWriteModel := func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error) {
		return c.otpSMSWriteModelByID(ctx, userID, resourceOwner)
	}
	codeSentEvent := func(ctx context.Context, aggregate *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanOTPSMSCodeSentEvent(ctx, aggregate, generatorInfo)
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
	commands, err := checkOTP(
		ctx,
		userID,
		code,
		resourceOwner,
		authRequest,
		writeModel,
		c.eventstore.FilterToQueryReducer,
		c.userEncryption,
		c.phoneCodeVerifier,
		succeededEvent,
		failedEvent,
	)
	if len(commands) > 0 {
		_, pushErr := c.eventstore.Push(ctx, commands...)
		logging.WithFields("userID", userID).OnError(pushErr).Error("otp failure check push failed")
	}
	return err
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
	if err := c.checkPermissionUpdateUserCredentials(ctx, otpWriteModel.ResourceOwner(), userID); err != nil {
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
	if err := c.checkPermissionUpdateUser(ctx, existingOTP.WriteModel.ResourceOwner, userID); err != nil {
		return nil, err
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
	codeAddedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo, _ string) eventstore.Command {
		return user.NewHumanOTPEmailCodeAddedEvent(ctx, aggregate, code, expiry, info)
	}
	generateCode := func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		code, err := c.newEncryptedCodeWithDefault(ctx, filter, typ, alg, defaultConfig)
		return code, "", err
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
		generateCode,
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
	commands, err := checkOTP(
		ctx,
		userID,
		code,
		resourceOwner,
		authRequest,
		writeModel,
		c.eventstore.FilterToQueryReducer,
		c.userEncryption,
		nil, // email currently always uses local code checks
		succeededEvent,
		failedEvent,
	)
	if len(commands) > 0 {
		_, pushErr := c.eventstore.Push(ctx, commands...)
		logging.WithFields("userID", userID).OnError(pushErr).Error("otp failure check push failed")
	}
	return err
}

// sendHumanOTP creates a code for a registered mechanism (sms / email), which is used for a check (during login)
func (c *Commands) sendHumanOTP(
	ctx context.Context,
	userID, resourceOwner string,
	authRequest *domain.AuthRequest,
	writeModelByID func(ctx context.Context, userID string, resourceOwner string) (OTPWriteModel, error),
	secretGeneratorType domain.SecretGeneratorType,
	defaultSecretGenerator *crypto.GeneratorConfig,
	codeAddedEvent func(ctx context.Context, aggregate *eventstore.Aggregate, code *crypto.CryptoValue, expiry time.Duration, info *user.AuthRequestInfo, generatorID string) eventstore.Command,
	generateCode func(ctx context.Context, filter preparation.FilterToQueryReducer, secretGeneratorType domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error),
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
	code, generatorID, err := generateCode(ctx, c.eventstore.Filter, secretGeneratorType, c.userEncryption, defaultSecretGenerator) //nolint:staticcheck
	if err != nil {
		return err
	}
	userAgg := &user.NewAggregate(userID, resourceOwner).Aggregate
	_, err = c.eventstore.Push(ctx, codeAddedEvent(ctx, userAgg, code.CryptedCode(), code.CodeExpiry(), authRequestDomainToAuthRequestInfo(authRequest), generatorID))
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

func checkOTP(
	ctx context.Context,
	userID, code, resourceOwner string,
	authRequest *domain.AuthRequest,
	writeModelByID func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error),
	queryReducer func(ctx context.Context, r eventstore.QueryReducer) error,
	alg crypto.EncryptionAlgorithm,
	getCodeVerifier func(ctx context.Context, id string) (senders.CodeGenerator, error),
	checkSucceededEvent, checkFailedEvent func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command,
) ([]eventstore.Command, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty")
	}
	existingOTP, err := writeModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingOTP.OTPAdded() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady")
	}
	if existingOTP.Code() == nil && existingOTP.GeneratorID() == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound")
	}
	userAgg := &user.NewAggregate(userID, existingOTP.ResourceOwner()).Aggregate
	verifyErr := verifyCode(
		ctx,
		existingOTP.CodeCreationDate(),
		existingOTP.CodeExpiry(),
		existingOTP.Code(),
		existingOTP.GeneratorID(),
		existingOTP.ProviderVerificationID(),
		code,
		alg,
		getCodeVerifier,
	)
	// recheck for additional events (failed OTP checks or locks)
	recheckErr := queryReducer(ctx, existingOTP)
	if recheckErr != nil {
		return nil, recheckErr
	}
	if existingOTP.UserLocked() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-S6h4R", "Errors.User.Locked")
	}

	// the OTP check succeeded and the user was not locked in the meantime
	if verifyErr == nil {
		return []eventstore.Command{checkSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest))}, nil
	}

	// the OTP check failed, therefore check if the limit was reached and the user must additionally be locked
	commands := make([]eventstore.Command, 0, 2)
	commands = append(commands, checkFailedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	lockoutPolicy, lockoutErr := getLockoutPolicy(ctx, existingOTP.ResourceOwner(), queryReducer)
	logging.OnError(lockoutErr).Error("unable to get lockout policy")
	if lockoutPolicy != nil && lockoutPolicy.MaxOTPAttempts > 0 && existingOTP.CheckFailedCount()+1 >= lockoutPolicy.MaxOTPAttempts {
		commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
	}
	return commands, verifyErr
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
