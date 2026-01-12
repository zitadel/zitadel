package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Querier[*Session] = (*GetSessionQuery)(nil)

type sessionTokenVerifierFunc func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error)

type GetSessionQuery struct {
	sessionID            string
	instanceID           string
	sessionToken         string
	sessionTokenVerifier sessionTokenVerifierFunc

	session *Session
}

func NewGetSessionQuery(sessionID, instanceID, sessionToken string, sessionTokenVerifier sessionTokenVerifierFunc) *GetSessionQuery {
	if sessionTokenVerifier == nil {
		sessionTokenVerifier = sessionTokenVerifierFn
	}
	return &GetSessionQuery{
		sessionID:            sessionID,
		instanceID:           instanceID,
		sessionToken:         sessionToken,
		sessionTokenVerifier: sessionTokenVerifier,
	}
}

// Validate implements [Querier].
func (g *GetSessionQuery) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if g.instanceID = strings.TrimSpace(g.instanceID); g.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "QUERY-3n9fGv", "Errors.Missing.InstanceID")
	}

	if g.sessionID = strings.TrimSpace(g.sessionID); g.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "QUERY-CtWgrV", "Errors.Missing.SessionID")
	}
	return nil
}

// Execute implements [Querier].
func (g *GetSessionQuery) Execute(ctx context.Context, opts *InvokeOpts) error {
	sessionRepo := opts.sessionRepo

	session, err := sessionRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			database.And(
				sessionRepo.InstanceIDCondition(g.instanceID),
				sessionRepo.IDCondition(g.sessionID),
			),
		),
	)
	if err = handleGetError(err, "DOM-QiiiFY", objectTypeSession); err != nil {
		return err
	}

	if g.sessionToken == "" {
		if err = checkSessionPermission(ctx, session, opts); err != nil {
			return err
		}
		g.session = session
		return nil
	}
	// verify session token
	if err = g.sessionTokenVerifier(ctx, g.sessionToken, g.sessionID, session.Token); err != nil {
		return zerrors.ThrowPermissionDenied(err, "QUERY-M3f4fS", "Errors.PermissionDenied")
	}
	g.session = session
	return nil
}

// String implements [Querier].
func (g *GetSessionQuery) String() string {
	return "GetSessionQuery"
}

// Result implements [Querier].
func (g *GetSessionQuery) Result() *Session {
	return g.session
}

func checkSessionPermission(ctx context.Context, session *Session, opts *InvokeOpts) error {
	data := authz.GetCtxData(ctx)
	// (@grvijayan) todo: 1. no permission check needed if the requestor created the session
	// creator is not yet available in the new session model

	// 2. no permission check needed if the requestor is retrieving their own session
	if session.UserID != "" &&
		data.UserID == session.UserID {
		return nil
	}

	// 3. no permission check needed for the same useragent
	if session.UserAgent != nil &&
		session.UserAgent.FingerprintID != nil &&
		*session.UserAgent.FingerprintID != "" &&
		data.AgentID != "" &&
		*session.UserAgent.FingerprintID == data.AgentID {
		return nil
	}

	// 4. check if the requestor has permissions on the user's organization
	// get user's organization
	userRepo := opts.userRepo
	user, err := userRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			userRepo.IDCondition(session.UserID),
		),
	)
	if err = handleGetError(err, "QUERY-4n9fGv", objectTypeUser); err != nil {
		return err
	}

	// check permission on the organization
	if user.OrganizationID != "" {
		if authZErr := opts.Permissions.CheckOrganizationPermission(ctx, SessionReadPermission, user.OrganizationID); authZErr != nil {
			return zerrors.ThrowPermissionDenied(authZErr, "QUERY-6l8oWp", "Errors.PermissionDenied")
		}
		return nil
	}

	// 5. check permission on the instance by default
	if authZErr := opts.Permissions.CheckInstancePermission(ctx, SessionReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "QUERY-5RJSUU", "Errors.PermissionDenied")
	}

	return nil
}
