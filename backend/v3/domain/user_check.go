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
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type UserCheckCommand struct {
	CheckUser  *session_grpc.CheckUser
	SessionID  string
	InstanceID string

	FetchedUser User

	// Out
	UserCheckedAt         time.Time
	PreferredUserLanguage *language.Tag
}

// NewUserCheckCommand returns a check Commander validating the input user.
//
// It assumes that a [Session] already exists: this check should be part of the
// batch call to create/set a session.
//
// The check will update the existing session or return an error if the session
// is not found or validation fails.
func NewUserCheckCommand(sessionID, instanceID string) *UserCheckCommand {
	return &UserCheckCommand{
		SessionID:  sessionID,
		InstanceID: instanceID,
	}
}

// RequiresTransaction implements [Transactional].
func (u *UserCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (u *UserCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if u.CheckUser == nil {
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
	if u.CheckUser == nil {
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
	if u.CheckUser == nil {
		return
	}
	var usrQueryOpt database.QueryOption
	userRepo := opts.userRepo

	switch searchType := u.CheckUser.GetSearch().(type) {
	case *session_grpc.CheckUser_UserId:
		usrQueryOpt = database.WithCondition(userRepo.IDCondition(searchType.UserId))
	case *session_grpc.CheckUser_LoginName:
		usrQueryOpt = database.WithCondition(userRepo.LoginNameCondition(database.TextOperationEqual, searchType.LoginName))
	default:
		return zerrors.ThrowInvalidArgumentf(nil, "DOM-7B2m0b", "user search %T not implemented", searchType)
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
