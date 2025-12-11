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
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var _ Commander = (*OTPEmailChallengeCommand)(nil)
var _ Transactional = (*OTPEmailChallengeCommand)(nil)

type OTPEmailChallengeCommand struct {
	RequestChallengeOTPEmail *session_grpc.RequestChallenges_OTPEmail

	sessionID                    string
	instanceID                   string
	defaultSecretGeneratorConfig *crypto.GeneratorConfig
	otpAlg                       crypto.EncryptionAlgorithm
	newEmailCode                 newOTPCodeFunc

	Session *Session
	User    *User

	ChallengeOTPEmail *SessionChallengeOTPEmail // the generated OTP SMS challenge that is stored in the session.
	OTPEmailChallenge *string                   // challenge to be set in the CreateSessionResponse
}

func NewOTPEmailChallengeCommand(
	requestChallengeOTPEmail *session_grpc.RequestChallenges_OTPEmail,
	sessionID string,
	instanceID string,
	secretGeneratorConfig *crypto.GeneratorConfig,
	otpAlgorithm crypto.EncryptionAlgorithm,
	newEmailCodeFn newOTPCodeFunc) *OTPEmailChallengeCommand {

	if secretGeneratorConfig == nil {
		secretGeneratorConfig = otpSMSSecretGeneratorConfig
	}
	if otpAlgorithm == nil {
		otpAlgorithm = mfaEncryptionAlgo
	}
	if newEmailCodeFn == nil {
		newEmailCodeFn = crypto.NewCode
	}

	return &OTPEmailChallengeCommand{
		RequestChallengeOTPEmail:     requestChallengeOTPEmail,
		sessionID:                    sessionID,
		instanceID:                   instanceID,
		defaultSecretGeneratorConfig: secretGeneratorConfig,
		otpAlg:                       otpAlgorithm,
		newEmailCode:                 newEmailCodeFn,
	}
}

// Validate implements [Commander].
// It validates that the session and user exist and that the user has email OTP enabled.
func (o *OTPEmailChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.RequestChallengeOTPEmail == nil {
		return nil
	}
	// validate required fields
	if o.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-BQ5UgK", "session id missing")
	}
	if o.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-kDnkDn", "instance id missing")
	}

	// get session
	sessionRepo := opts.sessionRepo
	o.Session, err = sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(o.sessionID)))
	if err := handleGetError(err, "DOM-JArUai", objectTypeSession); err != nil {
		return err
	}
	if o.Session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-wG2XoJ", "missing user id in session")
	}

	// get user
	userRepo := opts.userRepo
	o.User, err = userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(o.Session.UserID)),
	)
	if err := handleGetError(err, "DOM-56MWkg", objectTypeUser); err != nil {
		return err
	}

	// validate human user and user email
	if o.User.Human == nil || o.User.Human.Email.Address == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "user email not configured")
	}
	// validate email OTP is enabled
	if o.User.Human.Email.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4q", "email OTP not enabled")
	}

	return nil
}

// Events implements [Commander].
// It creates the OTPEmailChallengedEvent if an OTP email challenge was requested.
func (o *OTPEmailChallengeCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if o.RequestChallengeOTPEmail == nil {
		return nil, nil
	}
	return []eventstore.Command{
		session.NewOTPEmailChallengedEvent(
			ctx,
			&session.NewAggregate(o.sessionID, o.instanceID).Aggregate,
			o.ChallengeOTPEmail.Code,
			o.ChallengeOTPEmail.Expiry,
			o.ChallengeOTPEmail.CodeReturned,
			o.ChallengeOTPEmail.URLTmpl,
		),
	}, nil
}

// Execute implements [Commander].
// It updates the session with the generated OTP email challenge.
func (o *OTPEmailChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	if o.RequestChallengeOTPEmail == nil {
		return nil
	}

	// prepare the otp email challenge
	err := o.prepareOTPEmailChallenge(ctx, opts)
	if err != nil {
		return err
	}

	// update the session with the otp email challenge
	sessionRepo := opts.sessionRepo
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.IDCondition(o.Session.ID),
		sessionRepo.SetChallenge(o.ChallengeOTPEmail),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-YfQIA3", objectTypeSession); err != nil {
		return err
	}

	return nil
}

