package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	expectedUpdatedRows = 1
	objectTypeSession   = "Session"
	objectTypeUser      = "User"
)

var _ Commander = (*OTPSMSChallengeCommand)(nil)
var _ Transactional = (*OTPSMSChallengeCommand)(nil)

type ChallengeTypeOTPSMS struct {
	ReturnCode bool
}

// (@grvijayan) todo: to be implemented
// getActiveSMSProviderFn helps determine whether to generate an internal OTP code or use an external SMS provider.
// maybe temporary until we figure out how to integrate SMS providers from v1 API
type getActiveSMSProviderFn func(ctx context.Context, instanceID string) (string, error)

// OTPSMSChallengeCommand creates an OTP SMS challenge for a session.
type OTPSMSChallengeCommand struct {
	ChallengeTypeOTPSMS *ChallengeTypeOTPSMS

	SessionID  string
	InstanceID string

	defaultSecretGeneratorConfig *crypto.GeneratorConfig
	otpAlgorithm                 crypto.EncryptionAlgorithm
	smsProvider                  getActiveSMSProviderFn
	newPhoneCode                 newOTPCodeFunc

	session *Session
	user    *User

	challengeOTPSMS *SessionChallengeOTPSMS // the generated OTP SMS challenge that is stored in the session.
	otpSMSChallenge *string                 // challenge to be set in the CreateSessionResponse
}

// NewOTPSMSChallengeCommand creates a command to generate an OTP SMS challenge for a session.
func NewOTPSMSChallengeCommand(
	challengeTypeOTPSMS *ChallengeTypeOTPSMS,
	sessionID string,
	instanceID string,
	secretGeneratorConfig *crypto.GeneratorConfig,
	otpAlgorithm crypto.EncryptionAlgorithm,
	smsProvider getActiveSMSProviderFn,
	newPhoneCodeFn newOTPCodeFunc) *OTPSMSChallengeCommand {

	if secretGeneratorConfig == nil {
		secretGeneratorConfig = otpSMSSecretGeneratorConfig
	}
	if otpAlgorithm == nil {
		otpAlgorithm = mfaEncryptionAlgo
	}
	if newPhoneCodeFn == nil {
		newPhoneCodeFn = crypto.NewCode
	}

	return &OTPSMSChallengeCommand{
		ChallengeTypeOTPSMS:          challengeTypeOTPSMS,
		SessionID:                    sessionID,
		InstanceID:                   instanceID,
		defaultSecretGeneratorConfig: secretGeneratorConfig,
		otpAlgorithm:                 otpAlgorithm,
		smsProvider:                  smsProvider,
		newPhoneCode:                 newPhoneCodeFn,
	}
}

// Validate implements [Commander].
// It checks if the command has all required fields and fetches necessary data.
func (o *OTPSMSChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.ChallengeTypeOTPSMS == nil {
		return nil
	}

	err = o.validatePreConditions()
	if err != nil {
		return err
	}

	// get session
	sessionRepo := opts.sessionRepo
	retrievedSession, err := sessionRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(sessionRepo.PrimaryKeyCondition(o.InstanceID, o.SessionID)),
	)
	if err := handleGetError(err, "DOM-2aGWWE", objectTypeSession); err != nil {
		return err
	}
	if retrievedSession.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-Vi16Fs", "Errors.Missing.Session.UserID")
	}

	// get user
	userRepo := opts.userRepo
	retrievedUser, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.PrimaryKeyCondition(o.InstanceID, retrievedSession.UserID)),
	)
	if err := handleGetError(err, "DOM-3aGHDs", objectTypeUser); err != nil {
		return err
	}

	if retrievedUser.ID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-1bzvsh", "Errors.User.UserIDMissing")
	}

	// validate human user and user phone
	if retrievedUser.Human == nil || retrievedUser.Human.Phone == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "Errors.NotFound.User.Human.Phone")
	}
	// validate phone OTP is enabled
	if retrievedUser.Human.Phone.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4m", "Errors.OTPSMS.NotEnabled")
	}
	o.session = retrievedSession
	o.user = retrievedUser

	return nil
}

