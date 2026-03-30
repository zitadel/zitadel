package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// tarpitFn represents a tarpit function
//
// The input is the number of failed attempts after which the tarpit is started
type tarpitFn func(failedAttempts uint64)

type CheckPasswordType struct {
	Password string
}

type PasswordCheckCommand struct {
	CheckPassword *CheckPasswordType

	SessionID  string
	InstanceID string
	TarpitFunc tarpitFn
	VerifierFn verifierFn

	FetchedUser      User
	UpdatedHashedPsw string
	CheckTime        time.Time

	// IsValidated indicates if the password check has completed (no transaction errors)
	IsValidated bool

	// IsValidationSuccessful indicates if the completed password check is successful
	// (i.e. no functional error is returned)
	IsValidationSuccessful bool

	// IsUserLocked indicates if the user has been locked (i.e. lockout policy check failed)
	IsUserLocked bool
}

// NewPasswordCheckCommand initializes a new [PasswordCheckCommand]
//
// If tarpitFunc is nil, the default tarpit will be used.
//
// verifyFn is a function that takes as input a target encoded password
// and an input password to verify. It returns an updated hash and an error.
// It defaults to [passwap.Swapper.Verify]
//
// The command does not implement [Transactional] due verifyFn that might take a long time to execute.
// So the DB transaction will be started only after verifyFn has been run.
//
// Moreover, the command may return a functional error so manual management of the transaction is needed
// to avoid rollbacking a transaction despite having only a functional error.
func NewPasswordCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, verifyFn func(encoded, password string) (updated string, err error), request *CheckPasswordType) *PasswordCheckCommand {
	tf := sysConfig.Tarpit.Tarpit()
	if tarpitFunc != nil {
		tf = tarpitFunc
	}

	verifierFunction := passwordHasher.Verify
	if verifyFn != nil {
		verifierFunction = verifyFn
	}

	return &PasswordCheckCommand{
		CheckPassword: request,
		SessionID:     strings.TrimSpace(sessionID),
		InstanceID:    strings.TrimSpace(instanceID),
		TarpitFunc:    tf,
		VerifierFn:    verifierFunction,
	}
}

// Events implements [Commander].
func (p *PasswordCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if p.CheckPassword == nil || !p.IsValidated {
		return nil, nil
	}

	toReturn := make([]eventstore.Command, 1, 3)
	userAgg := &user.NewAggregate(p.FetchedUser.ID, p.FetchedUser.OrganizationID).Aggregate

	if p.IsValidationSuccessful {
		toReturn[0] = user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil)
		if p.UpdatedHashedPsw != "" {
			toReturn = append(toReturn, user.NewHumanPasswordHashUpdatedEvent(ctx, userAgg, p.UpdatedHashedPsw))
		}
	} else {
		toReturn[0] = user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil)
		if p.IsUserLocked {
			toReturn = append(toReturn, user.NewUserLockedEvent(ctx, userAgg))
		}
	}

	toReturn = append(toReturn, session.NewPasswordCheckedEvent(ctx, &session.NewAggregate(p.SessionID, p.InstanceID).Aggregate, p.CheckTime))

	return toReturn, nil
}

// Execute implements [Commander].
func (p *PasswordCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.CheckPassword == nil {
		return nil
	}

	humanRepo := opts.userRepo.Human()
	sessionRepo := opts.sessionRepo

	updatedHash, verifyErr := p.VerifierFn(p.FetchedUser.Human.Password.Hash, p.CheckPassword.Password)
	pswCheckType, err := p.GetPasswordCheckAndError(verifyErr)
	changes, changesErr := p.GetPasswordCheckChanges(ctx, opts, humanRepo, updatedHash, pswCheckType)
	if changesErr != nil {
		return changesErr
	}

	tx, txErr := opts.StartTransaction(ctx, nil)
	if txErr != nil {
		return zerrors.ThrowInternal(txErr, "DOM-IR1vH2", "failed starting transaction")
	}

	defer func() {
		if endErr := tx.End(ctx, txErr); endErr != nil {
			err = endErr
		}
	}()

	updateCount, updateErr := humanRepo.Update(
		ctx,
		tx,

		humanRepo.IDCondition(p.FetchedUser.ID),

		changes,
	)
	if updateErr != nil {
		txErr = zerrors.ThrowInternal(updateErr, "DOM-netNam", "failed updating user")
		return txErr
	}

	if updateCount == 0 {
		txErr = zerrors.ThrowNotFound(nil, "DOM-8wVrNc", "user not found")
		return txErr
	}
	if updateCount > 1 {
		txErr = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-D4hy9C", "unexpected number of rows updated")
		return txErr
	}

	var passwordFactor SessionFactor
	if err == nil {
		passwordFactor = &SessionFactorPassword{LastVerifiedAt: p.CheckTime}
	} else {
		passwordFactor = &SessionFactorPassword{LastFailedAt: p.CheckTime}
	}

	updateCount, updateErr = sessionRepo.Update(ctx, tx, sessionRepo.IDCondition(p.SessionID), sessionRepo.SetFactor(passwordFactor))
	if updateErr != nil {
		txErr = zerrors.ThrowInternal(updateErr, "DOM-IZagay", "failed updating session")
		return txErr
	}

	if updateCount == 0 {
		txErr = zerrors.ThrowNotFound(nil, "DOM-H9Q59c", "session not found")
		return txErr
	}
	if updateCount > 1 {
		txErr = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-Tbvpy8", "unexpected number of rows updated")
		return txErr
	}

	if err != nil && p.TarpitFunc != nil {
		p.TarpitFunc(uint64(p.FetchedUser.Human.Password.FailedAttempts + 1))
	}

	p.IsValidated = true
	p.IsValidationSuccessful = err == nil

	// This error is functional so it should be returned AND the transaction should succeed.
	// Hence, do not use [Transactional] and make sure txErr = nil
	return err
}

