package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type otpCodeVerifyFn func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue, verificationCode string, algorithm crypto.EncryptionAlgorithm) error
type phoneCodeVerifyFn func(ctx context.Context, id string) (senders.CodeGenerator, error)
type changesFn func(ctx context.Context, db database.QueryExecutor, lockoutSettingsRepo LockoutSettingsRepository, humanRepo HumanUserRepository, isVerifyErr bool, checkTime time.Time) database.Changes

type OTPRequestType uint

const (
	OTPSMSRequestType = iota
	OTPEmailRequestType
)

type OTPCheckCommand struct {
	CheckOTP *session_grpc.CheckOTP

	sessionID  string
	instanceID string

	tarpitFunc          tarpitFn
	otpCodeVerifyFunc   otpCodeVerifyFn
	phonecodeVerifyFunc phoneCodeVerifyFn

	requestType         OTPRequestType
	encryptionAlgorithm crypto.EncryptionAlgorithm

	smsChallenge   *SessionChallengeOTPSMS
	emailChallenge *SessionChallengeOTPEmail

	// For Events()
	IsSMSCheckSucceeded   bool
	IsEmailCheckSucceeded bool
	IsUserLocked          bool
	FetchedUser           *User

	CheckTime time.Time
}

func NewOTPCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, otpVerifyFunc otpCodeVerifyFn, phoneVerifyFunc phoneCodeVerifyFn, encryptionAlgo crypto.EncryptionAlgorithm, request *session_grpc.CheckOTP, kind OTPRequestType) *OTPCheckCommand {
	tf := sysConfig.Tarpit.Tarpit()
	if tarpitFunc != nil {
		tf = tarpitFunc
	}

	ea := otpEncryptionAlgo
	if encryptionAlgo != nil {
		ea = encryptionAlgo
	}

	otpVerify := crypto.VerifyCode
	if otpVerifyFunc != nil {
		otpVerify = otpVerifyFunc
	}

	phoneVerify := defaultPhoneVerifier
	if phoneVerifyFunc != nil {
		phoneVerify = phoneVerifyFunc
	}

	return &OTPCheckCommand{
		CheckOTP:            request,
		sessionID:           sessionID,
		instanceID:          instanceID,
		requestType:         kind,
		tarpitFunc:          tf,
		otpCodeVerifyFunc:   otpVerify,
		encryptionAlgorithm: ea,
		phonecodeVerifyFunc: phoneVerify,
	}
}

