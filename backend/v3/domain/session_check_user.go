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
			return zerrors.ThrowNotFound(err, "DOM-rbdCv3", "session not found")
		}
		return zerrors.ThrowInternal(err, "DOM-To1rLz", "failed fetching session")
	}

	if session.UserID != "" && u.FetchedUser.ID != "" && session.UserID != u.FetchedUser.ID {
		return zerrors.ThrowInvalidArgument(nil, "DOM-78g1TV", "user change not possible")
	}
	userFactor := &SessionFactorUser{
		UserID:         u.FetchedUser.ID,
		LastVerifiedAt: time.Now(),
	}

	updateCount, err := sessionRepo.Update(ctx, opts.DB(), sessionRepo.IDCondition(session.ID), sessionRepo.SetFactor(userFactor))
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-netNam", "failed updating session")
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-FszyWS", "session not found")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-SsIwDt", "unexpected number of rows updated")
		return err
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
		return zerrors.ThrowPreconditionFailed(nil, "DOM-00o0ys", "Errors.Missing.SessionID")
	}
	if u.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-Oe1dtz", "Errors.Missing.InstanceID")
	}

	if authZErr := opts.Permissions.CheckSessionPermission(ctx, SessionWritePermission, u.SessionID); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-4qz3mt", "Errors.PermissionDenied")
	}

	var usrQueryOpt database.QueryOption
	userRepo := opts.userRepo

	if loginName := u.InputCheckUser.LoginName; loginName != "" {
		usrQueryOpt = database.WithCondition(userRepo.LoginNameCondition(database.TextOperationEqual, loginName))
	} else if userID := u.InputCheckUser.UserID; userID != "" {
		usrQueryOpt = database.WithCondition(userRepo.IDCondition(userID))
	} else {
		return zerrors.ThrowInvalidArgument(nil, "DOM-7B2m0b", "no valid query option")
	}

	usr, err := userRepo.Get(ctx, opts.DB(), usrQueryOpt)
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-lcZeXI", "user not found")
		}
		return zerrors.ThrowInternal(err, "DOM-Y846I0", "failed fetching user")
	}

	if usr.State != UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-vgDIu9", "Errors.User.NotActive")
	}

	u.FetchedUser = *usr

	return nil
}

var _ Commander = (*UserCheckCommand)(nil)
var _ Transactional = (*UserCheckCommand)(nil)
