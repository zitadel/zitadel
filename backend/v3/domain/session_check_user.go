package domain

import (
	"context"
	"errors"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CheckUserType struct {
	UserID    string
	LoginName string
}

type UserCheckCommand struct {
	InputCheckUser *CheckUserType
	SessionID      string
	InstanceID     string

	FetchedUser User

	// Out
	PreferredUserLanguage *language.Tag
	UserCheckedAt         time.Time
}

// NewUserCheckCommand returns a check Commander validating the input user.
//
// It assumes that a [Session] already exists: this check should be part of the
// batch call to create/set a session.
//
// The check will update the existing session or return an error if the session
// is not found or validation fails.
func NewUserCheckCommand(sessionID, instanceID string, checkUser *CheckUserType) *UserCheckCommand {
	return &UserCheckCommand{
		SessionID:      sessionID,
		InstanceID:     instanceID,
		InputCheckUser: checkUser,
	}
}

// RequiresTransaction implements [Transactional].
func (u *UserCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (u *UserCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if u.InputCheckUser == nil {
		return nil, nil
	}

	return []eventstore.Command{
		session.NewUserCheckedEvent(
			ctx,
			&session.NewAggregate(u.SessionID, u.InstanceID).Aggregate,
			u.FetchedUser.ID,
			u.FetchedUser.OrganizationID,
			u.UserCheckedAt,
			u.PreferredUserLanguage,
		),
	}, nil
}

// Execute implements [Commander].
func (u *UserCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if u.InputCheckUser == nil {
		return err
	}
	sessionRepo := opts.sessionRepo

	if human := u.FetchedUser.Human; human != nil {
		if !human.PreferredLanguage.IsRoot() {
			u.PreferredUserLanguage = &human.PreferredLanguage
		}
	}

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(u.SessionID)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return ErrSessionNotFound(err, u.SessionID) // TODO: not found vs precondition failed
		}
		return ErrInternal(err, "failed to fetch the session")
	}

	if session.UserID != "" && u.FetchedUser.ID != "" && session.UserID != u.FetchedUser.ID {
		return ErrSessionUserChange()
	}
	userFactor := &SessionFactorUser{
		UserID:         u.FetchedUser.ID,
		LastVerifiedAt: time.Now(),
	}

	updateCount, err := sessionRepo.Update(ctx, opts.DB(), sessionRepo.IDCondition(session.ID), sessionRepo.SetFactor(userFactor))
	if err != nil {
		return ErrInternal(err, "failed to update the session")
	}

	if updateCount == 0 {
		return ErrSessionNotFound(err, u.SessionID) // TODO: not found vs precondition failed vs internal
	}
	if updateCount > 1 {
		return ErrMoreThanOneRowAffected("unexpected number of rows updated", updateCount)
	}

	u.UserCheckedAt = userFactor.LastVerifiedAt

	return err
}

// String implements [Commander].
func (u *UserCheckCommand) String() string {
	return "UserCheckCommand"
}

// Validate implements [Commander].
func (u *UserCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if u.InputCheckUser == nil {
		return nil
	}

	if u.SessionID == "" {
		return ErrIDMissing() // TODO: precondition failed vs invalid argument
	}
	if u.InstanceID == "" {
		return ErrInstanceIDMissing() // TODO: precondition failed vs invalid argument
	}

	if authZErr := opts.Permissions.CheckSessionPermission(ctx, SessionWritePermission, u.SessionID); authZErr != nil {
		// TODO: return a more specific error once permissions are implemented?
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-4qz3mt", "Errors.PermissionDenied")
	}

	var usrQueryOpt database.QueryOption
	var query string
	userRepo := opts.userRepo

	if loginName := u.InputCheckUser.LoginName; loginName != "" {
		usrQueryOpt = database.WithCondition(userRepo.LoginNameCondition(database.TextOperationEqual, loginName))
		query = loginName
	} else if userID := u.InputCheckUser.UserID; userID != "" {
		usrQueryOpt = database.WithCondition(userRepo.IDCondition(userID))
		query = userID
	} else {
		return ErrInvalidRequest("either login_name or user_id is required")
	}

	usr, err := userRepo.Get(ctx, opts.DB(), usrQueryOpt)
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return ErrUserNotFound(err, query)
		}
		return ErrInternal(err, "failed to fetch the user")
	}

	if usr.State != UserStateActive {
		return ErrUserNotActive(usr.ID, "user must be active for sign in")
	}

	u.FetchedUser = *usr

	return nil
}

var _ Commander = (*UserCheckCommand)(nil)
var _ Transactional = (*UserCheckCommand)(nil)
