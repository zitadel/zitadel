package domain

import (
	"context"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CheckTOTPType struct {
	Code string
}

type TOTPCheckCommand struct {
	CheckTOTP           *CheckTOTPType
	TarpitFunc          tarpitFn
	ValidateFunc        totpValidateFn
	EncryptionAlgorithm crypto.EncryptionAlgorithm
	SessionID           string
	InstanceID          string

	FetchedUser User

	// For Events()
	IsCheckSuccessful bool
	IsUserLocked      bool
	CheckedAt         time.Time
}

// NewTOTPCheckCommand initializes a new [TOTPCheckCommand]
//
// If tarpitFunc is nil, the default tarpit will be used.
//
// totpValidator is a function that takes as input a target TOTP to verify
// and an input ciphered secret.
// The secret is deciphered first using the input encryptionAlgo,
// then it is used to verify the TOTP. It returns true if the TOTP is validated successfully.
//
//   - totpValidator defaults to [totp.Validate]
//   - encryptionAlgo defaults to [crypto.NewAESCrypto] using the config specified in defaults.yaml
func NewTOTPCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, totpValidator totpValidateFn, encryptionAlgo crypto.EncryptionAlgorithm, request *CheckTOTPType) (*TOTPCheckCommand, error) {
	if sysConfig.Tarpit.Tarpit() == nil && tarpitFunc == nil {
		return nil, zerrors.ThrowInternal(nil, "DOM-o46bLe", "no tarpit function set")
	}

	cmd := &TOTPCheckCommand{
		CheckTOTP:           request,
		TarpitFunc:          sysConfig.Tarpit.Tarpit(),
		EncryptionAlgorithm: mfaEncryptionAlgo,
		SessionID:           sessionID,
		InstanceID:          instanceID,
		ValidateFunc:        totp.Validate,
	}
	if tarpitFunc != nil {
		cmd.TarpitFunc = tarpitFunc
	}

	if encryptionAlgo != nil {
		cmd.EncryptionAlgorithm = encryptionAlgo
	}

	if totpValidator != nil {
		cmd.ValidateFunc = totpValidator
	}

	return cmd, nil
}

// RequiresTransaction implements [Transactional].
func (t *TOTPCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (t *TOTPCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if t.CheckTOTP == nil {
		return nil, nil
	}

	events := make([]eventstore.Command, 1, 2)
	userAgg := &user.NewAggregate(t.FetchedUser.ID, t.FetchedUser.OrganizationID).Aggregate
	if t.IsCheckSuccessful {
		events[0] = user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, nil)
		return append(events, session.NewTOTPCheckedEvent(ctx, &session.NewAggregate(t.SessionID, t.InstanceID).Aggregate, t.CheckedAt)), nil
	}
	events[0] = user.NewHumanOTPCheckFailedEvent(ctx, userAgg, nil)

	if t.IsUserLocked {
		events = append(events, user.NewUserLockedEvent(ctx, userAgg))
	}

	return events, nil
}

