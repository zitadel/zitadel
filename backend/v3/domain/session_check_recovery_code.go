package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Commander = &RecoveryCodeCheckCommand{}
var _ Transactional = &RecoveryCodeCheckCommand{}

type CheckRecoveryCode struct {
	RecoveryCode string
}

type RecoveryCodeCheckCommand struct {
	CheckRecoveryCode *CheckRecoveryCode
	Hasher            *crypto.Hasher

	SessionID  string
	InstanceID string

	session *Session
	user    *User

	checkSucceeded     bool
	userLocked         bool
	checkedAt          time.Time
	hashedRecoveryCode string
}

// NewRecoveryCodeCheckCommand returns a check Commander validating the input recovery code
//
// It assumes that a [Session] already exists: this check should be part of the
// batch call to create/set a session.
//
// The check will update the existing session or return an error if the session
// is not found or validation fails.
func NewRecoveryCodeCheckCommand(sessionID, instanceID string, check *CheckRecoveryCode, hasher *crypto.Hasher) *RecoveryCodeCheckCommand {
	return &RecoveryCodeCheckCommand{
		SessionID:         sessionID,
		InstanceID:        instanceID,
		CheckRecoveryCode: check,
		Hasher:            hasher,
	}
}

// Validate implements [Commander].
func (rc *RecoveryCodeCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	// no recovery code check to be done
	if rc.CheckRecoveryCode == nil {
		return nil
	}

	// precondition checks
	if rc.CheckRecoveryCode.RecoveryCode == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-cEKxoG", "Errors.user.MFA.RecoveryCodes.Empty")
	}
	if rc.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-hsbVyd", "Errors.session.IDMissing")
	}
	if rc.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-lGIe1v", "Errors.Instance.IDMissing")
	}

	// todo: review permission check

	// get session
	sessionRepo := opts.sessionRepo
	retrievedSession, err := sessionRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(sessionRepo.PrimaryKeyCondition(rc.InstanceID, rc.SessionID)),
	)
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-Ot3qO6", "Errors.session.NotFound")
		}
		return zerrors.ThrowInternal(err, "DOM-2sF2kF", "Errors.Internal")
	}

	if retrievedSession.UserID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-EaLqwq", "Errors.user.UserIDMissing")
	}

	// get the user from the session
	userRepo := opts.userRepo
	retrievedUser, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(retrievedSession.UserID)),
	)
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-Ot3qO6", "Errors.user.NotFound")
		}
		return zerrors.ThrowInternal(err, "DOM-7sWTNf", "Errors.Internal")
	}

	// check user state
	// todo: review: checked twice in the eventstore model
	if retrievedUser.State == UserStateLocked {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-47H1Ii", "Errors.user.Locked")
	}

	// check user MFA state
	// todo: review MFA state check
	if retrievedUser.Human == nil || retrievedUser.Human.RecoveryCodes == nil || len(retrievedUser.Human.RecoveryCodes.Codes) == 0 {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.user.MFA.RecoveryCodes.NotReady")
	}

	rc.session = retrievedSession
	rc.user = retrievedUser

	return nil
}

// Execute implements [Commander].
func (rc *RecoveryCodeCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	// no recovery code check to be done
	if rc.CheckRecoveryCode == nil {
		return nil
	}

	hashedRecoveryCode, checkErr := validateRecoveryCode(rc.CheckRecoveryCode.RecoveryCode, rc.user.Human.RecoveryCodes.Codes, rc.Hasher)
	if checkErr != nil {
		err = rc.handleRecoveryCodeCheckFailed(ctx, opts)
		if err != nil {
			return errors.Join(checkErr, err) // todo: or just log the update error?
		}
		return checkErr
	}

	return rc.handleRecoveryCodeCheckSucceeded(ctx, opts, hashedRecoveryCode)
}

func (rc *RecoveryCodeCheckCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	// no recovery code check to be done
	if rc.CheckRecoveryCode == nil {
		return nil, nil
	}

	if !rc.checkSucceeded {
		events := make([]eventstore.Command, 0, 2)
		events = append(events,
			user.NewHumanRecoveryCodeCheckFailedEvent(
				ctx,
				&user.NewAggregate(rc.user.ID, rc.user.OrganizationID).Aggregate,
				nil,
			))
		if rc.userLocked {
			events = append(events,
				user.NewUserLockedEvent(
					ctx,
					&user.NewAggregate(rc.user.ID, rc.user.OrganizationID).Aggregate,
				),
			)
		}
		return events, nil
	}

	return []eventstore.Command{
		user.NewHumanRecoveryCodeCheckSucceededEvent(
			ctx,
			&user.NewAggregate(rc.user.ID, rc.user.OrganizationID).Aggregate,
			rc.hashedRecoveryCode,
			nil,
		),
		session.NewRecoveryCodeCheckedEvent(
			ctx,
			&session.NewAggregate(rc.SessionID, rc.InstanceID).Aggregate,
			rc.checkedAt,
		),
	}, nil
}

