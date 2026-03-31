package domain

import (
	"context"
	"io"
	"text/template"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Commander = (*OTPEmailChallengeCommand)(nil)
var _ Transactional = (*OTPEmailChallengeCommand)(nil)

type SendCode struct {
	URLTemplate string
}

type DeliveryType struct {
	SendCode   *SendCode
	ReturnCode bool
}
type ChallengeTypeOTPEmail struct {
	DeliveryType DeliveryType
}

type OTPEmailChallengeCommand struct {
	ChallengeTypeOTPEmail *ChallengeTypeOTPEmail

	SessionID  string
	InstanceID string

	defaultSecretGeneratorConfig *crypto.GeneratorConfig
	otpEncryptionAlgorithm       crypto.EncryptionAlgorithm
	newEmailCode                 newOTPCodeFunc

	sessionChallengeOTPEmail *SessionChallengeOTPEmail // the generated OTP Email challenge that is stored in the session.
	otpEmailChallenge        *string                   // challenge to be set in the CreateSessionResponse
}

func NewOTPEmailChallengeCommand(
	challengeTypeOTPEmail *ChallengeTypeOTPEmail,
	sessionID string,
	instanceID string,
	secretGeneratorConfig *crypto.GeneratorConfig,
	otpAlgorithm crypto.EncryptionAlgorithm,
	newEmailCodeFn newOTPCodeFunc) *OTPEmailChallengeCommand {

	if secretGeneratorConfig == nil {
		secretGeneratorConfig = otpEmailSecretGeneratorConfig
	}
	if otpAlgorithm == nil {
		otpAlgorithm = mfaEncryptionAlgo
	}
	if newEmailCodeFn == nil {
		newEmailCodeFn = crypto.NewCode
	}

	return &OTPEmailChallengeCommand{
		ChallengeTypeOTPEmail:        challengeTypeOTPEmail,
		SessionID:                    sessionID,
		InstanceID:                   instanceID,
		defaultSecretGeneratorConfig: secretGeneratorConfig,
		otpEncryptionAlgorithm:       otpAlgorithm,
		newEmailCode:                 newEmailCodeFn,
	}
}

// Validate implements [Commander].
// It validates that the session and user exist and that the user has email OTP enabled.
func (o *OTPEmailChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.ChallengeTypeOTPEmail == nil {
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
	if err := handleGetError(err, "DOM-JArUai", objectTypeSession); err != nil {
		return err
	}
	if retrievedSession.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-wG2XoJ", "Errors.Missing.Session.UserID")
	}

	// get user
	userRepo := opts.userRepo
	retrievedUser, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.PrimaryKeyCondition(o.InstanceID, retrievedSession.UserID)),
	)
	if err := handleGetError(err, "DOM-56MWkg", objectTypeUser); err != nil {
		return err
	}

	// validate human user and user email
	if retrievedUser.Human == nil || retrievedUser.Human.Email.Address == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "Errors.NotFound.User.Human.Email")
	}
	// validate email OTP is enabled
	if retrievedUser.Human.Email.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4q", "Errors.User.MFA.OTP.NotReady")
	}

	return nil
}

// Execute implements [Commander].
// It updates the session with the generated OTP email challenge.
func (o *OTPEmailChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	if o.ChallengeTypeOTPEmail == nil {
		return nil
	}

	// prepare the otp email challenge
	sessionChallengeOTPEmail, challenge, err := o.prepareOTPEmailChallenge(ctx, opts)
	if err != nil {
		return err
	}

	// update the session with the otp email challenge
	sessionRepo := opts.sessionRepo
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.PrimaryKeyCondition(o.InstanceID, o.SessionID),
		sessionRepo.SetChallenge(sessionChallengeOTPEmail),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-YfQIA3", objectTypeSession); err != nil {
		return err
	}
	o.sessionChallengeOTPEmail = sessionChallengeOTPEmail
	if o.ChallengeTypeOTPEmail.DeliveryType.ReturnCode { // only set when the delivery type is ReturnCode
		o.otpEmailChallenge = &challenge
	}

	return nil
}