// Execute implements [Commander].
func (t *TOTPCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if t.CheckTOTP == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	humanRepo := opts.userRepo.Human()

	verifyErr := t.verifyTOTP(t.FetchedUser.Human.TOTP.Secret)

	t.CheckedAt = time.Now()
	if verifyErr == nil {
		rowCount, err := humanRepo.Update(ctx, opts.DB(),
			humanRepo.PrimaryKeyCondition(t.InstanceID, t.FetchedUser.ID),
			humanRepo.SetLastSuccessfulTOTPCheck(t.CheckedAt),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-aoMAzO", "user"); err != nil {
			return err
		}

		rowCount, err = sessionRepo.Update(ctx, opts.DB(),
			sessionRepo.PrimaryKeyCondition(t.InstanceID, t.SessionID),
			sessionRepo.SetFactor(&SessionFactorTOTP{LastVerifiedAt: t.CheckedAt}),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-ymhCTD", "session"); err != nil {
			return err
		}

		t.IsCheckSuccessful = true
		return nil
	}

	changes := make(database.Changes, 1, 2)
	changes[0] = humanRepo.IncrementTOTPFailedAttempts()

	policy, err := GetLockoutPolicy(ctx, opts.DB(), opts.lockoutSettingRepo, t.InstanceID, t.FetchedUser.OrganizationID)
	if err != nil {
		return err
	}

	if policy != nil &&
		policy.MaxOTPAttempts != nil && *policy.MaxOTPAttempts > 0 &&
		uint64(t.FetchedUser.Human.TOTP.FailedAttempts)+1 >= *policy.MaxOTPAttempts {
		changes = append(changes, humanRepo.SetState(UserStateLocked))
		t.IsUserLocked = true
	}

	rowCount, err := humanRepo.Update(ctx, opts.DB(), humanRepo.PrimaryKeyCondition(t.InstanceID, t.FetchedUser.ID), changes)
	if err := handleUpdateError(err, 1, rowCount, "DOM-lQLpIa", "user"); err != nil {
		return err
	}

	rowCount, err = sessionRepo.Update(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(t.InstanceID, t.SessionID),
		sessionRepo.SetFactor(&SessionFactorTOTP{LastFailedAt: t.CheckedAt}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-rSa1yU", "session"); err != nil {
		return err
	}

	t.TarpitFunc(uint64(t.FetchedUser.Human.TOTP.FailedAttempts) + 1)

	// TODO(IAM-Marco): This error is a functional error and needs to be returned BUT
	// the transaction needs to NOT be rollbacked.
	//
	// As of now, this check doesn't work because the error will rollback the transaction that is
	// managed automatically by implementing the [Transactional] interface.
	//
	// Not implementing the [Transactional] interface and managing it manually is not possible either
	// because emitting the events (in [TOTPCheckCommand.Events]) need to happen in the same transaction.
	//
	// As of now, we do not have a solution for this.
	return verifyErr
}

// String implements [Commander].
func (t *TOTPCheckCommand) String() string {
	return "TOTPCheckCommand"
}

// Validate implements [Commander].
func (t *TOTPCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if t.CheckTOTP == nil {
		return nil
	}

	if t.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-ZNWO80", "Errors.Missing.SessionID")
	}
	if t.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-47G8S3", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.PrimaryKeyCondition(t.InstanceID, t.SessionID)))
	if err := handleGetError(err, "DOM-e4OuhO", "session"); err != nil {
		return err
	}
	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing")
	}

	user, err := userRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			userRepo.PrimaryKeyCondition(t.InstanceID, session.UserID),
		),
		// TODO(IAM-Marco): This might not work if we do manual transaction management. See https://github.com/zitadel/zitadel/pull/11886#discussion_r3014948862
		database.WithResultLock(),
	)
	if err := handleGetError(err, "DOM-PZvWq0", "user"); err != nil {
		return err
	}
	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-zzv1MO", "Errors.User.NotHuman")
	}

	if user.Human.TOTP == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-V6Av2a", "Errors.User.NoTOTP")
	}

	if user.Human.TOTP.Secret == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-b44CWR", "Errors.User.NoTOTPSecret")
	}

	if user.Human.TOTP.VerifiedAt.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-0g4ZAU", "Errors.User.MFA.OTP.NotReady")
	}

	if user.State == UserStateLocked {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-gM4SUh", "Errors.User.Locked")
	}

	t.FetchedUser = *user

	return nil
}

func (t *TOTPCheckCommand) verifyTOTP(existingTOTPSecret *crypto.CryptoValue) error {
	decryptedSecret, err := crypto.DecryptString(existingTOTPSecret, t.EncryptionAlgorithm)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-Yqhggx", "Errors.TOTP.FailedToDecryptSecret")
	}

	isValid := t.ValidateFunc(t.CheckTOTP.Code, decryptedSecret)
	if !isValid {
		return zerrors.ThrowInvalidArgument(nil, "DOM-o5cVir", "Errors.User.MFA.OTP.InvalidCode")
	}

	return nil
}

var _ Commander = (*TOTPCheckCommand)(nil)
var _ Transactional = (*TOTPCheckCommand)(nil)