// RequiresTransaction implements [Transactional].
func (rc *RecoveryCodeCheckCommand) RequiresTransaction() {}

// String implements [Commander].
func (rc *RecoveryCodeCheckCommand) String() string {
	return "RecoveryCodeCheckCommand"
}

func (rc *RecoveryCodeCheckCommand) handleRecoveryCodeCheckFailed(ctx context.Context, opts *InvokeOpts) error {
	checkTime := time.Now()

	lockoutPolicy, err := getLockoutPolicy(ctx, opts, rc.InstanceID, rc.user.OrganizationID)
	logging.OnError(err).Error("failed to get lockout policy") // todo: review

	// update user state and recovery_code_failed_attempts
	userRepo := opts.userRepo.Human()
	userUpdates := make([]database.Change, 0, 2)

	// update recovery_code_failed_attempts for the user
	userUpdates = append(userUpdates, userRepo.IncrementRecoveryCodeFailedAttempts())

	// update the user's state to locked if the failed recovery code check attempts exceed the configured max value
	if lockoutPolicy != nil &&
		lockoutPolicy.MaxOTPAttempts != nil &&
		*lockoutPolicy.MaxOTPAttempts > 0 &&
		(uint64(rc.user.Human.RecoveryCodes.FailedAttempts)+1 >= *lockoutPolicy.MaxOTPAttempts) {
		userUpdates = append(userUpdates, userRepo.SetState(UserStateLocked))
		rc.userLocked = true
	}

	err = rc.updateUser(
		ctx,
		opts,
		userUpdates...,
	)
	if err != nil {
		return err
	}

	// update the session's recovery code check factor
	recoveryCodeCheckFactor := &SessionFactorRecoveryCode{
		LastFailedAt: checkTime,
	}
	sessionRepo := opts.sessionRepo
	err = rc.updateSession(
		ctx,
		opts,
		sessionRepo.SetFactor(recoveryCodeCheckFactor),
	)
	if err != nil {
		return err
	}

	rc.checkedAt = checkTime
	return nil
}

func (rc *RecoveryCodeCheckCommand) handleRecoveryCodeCheckSucceeded(ctx context.Context, opts *InvokeOpts, hashedRecoveryCode string) error {
	checkTime := time.Now()

	// update the recovery_code_last_successful_check timestamp for the user
	// and remove the used recovery code from the user's list of recovery codes
	userRepo := opts.userRepo.Human()
	err := rc.updateUser(
		ctx,
		opts,
		userRepo.SetLastSuccessfulRecoveryCodeCheck(checkTime),
		userRepo.RemoveRecoveryCode(hashedRecoveryCode),
	)
	if err != nil {
		return err
	}

	// update the session's recovery code check factor
	recoveryCodeCheckFactor := &SessionFactorRecoveryCode{
		LastVerifiedAt: checkTime,
	}
	sessionRepo := opts.sessionRepo
	err = rc.updateSession(
		ctx,
		opts,
		sessionRepo.SetFactor(recoveryCodeCheckFactor),
	)
	if err != nil {
		return err
	}

	rc.checkSucceeded = true
	rc.hashedRecoveryCode = hashedRecoveryCode
	rc.checkedAt = checkTime

	return nil
}

func (rc *RecoveryCodeCheckCommand) updateSession(ctx context.Context, opts *InvokeOpts, changes ...database.Change) error {
	sessionRepo := opts.sessionRepo
	updateCount, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.PrimaryKeyCondition(rc.InstanceID, rc.SessionID),
		changes...,
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-fhR4N3", "session update failed") // todo: review error message format
	}
	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-Frnwxt", "Errors.session.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-gYp8tG", "unexpected number of rows updateCount")
	}
	return nil
}

func (rc *RecoveryCodeCheckCommand) updateUser(ctx context.Context, opts *InvokeOpts, changes ...database.Change) error {
	userRepo := opts.userRepo.Human()
	updateCount, err := userRepo.Update(
		ctx,
		opts.DB(),
		userRepo.PrimaryKeyCondition(rc.InstanceID, rc.user.ID),
		changes...,
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-XGf3Tk", "user update failed")
	}
	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-hQu5ns", "Errors.user.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-EWjTOH", "Errors.Internal")
	}
	return nil
}

// todo: review: duplicated from the domain layer
func validateRecoveryCode(reqRecoveryCode string, recoveryCodes []string, hasher *crypto.Hasher) (string, error) {
	if reqRecoveryCode == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "DOM-dk1MaX", "Errors.user.MFA.RecoveryCodes.InvalidCode")
	}
	for _, recoveryCode := range recoveryCodes {
		_, verifyErr := hasher.Verify(recoveryCode, reqRecoveryCode)
		if verifyErr != nil {
			continue
		}
		return recoveryCode, nil
	}
	return "", zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.user.MFA.RecoveryCodes.InvalidCode")
}
