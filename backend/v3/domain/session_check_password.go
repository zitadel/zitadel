package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

// tarpitFn represents a tarpit function
//
// The input is the number of failed attempts after which the tarpit is started
type tarpitFn func(failedAttempts uint64)

type PasswordCheckCommand struct {
	CheckPassword *session_grpc.CheckPassword

	sessionID  string
	instanceID string
	tarpitFunc tarpitFn
	verifierFn func(encoded, password string) (updated string, err error)

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
func NewPasswordCheckCommand(sessionID, instanceID string, tarpitFunc tarpitFn, verifyFn func(encoded, password string) (updated string, err error), request *session_grpc.CheckPassword) *PasswordCheckCommand {
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
		sessionID:     strings.TrimSpace(sessionID),
		instanceID:    strings.TrimSpace(instanceID),
		tarpitFunc:    tf,
		verifierFn:    verifierFunction,
	}
}

// RequiresTransaction implements [Transactional].
func (p *PasswordCheckCommand) RequiresTransaction() {}

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

	toReturn = append(toReturn, session.NewPasswordCheckedEvent(ctx, &session.NewAggregate(p.sessionID, p.instanceID).Aggregate, p.CheckTime))

	return toReturn, nil
}

// Execute implements [Commander].
func (p *PasswordCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.CheckPassword == nil {
		return nil
	}

	humanRepo := opts.userRepo.Human()
	sessionRepo := opts.sessionRepo

	updatedHash, err := p.verifierFn(p.FetchedUser.Human.Password.Password, p.CheckPassword.GetPassword())
	pswCheckType, err := p.GetPasswordCheckAndError(err)
	changes, changesErr := p.GetPasswordCheckChanges(ctx, opts, humanRepo, updatedHash, pswCheckType)
	if changesErr != nil {
		return changesErr
	}

	updateCount, updateErr := humanRepo.Update(
		ctx,
		opts.DB(),

		humanRepo.IDCondition(p.FetchedUser.ID),

		changes,
	)
	if updateErr != nil {
		return zerrors.ThrowInternal(updateErr, "DOM-netNam", "failed updating user")
	}

	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-8wVrNc", "user not found")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-D4hy9C", "unexpected number of rows updated")
	}

	var passwordFactor SessionFactor
	if err == nil {
		passwordFactor = &SessionFactorPassword{LastVerifiedAt: p.CheckTime}
	} else {
		passwordFactor = &SessionFactorPassword{LastFailedAt: p.CheckTime}
	}

	updateCount, updateErr = sessionRepo.Update(ctx, opts.DB(), sessionRepo.IDCondition(p.sessionID), sessionRepo.SetFactor(passwordFactor))
	if updateErr != nil {
		return zerrors.ThrowInternal(updateErr, "DOM-IZagay", "failed updating session")
	}

	if updateCount == 0 {
		return zerrors.ThrowNotFound(nil, "DOM-H9Q59c", "session not found")
	}
	if updateCount > 1 {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-Tbvpy8", "unexpected number of rows updated")
	}

	if err != nil && p.tarpitFunc != nil {
		p.tarpitFunc(uint64(p.FetchedUser.Human.Password.FailedAttempts + 1))
	}

	p.IsValidated = true
	p.IsValidationSuccessful = err == nil

	// TODO(IAM-Marco): The error returned here will block the transaction and stop events from being emitted.
	// This error is functional so it should be returned AND the transaction should succeed. How can we fix it?
	return err
}

func (p *PasswordCheckCommand) getLockoutPolicy(ctx context.Context, opts *InvokeOpts, orgID string) (*LockoutSettings, error) {
	lockoutSettingRepo := opts.lockoutSettingRepo

	// We need the organization lockout policy first, and if not available, the instance (default) policy.
	// So we retrieve all records with a matching instance ID and organization ID OR
	// all records with a matching instance ID and NULL (or empty) organization ID.
	// Then we assume NULLs are sorted as largest numbers (that's the case in Postgres),
	// so we sort ascending by organization ID.
	// We limit the result to 1 so that we get either the org policy or the instance one.
	settings, err := lockoutSettingRepo.List(ctx, opts.DB(),
		p.listLockoutSettingCondition(lockoutSettingRepo, orgID),
		database.WithOrderByAscending(lockoutSettingRepo.OrganizationIDColumn(), lockoutSettingRepo.InstanceIDColumn()),
		database.WithLimit(1),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOM-3B8Z6s", "failed fetching lockout settings")
	}

	if rowsReturned := len(settings); rowsReturned != 1 {
		return nil, zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, int64(rowsReturned)), "DOM-mmsrCt", "unexpected number of rows returned")
	}

	return settings[0], nil
}

func (p *PasswordCheckCommand) listLockoutSettingCondition(repo LockoutSettingsRepository, orgID string) database.QueryOption {
	instanceAndOrg := database.And(repo.InstanceIDCondition(p.instanceID), repo.OrganizationIDCondition(&orgID))
	orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
	onlyInstance := database.And(repo.InstanceIDCondition(p.instanceID), orgNullOrEmpty)

	return database.WithCondition(database.Or(instanceAndOrg, onlyInstance))
}

func (p *PasswordCheckCommand) GetPasswordCheckChanges(ctx context.Context, opts *InvokeOpts, humanRepo HumanUserRepository, updatedHash string, checkType PasswordCheckType) (database.Changes, error) {
	dbUpdates := make(database.Changes, 1, 2)
	dbUpdates[0] = humanRepo.CheckPassword(checkType)

	switch checkType.(type) {
	case *CheckTypeSucceeded:
		if updatedHash != "" {
			p.UpdatedHashedPsw = updatedHash
			dbUpdates = append(dbUpdates, humanRepo.SetPassword(&VerificationTypeSkipped{Value: &updatedHash, VerifiedAt: p.CheckTime}))
		}
	case *CheckTypeFailed:
		lockoutPolicy, err := getLockoutPolicy(ctx, opts.DB(), opts.lockoutSettingRepo, p.instanceID, p.FetchedUser.OrganizationID)
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

func (p *PasswordCheckCommand) GetPasswordCheckAndError(err error) (PasswordCheckType, error) {
	p.CheckTime = time.Now()
	if err == nil {
		return &CheckTypeSucceeded{SucceededAt: p.CheckTime}, nil
	}

	// TODO(IAM-Marco): Do we actually want to differentiate? I feel that it's giving away relevant info
	// about the password
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		err = zerrors.ThrowInvalidArgument(
			NewPasswordVerificationError(p.FetchedUser.Human.Password.FailedAttempts+1),
			"DOM-3gcfDV",
			"Errors.User.Password.Invalid",
		)
		return &CheckTypeFailed{FailedAt: p.CheckTime}, err
	}

	return &CheckTypeFailed{FailedAt: p.CheckTime}, zerrors.ThrowInternal(err, "DOM-xceNzI", "Errors.Internal")
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

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.sessionID)))
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

	if human.Password.Password == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-gklgos", "Errors.User.Password.NotSet")
	}

	p.FetchedUser = *user

	return nil
}

var _ Commander = (*PasswordCheckCommand)(nil)
var _ Transactional = (*PasswordCheckCommand)(nil)
