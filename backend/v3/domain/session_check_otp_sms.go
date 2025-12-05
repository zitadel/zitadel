package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type otpCodeVerifyFn func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue, verificationCode string, algorithm crypto.EncryptionAlgorithm) error
type phoneCodeVerifyFn func(ctx context.Context, id string) (senders.CodeGenerator, error)

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

	smsChallenge *SessionChallengeOTPSMS
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

// RequiresTransaction implements Transactional.
func (o *OTPCheckCommand) RequiresTransaction() {}

// Events implements Commander.
func (o *OTPCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if o.CheckOTP == nil {
		return nil, nil
	}
	return nil, nil
}

// Execute implements Commander.
func (o *OTPCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.CheckOTP == nil {
		return nil
	}

	var verificationError error
	if o.requestType == OTPSMSRequestType {
		verificationError = o.executeSMSCheck(ctx)

	}
	// TODO(IAM-Marco): email check

	if verificationError == nil {
		return nil
	}
	return nil
}

func (o *OTPCheckCommand) executeSMSCheck(ctx context.Context) (err error) {
	if o.smsChallenge == nil {
		return zerrors.ThrowInternal(nil, "DOM-asd", "no sms challenge set")
	}

	if o.smsChallenge.GeneratorID == "" {
		if o.smsChallenge.Code == nil {
			return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.Code.NotFound")
		}
		return o.otpCodeVerifyFunc(
			time.Now(), // o.smsChallenge is missing creation date
			o.smsChallenge.Expiry,
			o.smsChallenge.Code,
			o.CheckOTP.GetCode(),
			o.encryptionAlgorithm,
		)
	}

	if o.phonecodeVerifyFunc == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.Code.NotConfigured")
	}
	verifier, err := o.phonecodeVerifyFunc(ctx, o.smsChallenge.GeneratorID)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-asd", "failed fetching phone verifier")
	}

	// TODO(IAM-Marco): o.smsChallenge should have a VerificationID
	return verifier.VerifyCode("Missing Verification ID", o.CheckOTP.GetCode())
}

// String implements Commander.
func (o *OTPCheckCommand) String() string {
	return "OTPCheckCommand"
}

// Validate implements Commander.
func (o *OTPCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if o.CheckOTP == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(o.sessionID)))
	if err := handleGetError(err, "DOM-asd", "session"); err != nil {
		return err
	}

	user, err := userRepo.Get(ctx, opts.DB(),
		database.WithCondition(userRepo.IDCondition(session.UserID)),
	)
	if err := handleGetError(err, "DOM-asd", "user"); err != nil {
		return err
	}

	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "user not human")
	}

	if o.CheckOTP.GetCode() == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-asd", "Errors.User.Code.Empty")
	}

	if o.requestType == OTPSMSRequestType {
		return o.validateSMSOTP(session.Challenges.GetOTPSMSChallenge(), user.Human.Phone)
	}

	// TODO(IAM-Marco): Email check
	return nil
}

func (o *OTPCheckCommand) validateSMSOTP(otpChallenge *SessionChallengeOTPSMS, userPhone *HumanPhone) (err error) {
	if otpChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "no OTP SMS challenge set")
	}

	if userPhone == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "no phone set")
	}

	if userPhone.OTP.EnabledAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.MFA.OTP.NotReady")
	}

	if otpChallenge.Code == nil && otpChallenge.GeneratorID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.Code.NotFound")
	}

	return nil
}

var _ Commander = (*OTPCheckCommand)(nil)
var _ Transactional = (*OTPCheckCommand)(nil)