// Execute implements [Commander].
// It generates the OTP SMS challenge and updates the session.
func (o *OTPSMSChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	if o.ChallengeTypeOTPSMS == nil {
		return nil
	}

	// generate phone code
	code, plain, generatorID, expiry, err := o.createPhoneCode(ctx, opts)
	if err != nil {
		return err
	}

	// update the session with the otp sms challenge
	challengeOTPSMS := &SessionChallengeOTPSMS{
		LastChallengedAt:  time.Now(),
		Code:              code,
		Expiry:            expiry,
		CodeReturned:      o.ChallengeTypeOTPSMS.ReturnCode,
		GeneratorID:       generatorID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
	sessionRepo := opts.sessionRepo
	updateCount, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.PrimaryKeyCondition(o.InstanceID, o.SessionID),
		sessionRepo.SetChallenge(challengeOTPSMS),
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-AigB0Z", "session update failed")
	}
	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-QThZH7", "Errors.Session.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(expectedUpdatedRows, updateCount), "DOM-gYp8tG", "unexpected number of rows")
	}
	// todo (@grvijayan): uncomment after these changes are available
	// if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-AigB0Z", objectTypeSession); err != nil {
	//	return err
	// }
	o.challengeOTPSMS = challengeOTPSMS
	if o.ChallengeTypeOTPSMS.ReturnCode {
		o.otpSMSChallenge = &plain
	}

	return nil
}

// Events implements [eventstore.Command].
// It returns an OTPSMSChallengedEvent.
func (o *OTPSMSChallengeCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	if o.ChallengeTypeOTPSMS == nil {
		return nil, nil
	}

	return []eventstore.Command{
		session.NewOTPSMSChallengedEvent(
			ctx,
			&session.NewAggregate(o.SessionID, o.InstanceID).Aggregate,
			o.challengeOTPSMS.Code,
			o.challengeOTPSMS.Expiry,
			o.challengeOTPSMS.CodeReturned,
			o.challengeOTPSMS.GeneratorID,
		),
	}, nil
}

func (o *OTPSMSChallengeCommand) GetOTPSMSChallenge() *string {
	return o.otpSMSChallenge
}

// String implements [Commander].
func (o *OTPSMSChallengeCommand) String() string {
	return "OTPSMSChallengeCommand"
}

// RequiresTransaction implements [Transactional].
func (o *OTPSMSChallengeCommand) RequiresTransaction() {}

// createPhoneCode generates an OTP code or retrieves the external provider ID if an external SMS provider is active.
// In the case of an external provider, the code generation is skipped (e.g., when using Twilio verification API).
func (o *OTPSMSChallengeCommand) createPhoneCode(ctx context.Context, opts *InvokeOpts) (code *crypto.CryptoValue, plain string, externalID string, expiry time.Duration, err error) {
	externalID, err = o.smsProvider(ctx, o.InstanceID)
	if err != nil {
		return nil, "", "", expiry, err
	}
	if externalID != "" {
		return nil, "", externalID, expiry, nil
	}

	config, err := GetOTPCryptoGeneratorConfigWithDefault(ctx, o.InstanceID, opts, o.defaultSecretGeneratorConfig, OTPTypeSMS)
	if err != nil {
		return nil, "", "", expiry, err
	}
	codeGenerator := crypto.NewEncryptionGenerator(*config, o.otpAlgorithm)
	crypted, plain, err := o.newPhoneCode(codeGenerator)
	if err != nil {
		return nil, "", "", expiry, err
	}
	return crypted, plain, "", config.Expiry, nil
}

func (o *OTPSMSChallengeCommand) validatePreConditions() error {
	// validate required fields
	if o.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-3XpM6A", "Errors.Missing.SessionID")
	}
	if o.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-jNNJ9f", "Errors.Missing.InstanceID")
	}

	// validate that sms provider is set
	if o.smsProvider == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-1aGWeE", "missing sms provider")
	}

	// validate that default secret generator config is set
	if o.defaultSecretGeneratorConfig == nil {
		return zerrors.ThrowInternal(nil, "DOM-IDcOzP", "missing default secret generator config")
	}

	// validate that otp algorithm is set
	if o.otpAlgorithm == nil {
		return zerrors.ThrowInternal(nil, "DOM-VeBPmV", "missing MFA encryption algorithm")
	}
	return nil
}