// RequiresTransaction implements [Transactional].
func (o *OTPCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (o *OTPCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if o.CheckOTP == nil {
		return nil, nil
	}

	toReturn := make([]eventstore.Command, 2, 3)
	userAgg := user.NewAggregate(o.FetchedUser.ID, o.FetchedUser.OrganizationID).Aggregate
	sessionAgg := session.NewAggregate(o.sessionID, o.instanceID).Aggregate

	switch o.requestType {
	case OTPSMSRequestType:
		if o.IsSMSCheckSucceeded {
			toReturn[0] = user.NewHumanOTPSMSCheckSucceededEvent(ctx, &userAgg, nil)
			toReturn[1] = session.NewOTPSMSCheckedEvent(ctx, &sessionAgg, o.CheckTime)
			return toReturn, nil
		}
		toReturn[0] = user.NewHumanOTPSMSCheckFailedEvent(ctx, &userAgg, nil)
		if o.IsUserLocked {
			toReturn[1] = user.NewUserLockedEvent(ctx, &userAgg)
			return append(toReturn, session.NewOTPSMSCheckedEvent(ctx, &sessionAgg, o.CheckTime)), nil
		}
		toReturn[1] = session.NewOTPSMSCheckedEvent(ctx, &sessionAgg, o.CheckTime)
	case OTPEmailRequestType:
		if o.IsEmailCheckSucceeded {
			toReturn[0] = user.NewHumanOTPEmailCheckSucceededEvent(ctx, &userAgg, nil)
			toReturn[1] = session.NewOTPEmailCheckedEvent(ctx, &sessionAgg, o.CheckTime)
			return toReturn, nil
		}
		toReturn[0] = user.NewHumanOTPEmailCheckFailedEvent(ctx, &userAgg, nil)
		if o.IsUserLocked {
			toReturn[1] = user.NewUserLockedEvent(ctx, &userAgg)
			return append(toReturn, session.NewOTPEmailCheckedEvent(ctx, &sessionAgg, o.CheckTime)), nil
		}
		toReturn[1] = session.NewOTPEmailCheckedEvent(ctx, &sessionAgg, o.CheckTime)
	}

	return toReturn, nil
}

// Execute implements [Commander].
func (o *OTPCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.CheckOTP == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	var verificationError error
	var sessionFactorToUpdate func(lastVerified, lastFailed time.Time) database.Change
	var changesFn changesFn
	var lastVerified, lastFailed time.Time

	switch o.requestType {
	case OTPSMSRequestType:
		verificationError = o.executeCheck(
			ctx,
			o.smsChallenge.GeneratorID,
			"Missing Verification ID", // TODO(IAM-Marco): o.smsChallenge should have a VerificationID
			o.smsChallenge.Expiry,
			o.smsChallenge.Code,
		)
		o.IsSMSCheckSucceeded = verificationError == nil

		sessionFactorToUpdate = func(lastVerified, lastFailed time.Time) database.Change {
			return sessionRepo.SetFactor(&SessionFactorOTPSMS{
				LastVerifiedAt: lastVerified,
				LastFailedAt:   lastFailed,
			})
		}

		changesFn = o.getSMSOTPChanges

	case OTPEmailRequestType:
		verificationError = o.executeCheck(
			ctx,
			"",
			"Missing Verification ID", // TODO(IAM-Marco): o.emailChallenge should have a VerificationID
			o.emailChallenge.Expiry,
			o.emailChallenge.Code,
		)
		o.IsEmailCheckSucceeded = verificationError == nil
		sessionFactorToUpdate = func(lastVerified, lastFailed time.Time) database.Change {
			return sessionRepo.SetFactor(&SessionFactorOTPEmail{
				LastVerifiedAt: lastVerified,
				LastFailedAt:   lastFailed,
			})
		}

		changesFn = o.getEmailOTPChanges

	default:
		return zerrors.ThrowInvalidArgument(nil, "DOM-LFE7va", "invalid OTP request type")
	}
	o.CheckTime = time.Now()

	if verificationError == nil {
		lastVerified = o.CheckTime
	} else {
		lastFailed = o.CheckTime
	}

	sessionChange := sessionFactorToUpdate(lastVerified, lastFailed)
	rowCount, err := sessionRepo.Update(ctx, opts.DB(), sessionRepo.IDCondition(o.sessionID),
		sessionChange,
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-h8SL2F", "session"); err != nil {
		return err
	}

	userChanges := changesFn(ctx, opts.DB(), opts.lockoutSettingRepo, userRepo.Human(), verificationError != nil, o.CheckTime)
	rowCount, err = userRepo.Update(ctx, opts.DB(), userRepo.IDCondition(o.FetchedUser.ID), userChanges)
	if err := handleUpdateError(err, 1, rowCount, "DOM-k04c4g", "user"); err != nil {
		return err
	}

	return verificationError
}

func (o *OTPCheckCommand) executeCheck(ctx context.Context,
	generatorID, verificationID string,
	challengeExpiry time.Duration, challengeCode *crypto.CryptoValue) (err error) {
	if generatorID == "" {
		if challengeCode == nil {
			return zerrors.ThrowPreconditionFailed(nil, "DOM-LOWT6u", "Errors.User.Code.NotFound")
		}
		return o.otpCodeVerifyFunc(
			time.Now(), // o.smsChallenge is missing creation date (neither email challenge has)
			challengeExpiry,
			challengeCode,
			o.CheckOTP.GetCode(),
			o.encryptionAlgorithm,
		)
	}

	if o.phonecodeVerifyFunc == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-Cq1y9C", "Errors.User.Code.NotConfigured")
	}
	verifier, err := o.phonecodeVerifyFunc(ctx, generatorID)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-JB24yW", "failed fetching phone verifier")
	}

	return verifier.VerifyCode(verificationID, o.CheckOTP.GetCode())
}

func (o *OTPCheckCommand) getSMSOTPChanges(ctx context.Context, db database.QueryExecutor, lockoutSettingsRepo LockoutSettingsRepository, humanRepo HumanUserRepository, isVerifyErr bool, checkTime time.Time) database.Changes {
	if !isVerifyErr {
		return database.Changes{
			humanRepo.CheckSMSOTP(&CheckTypeSucceeded{SucceededAt: checkTime}),
		}
	}

	toReturn := make(database.Changes, 1, 2)
	toReturn[0] = humanRepo.CheckSMSOTP(&CheckTypeFailed{FailedAt: checkTime})

	lockoutSetting, err := getLockoutPolicy(ctx, db, lockoutSettingsRepo, o.instanceID, o.FetchedUser.OrganizationID)
	if err != nil {
		logger.Error(err.Error())
	}

	if shouldLockUser(lockoutSetting, uint64(o.FetchedUser.Human.Phone.OTP.Check.FailedAttempts)) {
		toReturn = append(toReturn, humanRepo.SetState(UserStateLocked))
		o.IsUserLocked = true
	}

	o.tarpitFunc(uint64(o.FetchedUser.Human.Phone.OTP.Check.FailedAttempts + 1))

	return toReturn
}