// prepareOTPEmailChallenge generates the OTP email challenge based on the delivery type in the request.
func (o *OTPEmailChallengeCommand) prepareOTPEmailChallenge(ctx context.Context, opts *InvokeOpts) error {
	// generate email code
	// todo: when do we set the default config?
	config, err := getOTPEmailCryptoGeneratorConfigWithDefault(ctx, o.instanceID, opts, o.defaultSecretGeneratorConfig)
	if err != nil {
		return err
	}
	codeGenerator := crypto.NewEncryptionGenerator(*config, o.otpAlg)
	crypted, plain, err := o.newEmailCode(codeGenerator)
	if err != nil {
		return err
	}

	o.ChallengeOTPEmail = &SessionChallengeOTPEmail{
		LastChallengedAt:  time.Now(),
		Code:              crypted,
		Expiry:            config.Expiry,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}

	o.OTPEmailChallenge = new(string)
	switch t := o.RequestChallengeOTPEmail.GetDeliveryType().(type) {
	case *session_grpc.RequestChallenges_OTPEmail_SendCode_:
		urlTmpl := t.SendCode.GetUrlTemplate()
		if err := renderURLTemplate(io.Discard, urlTmpl); err != nil {
			return err
		}
		o.ChallengeOTPEmail.URLTmpl = urlTmpl
	case *session_grpc.RequestChallenges_OTPEmail_ReturnCode_:
		o.ChallengeOTPEmail.CodeReturned = true
		*o.OTPEmailChallenge = plain
	case nil:
		//
	default:
		return zerrors.ThrowUnimplementedf(nil, "DOM-cc1bRa", "delivery_type oneOf %T in OTPEmailChallenge not implemented", t)
	}
	return nil
}

// String implements [Commander].
func (o *OTPEmailChallengeCommand) String() string {
	return "OTPEmailChallengeCommand"
}

// RequiresTransaction implements [Transactional].
func (o *OTPEmailChallengeCommand) RequiresTransaction() {}

// renderURLTemplate renders the given URL template with sample data to validate its correctness.
func renderURLTemplate(w io.Writer, tmpl string) error {
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
		SessionID:         "sessionID",
	}
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOM-wkDwQM", "invalid URL template")
	}
	if err = parsed.Execute(w, otpEmailURLData); err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOM-F5Yv8l", "invalid URL template")
	}
	return nil
}

func getOTPEmailCryptoGeneratorConfigWithDefault(ctx context.Context, instanceID string, opts *InvokeOpts, defaultConfig *crypto.GeneratorConfig) (*crypto.GeneratorConfig, error) {
	settingsRepo := opts.secretGeneratorSettingsRepo
	cfg, err := settingsRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(
			database.And(
				settingsRepo.InstanceIDCondition(instanceID),
				database.NewTextCondition( // todo: check TypeCondition
					settingsRepo.TypeColumn(),
					database.TextOperationEqual,
					SettingTypeSecretGenerator.String(),
				),
			),
		),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOM-x7Yd3E", "failed to get OTP email secret generator config")
	}
	if cfg.Settings.State != SettingStateActive || cfg.OTPEmail == nil {
		return defaultConfig, nil
	}
	return &crypto.GeneratorConfig{
		Length:              *cfg.OTPEmail.Length,
		Expiry:              *cfg.OTPEmail.Expiry,
		IncludeLowerLetters: *cfg.OTPEmail.IncludeLowerLetters,
		IncludeUpperLetters: *cfg.OTPEmail.IncludeUpperLetters,
		IncludeDigits:       *cfg.OTPEmail.IncludeDigits,
		IncludeSymbols:      *cfg.OTPEmail.IncludeSymbols,
	}, nil
}
