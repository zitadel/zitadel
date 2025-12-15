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
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var _ Commander = (*OTPSMSChallengeCommand)(nil)
var _ Transactional = (*OTPSMSChallengeCommand)(nil)

type newOTPCodeFunc func(g crypto.Generator) (*crypto.CryptoValue, string, error)

// (@grvijayan) todo: to be implemented
// getActiveSMSProviderFn helps determine whether to generate an internal OTP code or use an external SMS provider.
// maybe temporary until we figure out how to integrate SMS providers from v1 API
type getActiveSMSProviderFn func(ctx context.Context, instanceID string) (string, error)

// OTPSMSChallengeCommand creates an OTP SMS challenge for a session.
type OTPSMSChallengeCommand struct {
	RequestChallengeOTPSMS *session_grpc.RequestChallenges_OTPSMS

	sessionID                    string
	instanceID                   string
	defaultSecretGeneratorConfig *crypto.GeneratorConfig
	otpAlg                       crypto.EncryptionAlgorithm
	smsProvider                  getActiveSMSProviderFn
	newPhoneCode                 newOTPCodeFunc

	Session *Session
	User    *User

	ChallengeOTPSMS *SessionChallengeOTPSMS // the generated OTP SMS challenge that is stored in the session.
	OTPSMSChallenge *string                 // challenge to be set in the CreateSessionResponse
}

// NewOTPSMSChallengeCommand creates a command to generate an OTP SMS challenge for a session.
func NewOTPSMSChallengeCommand(
	requestChallengeOTPSMS *session_grpc.RequestChallenges_OTPSMS,
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
		RequestChallengeOTPSMS:       requestChallengeOTPSMS,
		sessionID:                    sessionID,
		instanceID:                   instanceID,
		defaultSecretGeneratorConfig: secretGeneratorConfig,
		otpAlg:                       otpAlgorithm,
		smsProvider:                  smsProvider,
		newPhoneCode:                 newPhoneCodeFn,
	}
}

// Validate implements [Commander].
// It checks if the command has all required fields and fetches necessary data.
func (o *OTPSMSChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.RequestChallengeOTPSMS == nil {
		return nil
	}

	// validate required fields
	if o.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-3XpM6A", "session id missing")
	}
	if o.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-jNNJ9f", "instance id missing")
	}

	// validate that sms provider is set
	if o.smsProvider == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-1aGWeE", "sms provider not configured")
	}

	// get session
	sessionRepo := opts.sessionRepo
	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(o.sessionID)))
	if err := handleGetError(err, "DOM-2aGWWE", objectTypeSession); err != nil {
		return err
	}
	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-Vi16Fs", "missing user id in session")
	}

	// get user
	userRepo := opts.userRepo
	user, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(session.UserID)),
	)
	if err := handleGetError(err, "DOM-3aGHDs", objectTypeUser); err != nil {
		return err
	}

	// validate human user and user phone
	if user.Human == nil || user.Human.Phone == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "user phone not configured")
	}
	// validate phone OTP is enabled
	if user.Human.Phone.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4m", "phone OTP not enabled")
	}
	o.Session = session
	o.User = user
	return nil
}

// Execute implements [Commander].
// It generates the OTP SMS challenge and updates the session.
func (o *OTPSMSChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	if o.RequestChallengeOTPSMS == nil {
		return nil
	}

	returnCode := o.RequestChallengeOTPSMS.GetReturnCode()

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
		CodeReturned:      returnCode,
		GeneratorID:       generatorID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
	sessionRepo := opts.sessionRepo
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.IDCondition(o.Session.ID),
		sessionRepo.SetChallenge(challengeOTPSMS),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-AigB0Z", objectTypeSession); err != nil {
		return err
	}
	o.ChallengeOTPSMS = challengeOTPSMS
	if returnCode {
		o.OTPSMSChallenge = &plain
	}

	return nil
}

// Events implements [eventstore.Command].
// It returns an OTPSMSChallengedEvent.
func (o *OTPSMSChallengeCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if o.RequestChallengeOTPSMS == nil {
		return nil, nil
	}
	return []eventstore.Command{
		session.NewOTPSMSChallengedEvent(
			ctx,
			&session.NewAggregate(o.sessionID, o.instanceID).Aggregate,
			o.ChallengeOTPSMS.Code,
			o.ChallengeOTPSMS.Expiry,
			o.ChallengeOTPSMS.CodeReturned,
			o.ChallengeOTPSMS.GeneratorID,
		),
	}, nil
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
	externalID, err = o.smsProvider(ctx, o.instanceID)
	if err != nil {
		return nil, "", "", expiry, err
	}
	if externalID != "" {
		return nil, "", externalID, expiry, nil
	}

	config, err := getOTPCryptoGeneratorConfigWithDefault(ctx, o.instanceID, opts, o.defaultSecretGeneratorConfig, OTPSMSRequestType)
	if err != nil {
		return nil, "", "", expiry, err
	}
	codeGenerator := crypto.NewEncryptionGenerator(*config, o.otpAlg)
	crypted, plain, err := o.newPhoneCode(codeGenerator)
	if err != nil {
		return nil, "", "", expiry, err
	}
	return crypted, plain, "", config.Expiry, nil
}

func getOTPCryptoGeneratorConfigWithDefault(ctx context.Context, instanceID string, opts *InvokeOpts, defaultConfig *crypto.GeneratorConfig, otpType OTPRequestType) (*crypto.GeneratorConfig, error) {
	settingsRepo := opts.secretGeneratorSettingsRepo
	cfg, err := settingsRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(
			database.And(
				settingsRepo.InstanceIDCondition(instanceID),
				database.NewTextCondition( // (@grvijayan) todo: check TypeCondition
					settingsRepo.TypeColumn(),
					database.TextOperationEqual,
					SettingTypeSecretGenerator.String(),
				),
			),
		),
	)
	if err := handleGetError(err, "DOM-x7Yd3E", "secret_generator_settings"); err != nil {
		return nil, err
	}

	if cfg.State != SettingStateActive {
		return defaultConfig, nil
	}

	var attrs SecretGeneratorAttrsWithExpiry
	switch otpType {
	case OTPSMSRequestType:
		if cfg.OTPSMS == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPSMS.SecretGeneratorAttrsWithExpiry
	case OTPEmailRequestType:
		if cfg.OTPEmail == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPEmail.SecretGeneratorAttrsWithExpiry
	default:
		return nil, zerrors.ThrowInternal(nil, "DOM-3AcM0U", "invalid otp request type")
	}
	return &crypto.GeneratorConfig{
		Length:              *attrs.Length,
		Expiry:              *attrs.Expiry,
		IncludeLowerLetters: *attrs.IncludeLowerLetters,
		IncludeUpperLetters: *attrs.IncludeUpperLetters,
		IncludeDigits:       *attrs.IncludeDigits,
		IncludeSymbols:      *attrs.IncludeSymbols,
	}, nil
}