func (o *OTPCheckCommand) getEmailOTPChanges(ctx context.Context, db database.QueryExecutor, lockoutSettingsRepo LockoutSettingsRepository, humanRepo HumanUserRepository, isVerifyErr bool, checkTime time.Time) database.Changes {
	if !isVerifyErr {
		return database.Changes{
			humanRepo.CheckEmailOTP(&CheckTypeSucceeded{SucceededAt: checkTime}),
		}
	}

	toReturn := make(database.Changes, 1, 2)
	toReturn[0] = humanRepo.CheckEmailOTP(&CheckTypeFailed{FailedAt: checkTime})

	lockoutSetting, err := getLockoutPolicy(ctx, db, lockoutSettingsRepo, o.instanceID, o.FetchedUser.OrganizationID)
	if err != nil {
		logger.Error(err.Error())
	}

	if shouldLockUser(lockoutSetting, uint64(o.FetchedUser.Human.Email.OTP.Check.FailedAttempts)) {
		toReturn = append(toReturn, humanRepo.SetState(UserStateLocked))
		o.IsUserLocked = true
	}
	o.tarpitFunc(uint64(o.FetchedUser.Human.Email.OTP.Check.FailedAttempts + 1))

	return toReturn
}

// String implements [Commander].
func (o *OTPCheckCommand) String() string {
	return "OTPCheckCommand"
}

// Validate implements [Commander].
func (o *OTPCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.CheckOTP == nil {
		return nil
	}

	if o.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-J2gf7h", "Errors.Missing.SessionID")
	}
	if o.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-RpM3IG", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(o.sessionID)))
	if err := handleGetError(err, "DOM-eppPwQ", "session"); err != nil {
		return err
	}

	user, err := userRepo.Get(ctx, opts.DB(),
		database.WithCondition(userRepo.IDCondition(session.UserID)),
	)
	if err := handleGetError(err, "DOM-TxDSma", "user"); err != nil {
		return err
	}

	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-pBmqRN", "user not human")
	}

	if o.CheckOTP.GetCode() == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-u7KQi4", "Errors.User.Code.Empty")
	}

	o.FetchedUser = user

	if o.requestType == OTPSMSRequestType {
		o.smsChallenge = session.Challenges.GetOTPSMSChallenge()
		return o.validateSMSOTP(o.smsChallenge, user.Human.Phone)
	}

	if o.requestType == OTPEmailRequestType {
		o.emailChallenge = session.Challenges.GetOTPEmailChallenge()
		return o.validateEmailOTP(session.Challenges.GetOTPEmailChallenge(), user.Human.Email)
	}

	return zerrors.ThrowInvalidArgument(nil, "DOM-oeUlih", "invalid OTP request type")
}

func (o *OTPCheckCommand) validateEmailOTP(otpChallenge *SessionChallengeOTPEmail, userEmail HumanEmail) error {
	if otpChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-2DlM76", "no OTP Email challenge set")
	}

	if userEmail.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-2uf0SY", "Errors.User.MFA.OTP.NotReady")
	}

	// TODO(IAM-Marco): otpChallenge.GeneratorID doesn't exist on EmailChallenge, I assume it's not a mistake
	if otpChallenge.Code == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-RegOgD", "Errors.User.Code.NotFound")
	}

	return nil
}

func (o *OTPCheckCommand) validateSMSOTP(otpChallenge *SessionChallengeOTPSMS, userPhone *HumanPhone) (err error) {
	if otpChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-UpslUc", "no OTP SMS challenge set")
	}

	if userPhone == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-fzSWTO", "no phone set")
	}

	if userPhone.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-iJZ4jp", "Errors.User.MFA.OTP.NotReady")
	}

	if otpChallenge.Code == nil && otpChallenge.GeneratorID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-tAK4Cc", "Errors.User.Code.NotFound")
	}

	return nil
}

var _ Commander = (*OTPCheckCommand)(nil)
var _ Transactional = (*OTPCheckCommand)(nil)
