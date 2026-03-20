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
// and an input cyphered secret.
// The secret is decyphered first using the input encryptionAlgo,
// then it is used to verify the TOTP. It returns true if the TOTP is validated successfully.
//
//   - totpValidator defaults to [totp.Validate]
//   - encryptionAlgo defaults to [crypto.NewAESCrypto] using the config specified in defaults.yaml
//
// The command does not implement [Transactional] due totpValidator that might take a long time to execute.
// So the DB transaction will be started only after totpValidator has been run.
//
// Moreover, the command may return a functional error so manual management of the transaction is needed
// to avoid rollbacking a transaction despite having only a functional error.
func NewTOTPCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, totpValidator totpValidateFn, encryptionAlgo crypto.EncryptionAlgorithm, request *CheckTOTPType) *TOTPCheckCommand {
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

	return cmd
}

// Events implements [Commander].
func (t *TOTPCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if t.CheckTOTP == nil {
		return nil, nil
	}

	toReturn := make([]eventstore.Command, 2, 3)
	userAgg := &user.NewAggregate(t.FetchedUser.ID, t.FetchedUser.OrganizationID).Aggregate
	if t.IsCheckSuccessful {
		toReturn[0] = user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, nil)
	} else {
		toReturn[0] = user.NewHumanOTPCheckFailedEvent(ctx, userAgg, nil)
	}

	if t.IsUserLocked {
		toReturn[1] = user.NewUserLockedEvent(ctx, userAgg)
		toReturn = append(toReturn, session.NewTOTPCheckedEvent(ctx, &session.NewAggregate(t.SessionID, t.InstanceID).Aggregate, t.CheckedAt))
	} else {
		toReturn[1] = session.NewTOTPCheckedEvent(ctx, &session.NewAggregate(t.SessionID, t.InstanceID).Aggregate, t.CheckedAt)
	}

	return toReturn, nil
}

// Execute implements [Commander].
func (t *TOTPCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if t.CheckTOTP == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	humanRepo := opts.userRepo.Human()

	verifyErr := t.verifyTOTP(t.FetchedUser.Human.TOTP.Secret)

	beginner, ok := opts.DB().(database.Beginner)
	if !ok {
		return zerrors.ThrowInternal(nil, "DOM-Ug9936", "database doesn't implement database.Beginner")
	}

	tx, txErr := beginner.Begin(ctx, nil)
	if txErr != nil {
		return zerrors.ThrowInternal(txErr, "DOM-pw7gF8", "failed starting transaction")
	}

	defer func() {
		if endErr := tx.End(ctx, txErr); endErr != nil {
			err = endErr
		}
	}()

	if verifyErr == nil {
		t.CheckedAt = time.Now()
		rowCount, err := humanRepo.Update(ctx, tx,
			humanRepo.PrimaryKeyCondition(t.InstanceID, t.FetchedUser.ID),
			humanRepo.SetLastSuccessfulTOTPCheck(t.CheckedAt),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-aoMAzO", "user"); err != nil {
			txErr = err
			return err
		}

		rowCount, err = sessionRepo.Update(ctx, tx,
			sessionRepo.PrimaryKeyCondition(t.InstanceID, t.SessionID),
			sessionRepo.SetFactor(&SessionFactorTOTP{LastVerifiedAt: t.CheckedAt}),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-ymhCTD", "session"); err != nil {
			txErr = err
			return err
		}

		t.IsCheckSuccessful = true
		return nil
	}

	t.CheckedAt = time.Now()
	changes := make(database.Changes, 1, 2)
	changes[0] = humanRepo.IncrementTOTPFailedAttempts()

	policy, err := getLockoutPolicy(ctx, tx, opts.lockoutSettingRepo, t.InstanceID, t.FetchedUser.OrganizationID)
	if err != nil {
		txErr = err
		return err
	}

	if policy != nil &&
		policy.MaxOTPAttempts != nil && *policy.MaxOTPAttempts > 0 &&
		uint64(t.FetchedUser.Human.TOTP.FailedAttempts+1) >= *policy.MaxOTPAttempts {
		changes = append(changes, humanRepo.SetState(UserStateLocked))
		t.IsUserLocked = true
	}

	rowCount, err := humanRepo.Update(ctx, tx, humanRepo.PrimaryKeyCondition(t.InstanceID, t.FetchedUser.ID), changes)
	if err := handleUpdateError(err, 1, rowCount, "DOM-lQLpIa", "user"); err != nil {
		txErr = err
		return err
	}

	rowCount, err = sessionRepo.Update(ctx, tx,
		sessionRepo.PrimaryKeyCondition(t.InstanceID, t.SessionID),
		sessionRepo.SetFactor(&SessionFactorTOTP{LastVerifiedAt: t.CheckedAt}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-rSa1yU", "session"); err != nil {
		txErr = err
		return err
	}

	t.TarpitFunc(uint64(t.FetchedUser.Human.TOTP.FailedAttempts + 1))

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

	user, err := userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.PrimaryKeyCondition(t.InstanceID, session.UserID)))
	if err := handleGetError(err, "DOM-PZvWq0", "user"); err != nil {
		return err
	}
	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-zzv1MO", "user not human")

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
		return zerrors.ThrowInternal(err, "DOM-Yqhggx", "failed decrypting TOTP secret")
	}

	isValid := t.ValidateFunc(t.CheckTOTP.Code, decryptedSecret)
	if !isValid {
		return zerrors.ThrowInvalidArgument(nil, "DOM-o5cVir", "Errors.User.MFA.OTP.InvalidCode")
	}

	return nil
}

var _ Commander = (*TOTPCheckCommand)(nil)
