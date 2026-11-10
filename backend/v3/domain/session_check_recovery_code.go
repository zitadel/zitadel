package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Commander = &RecoveryCodeCheckCommand{}
var _ Transactional = &RecoveryCodeCheckCommand{}

type CheckTypeRecoveryCode struct {
	RecoveryCode string
}

type RecoveryCodeCheckCommand struct {
	CheckRecoveryCode *CheckTypeRecoveryCode

	SessionID  string
	InstanceID string

	verify  verifierFn
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
func NewRecoveryCodeCheckCommand(sessionID, instanceID string, check *CheckTypeRecoveryCode, verifier verifierFn) *RecoveryCodeCheckCommand {
	if verifier == nil {
		verifier = passwordHasher.Verify
	}

	return &RecoveryCodeCheckCommand{
		SessionID:         sessionID,
		InstanceID:        instanceID,
		CheckRecoveryCode: check,
		verify:            verifier,
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
		return zerrors.ThrowInvalidArgument(nil, "DOM-cEKxoG", "Errors.User.MFA.RecoveryCodes.Empty")
	}
	if rc.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-hsbVyd", "Errors.Session.IDMissing")
	}
	if rc.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-lGIe1v", "Errors.Instance.IDMissing")
	}

	// check if the password hash verifier is set
	if rc.verify == nil {
		return zerrors.ThrowInternal(nil, "DOM-rhMvn5", "Errors.Internal")
	}

	// (@grvijayan) todo: review permission check

	// get session
	sessionRepo := opts.sessionRepo
	retrievedSession, err := sessionRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(sessionRepo.PrimaryKeyCondition(rc.InstanceID, rc.SessionID)),
	)
	if err := handleGetError(err, "DOM-2sF2kF", objectTypeSession); err != nil {
		return err
	}

	if retrievedSession.UserID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-EaLqwq", "Errors.User.UserIDMissing")
	}

	// get the user from the session
	userRepo := opts.userRepo
	retrievedUser, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(retrievedSession.UserID)),
		database.WithResultLock(),
	)
	if err := handleGetError(err, "DOM-Ot3qO6", objectTypeUser); err != nil {
		return err
	}

	// check user state
	if retrievedUser.State == UserStateLocked {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-47H1Ii", "Errors.User.Locked")
	}

	// check user MFA state
	if retrievedUser.Human == nil || retrievedUser.Human.RecoveryCodes == nil || len(retrievedUser.Human.RecoveryCodes.Codes) == 0 {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.User.MFA.RecoveryCodes.NotReady")
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

	hashedRecoveryCode, checkErr := validateRecoveryCode(rc.CheckRecoveryCode.RecoveryCode, rc.user.Human.RecoveryCodes.Codes, rc.verify)
	if checkErr != nil {
		err = rc.handleRecoveryCodeCheckFailed(ctx, opts)
		// logging the update error and returning the original error related to recovery code verification
		logging.OnError(ctx, err).Error("failed to update user and session after recovery code check failure")
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
			),
		)
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

	lockoutPolicy, err := GetLockoutPolicy(ctx, opts.DB(), opts.lockoutSettingRepo, rc.InstanceID, rc.user.OrganizationID)
	logging.OnError(ctx, err).Error("failed to get lockout policy")

	// update user state and recovery_code_failed_attempts
	humanRepo := opts.userRepo.Human()
	userUpdates := make([]database.Change, 0, 2)

	// update recovery_code_failed_attempts for the user
	userUpdates = append(userUpdates, humanRepo.IncrementRecoveryCodeFailedAttempts())

	// update the user's state to locked if the failed recovery code check attempts exceed the configured max value
	if lockoutPolicy != nil &&
		lockoutPolicy.MaxOTPAttempts != nil &&
		*lockoutPolicy.MaxOTPAttempts > 0 &&
		(uint64(rc.user.Human.RecoveryCodes.FailedAttempts)+1 >= *lockoutPolicy.MaxOTPAttempts) {
		userUpdates = append(userUpdates, humanRepo.SetState(UserStateLocked))
		rc.userLocked = true
	}

	err = rc.updateUser(
		ctx,
		opts,
		humanRepo,
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
	humanRepo := opts.userRepo.Human()
	err := rc.updateUser(
		ctx,
		opts,
		humanRepo,
		humanRepo.SetLastSuccessfulRecoveryCodeCheck(checkTime),
		humanRepo.RemoveRecoveryCode(hashedRecoveryCode),
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
		return zerrors.ThrowInternal(err, "DOM-fhR4N3", "session update failed") // (@grvijayan) todo: review error message format
	}
	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-Frnwxt", "Errors.Session.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-gYp8tG", "unexpected number of rows updateCount")
	}
	return nil
}

func (rc *RecoveryCodeCheckCommand) updateUser(
	ctx context.Context,
	opts *InvokeOpts,
	humanRepo HumanUserRepository,
	changes ...database.Change,
) error {
	updateCount, err := humanRepo.Update(
		ctx,
		opts.DB(),
		humanRepo.PrimaryKeyCondition(rc.InstanceID, rc.user.ID),
		changes...,
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-XGf3Tk", "user update failed")
	}
	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-hQu5ns", "Errors.User.NotFound")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-EWjTOH", "Errors.Internal")
	}
	return nil

}

func validateRecoveryCode(reqRecoveryCode string, recoveryCodes []string, verify verifierFn) (string, error) {
	if reqRecoveryCode == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "DOM-dk1MaX", "Errors.User.MFA.RecoveryCodes.InvalidCode")
	}
	var matchedCode string
	for _, recoveryCode := range recoveryCodes {
		_, err := verify(recoveryCode, reqRecoveryCode)
		if err == nil && matchedCode == "" {
			matchedCode = recoveryCode
		}
	}
	if matchedCode != "" {
		return matchedCode, nil
	}
	return "", zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode")
}