func (p *PasswordCheckCommand) GetPasswordCheckChanges(ctx context.Context, opts *InvokeOpts, humanRepo HumanUserRepository, updatedHash string, checkType VerificationType) (database.Changes, error) {
	dbUpdates := make(database.Changes, 1)
	switch ct := checkType.(type) {
	case *VerificationTypeSucceeded:
		dbUpdates[0] = humanRepo.SetLastSuccessfulPasswordCheck(ct.VerifiedAt)
		if updatedHash != "" {
			p.UpdatedHashedPsw = updatedHash
			dbUpdates = append(dbUpdates, humanRepo.SetPassword(updatedHash))
		}
	case *VerificationTypeFailed:
		dbUpdates[0] = humanRepo.IncrementPasswordFailedAttempts()
		lockoutPolicy, err := GetLockoutPolicy(ctx, opts, p.InstanceID, p.FetchedUser.OrganizationID)
		if err != nil {
			return nil, err
		}

		if lockoutPolicy != nil &&
			lockoutPolicy.MaxPasswordAttempts != nil && *lockoutPolicy.MaxPasswordAttempts > 0 &&
			uint64(p.FetchedUser.Human.Password.FailedAttempts+1) >= *lockoutPolicy.MaxPasswordAttempts {

			dbUpdates = append(dbUpdates, humanRepo.SetState(UserStateLocked))
			p.IsUserLocked = true
		}
	}

	return dbUpdates, nil
}

func (p *PasswordCheckCommand) GetPasswordCheckAndError(err error) (VerificationType, error) {
	p.CheckTime = time.Now()
	if err == nil {
		return &VerificationTypeSucceeded{VerifiedAt: p.CheckTime}, nil
	}

	// TODO(IAM-Marco): Do we actually want to differentiate? I feel that it's giving away relevant info
	// about the password
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		err = zerrors.ThrowInvalidArgument(
			NewPasswordVerificationError(p.FetchedUser.Human.Password.FailedAttempts+1),
			"DOM-3gcfDV",
			"Errors.User.Password.Invalid",
		)
		return &VerificationTypeFailed{FailedAt: p.CheckTime}, err
	}

	return &VerificationTypeFailed{FailedAt: p.CheckTime}, zerrors.ThrowInternal(err, "DOM-xceNzI", "Errors.Internal")
}

// String implements [Commander].
func (p *PasswordCheckCommand) String() string {
	return "PasswordCheckCommand"
}

// Validate implements [Commander].
func (p *PasswordCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.CheckPassword == nil {
		return nil
	}

	if p.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-cRKWNx", "Errors.Missing.SessionID")
	}
	if p.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-JnEtcJ", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.SessionID)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-0XRmp8", "session not found")
		}
		return zerrors.ThrowInternal(err, "DOM-qAoQrg", "failed fetching session")
	}

	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing")
	}

	user, err := userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.IDCondition(session.UserID)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-zxKosn", "Errors.User.NotFound")
		}
		return zerrors.ThrowInternal(err, "DOM-nKD4Gq", "failed fetching user")
	}
	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-ADhxAx", "user not human")
	}
	human := user.Human

	if user.State == UserStateLocked {
		return zerrors.ThrowPreconditionFailedf(
			NewPasswordVerificationError(user.Human.Password.FailedAttempts),
			"DOM-D804Sj",
			"Errors.User.Locked",
		)
	}

	if human.Password.Hash == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-gklgos", "Errors.User.Password.NotSet")
	}

	p.FetchedUser = *user

	return nil
}

var _ Commander = (*PasswordCheckCommand)(nil)