// Events implements [Commander].
// It creates the OTPEmailChallengedEvent if an OTP email challenge was requested.
func (o *OTPEmailChallengeCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if o.ChallengeTypeOTPEmail == nil {
		return nil, nil
	}

	return []eventstore.Command{
		session.NewOTPEmailChallengedEvent(
			ctx,
			&session.NewAggregate(o.SessionID, o.InstanceID).Aggregate,
			o.sessionChallengeOTPEmail.Code,
			o.sessionChallengeOTPEmail.Expiry,
			o.sessionChallengeOTPEmail.CodeReturned,
			o.sessionChallengeOTPEmail.URLTemplate,
		),
	}, nil
}

func (o *OTPEmailChallengeCommand) GetOTPEmailChallenge() *string {
	return o.otpEmailChallenge
}

// prepareOTPEmailChallenge generates the OTP email challenge based on the delivery type in the request.
func (o *OTPEmailChallengeCommand) prepareOTPEmailChallenge(ctx context.Context, opts *InvokeOpts) (*SessionChallengeOTPEmail, string, error) {
	// generate email code
	config, err := GetOTPCryptoGeneratorConfigWithDefault(ctx, o.InstanceID, opts, o.defaultSecretGeneratorConfig, OTPTypeEmail)
	if err != nil {
		return nil, "", err
	}
	codeGenerator := crypto.NewEncryptionGenerator(*config, o.otpEncryptionAlgorithm)
	crypted, plain, err := o.newEmailCode(codeGenerator)
	if err != nil {
		return nil, "", err
	}

	challengeOTPEmail := &SessionChallengeOTPEmail{
		LastChallengedAt:  time.Now(),
		Code:              crypted,
		Expiry:            config.Expiry,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}

	var otpEmailChallenge string
	switch {
	case o.ChallengeTypeOTPEmail.DeliveryType.SendCode != nil:
		challengeOTPEmail.URLTemplate = o.ChallengeTypeOTPEmail.DeliveryType.SendCode.URLTemplate
	case o.ChallengeTypeOTPEmail.DeliveryType.ReturnCode:
		challengeOTPEmail.CodeReturned = true
		otpEmailChallenge = plain
	default:
		// no additional action needed
	}
	return challengeOTPEmail, otpEmailChallenge, nil
}

// String implements [Commander].
func (o *OTPEmailChallengeCommand) String() string {
	return "OTPEmailChallengeCommand"
}

// RequiresTransaction implements [Transactional].
func (o *OTPEmailChallengeCommand) RequiresTransaction() {}

// validateURLTemplate renders the given URL template with sample data to validate its correctness.
func validateURLTemplate(w io.Writer, tmpl string) error {
	otpEmailURLData := &struct {
		Code              string
		UserID            string
		LoginName         string
		DisplayName       string
		PreferredLanguage language.Tag
		SessionID         string
	}{
		Code:              "code",
		UserID:            "userID",
		LoginName:         "loginName",
		DisplayName:       "displayName",
		PreferredLanguage: language.English,
		SessionID:         "SessionID",
	}
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOM-wkDwQM", "Errors.Invalid.URLTemplate")
	}
	if err = parsed.Execute(w, otpEmailURLData); err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOM-F5Yv8l", "Errors.Invalid.URLTemplate")
	}
	return nil
}

func (o *OTPEmailChallengeCommand) validatePreConditions() error {
	// validate required fields
	if o.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-BQ5UgK", "Errors.Missing.SessionID")
	}
	if o.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-kDnkDn", "Errors.Missing.InstanceID")
	}

	// validate that default secret generator config is set
	if o.defaultSecretGeneratorConfig == nil {
		return zerrors.ThrowInternal(nil, "DOM-nnB9MS", "missing default secret generator config")
	}

	// validate that otp encryption algorithm is set
	if o.otpEncryptionAlgorithm == nil {
		return zerrors.ThrowInternal(nil, "DOM-kuG75Q", "missing MFA encryption algorithm")
	}

	// validate the URL template
	if sc := o.ChallengeTypeOTPEmail.DeliveryType.SendCode; sc != nil {
		if sc.URLTemplate != "" {
			if err := validateURLTemplate(io.Discard, sc.URLTemplate); err != nil {
				return err
			}
		}
	}
	return nil
}
