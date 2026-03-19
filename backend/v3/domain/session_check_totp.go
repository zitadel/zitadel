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
	tarpitFunc          tarpitFn
	validateFunc        totpValidateFn
	encryptionAlgorithm crypto.EncryptionAlgorithm

	sessionID  string
	instanceID string

	FetchedUser User

	// For Events()
	IsCheckSuccessful bool
	IsUserLocked      bool
	CheckedAt         time.Time
}

func NewTOTPCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, totpValidator totpValidateFn, encryptionAlgo crypto.EncryptionAlgorithm, request *CheckTOTPType) *TOTPCheckCommand {
	tf := sysConfig.Tarpit.Tarpit()
	if tarpitFunc != nil {
		tf = tarpitFunc
	}

	ea := mfaEncryptionAlgo
	if encryptionAlgo != nil {
		ea = encryptionAlgo
	}

	totpValidateFunc := totp.Validate
	if totpValidator != nil {
		totpValidateFunc = totpValidator
	}

	return &TOTPCheckCommand{
		CheckTOTP:           request,
		tarpitFunc:          tf,
		encryptionAlgorithm: ea,
		sessionID:           sessionID,
		instanceID:          instanceID,
		validateFunc:        totpValidateFunc,
	}
}

// RequiresTransaction implements [Transactional].
func (t *TOTPCheckCommand) RequiresTransaction() {}

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
		toReturn = append(toReturn, session.NewTOTPCheckedEvent(ctx, &session.NewAggregate(t.sessionID, t.instanceID).Aggregate, t.CheckedAt))
	} else {
		toReturn[1] = session.NewTOTPCheckedEvent(ctx, &session.NewAggregate(t.sessionID, t.instanceID).Aggregate, t.CheckedAt)
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

	if verifyErr == nil {
		t.CheckedAt = time.Now()
		rowCount, err := humanRepo.Update(ctx, opts.DB(),
			humanRepo.PrimaryKeyCondition(t.instanceID, t.FetchedUser.ID),
			humanRepo.SetLastSuccessfulTOTPCheck(t.CheckedAt),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-aoMAzO", "user"); err != nil {
			return err
		}

		rowCount, err = sessionRepo.Update(ctx, opts.DB(),
			sessionRepo.PrimaryKeyCondition(t.instanceID, t.sessionID),
			sessionRepo.SetFactor(&SessionFactorTOTP{LastVerifiedAt: t.CheckedAt}),
		)
		if err := handleUpdateError(err, 1, rowCount, "DOM-ymhCTD", "session"); err != nil {
			return err
		}

		t.IsCheckSuccessful = true
		return nil
	}

	t.CheckedAt = time.Now()
	changes := make(database.Changes, 1, 2)
	changes[0] = humanRepo.IncrementTOTPFailedAttempts()

	policy, err := getLockoutPolicy(ctx, opts.DB(), opts.lockoutSettingRepo, t.instanceID, t.FetchedUser.OrganizationID)
	if err != nil {
		return err
	}

	if policy != nil &&
		policy.MaxOTPAttempts != nil && *policy.MaxOTPAttempts > 0 &&
		uint64(t.FetchedUser.Human.TOTP.FailedAttempts+1) >= *policy.MaxOTPAttempts {
		changes = append(changes, humanRepo.SetState(UserStateLocked))
		t.IsUserLocked = true
	}

	rowCount, err := humanRepo.Update(ctx, opts.DB(), humanRepo.PrimaryKeyCondition(t.instanceID, t.FetchedUser.ID), changes)
	if err := handleUpdateError(err, 1, rowCount, "DOM-lQLpIa", "user"); err != nil {
		return err
	}

	rowCount, err = sessionRepo.Update(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(t.instanceID, t.sessionID),
		sessionRepo.SetFactor(&SessionFactorTOTP{LastVerifiedAt: t.CheckedAt}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-rSa1yU", "session"); err != nil {
		return err
	}

	t.tarpitFunc(uint64(t.FetchedUser.Human.TOTP.FailedAttempts + 1))

	// TODO(IAM-Marco): The error returned here will block the transaction and stop events from being emitted.
	// This error is functional so it should be returned AND the transaction should succeed. How can we fix it?
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

	if t.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-ZNWO80", "Errors.Missing.SessionID")
	}
	if t.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-47G8S3", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.PrimaryKeyCondition(t.instanceID, t.sessionID)))
	if err := handleGetError(err, "DOM-e4OuhO", "session"); err != nil {
		return err
	}
	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing")
	}

	user, err := userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.PrimaryKeyCondition(t.instanceID, session.UserID)))
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
	decryptedSecret, err := crypto.DecryptString(existingTOTPSecret, t.encryptionAlgorithm)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-Yqhggx", "failed decrypting TOTP secret")
	}

	isValid := t.validateFunc(t.CheckTOTP.Code, decryptedSecret)
	if !isValid {
		return zerrors.ThrowInvalidArgument(nil, "DOM-o5cVir", "Errors.User.MFA.OTP.InvalidCode")
	}

	return nil
}

var _ Commander = (*TOTPCheckCommand)(nil)
var _ Transactional = (*TOTPCheckCommand)(nil)
