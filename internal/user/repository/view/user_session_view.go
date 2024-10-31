package view

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed user_session_by_id.sql
var userSessionByIDQuery string

//go:embed user_sessions_by_user_agent.sql
var userSessionsByUserAgentQuery string

//go:embed user_agent_by_user_session_id.sql
var userAgentByUserSessionIDQuery string

//go:embed active_user_sessions_by_session_id.sql
var activeUserSessionsBySessionIDQuery string

func UserSessionByIDs(ctx context.Context, db *database.DB, agentID, userID, instanceID string) (userSession *model.UserSessionView, err error) {
	err = db.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			userSession, err = scanUserSession(row)
			return err
		},
		userSessionByIDQuery,
		agentID,
		userID,
		instanceID,
	)
	return userSession, err
}

func UserSessionsByAgentID(ctx context.Context, db *database.DB, agentID, instanceID string) (userSessions []*model.UserSessionView, err error) {
	err = db.QueryContext(
		ctx,
		func(rows *sql.Rows) error {
			userSessions, err = scanUserSessions(rows)
			return err
		},
		userSessionsByUserAgentQuery,
		agentID,
		instanceID,
	)
	return userSessions, err
}

func UserAgentIDBySessionID(ctx context.Context, db *database.DB, sessionID, instanceID string) (userAgentID string, err error) {
	err = db.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(&userAgentID)
		},
		userAgentByUserSessionIDQuery,
		sessionID,
		instanceID,
	)
	return userAgentID, err
}

// ActiveUserSessionsBySessionID returns all sessions (sessionID:userID map) with an active session on the same user agent (its id is also returned) based on a sessionID
func ActiveUserSessionsBySessionID(ctx context.Context, db *database.DB, sessionID, instanceID string) (userAgentID string, sessions map[string]string, err error) {
	err = db.QueryContext(
		ctx,
		func(rows *sql.Rows) error {
			userAgentID, sessions, err = scanActiveUserAgentUserIDs(rows)
			return err
		},
		activeUserSessionsBySessionIDQuery,
		sessionID,
		instanceID,
	)
	return userAgentID, sessions, err
}

func scanActiveUserAgentUserIDs(rows *sql.Rows) (userAgentID string, sessions map[string]string, err error) {
	sessions = make(map[string]string)
	for rows.Next() {
		var userID, sessionID string
		err := rows.Scan(
			&userAgentID,
			&userID,
			&sessionID,
		)
		if err != nil {
			return "", nil, err
		}
		sessions[sessionID] = userID
	}
	if err := rows.Close(); err != nil {
		return "", nil, zerrors.ThrowInternal(err, "VIEW-Sbrws", "Errors.Query.CloseRows")
	}
	return userAgentID, sessions, nil
}

func scanUserSession(row *sql.Row) (*model.UserSessionView, error) {
	session := new(model.UserSessionView)
	err := row.Scan(
		&session.CreationDate,
		&session.ChangeDate,
		&session.ResourceOwner,
		&session.State,
		&session.UserAgentID,
		&session.UserID,
		&session.UserName,
		&session.LoginName,
		&session.DisplayName,
		&session.AvatarKey,
		&session.SelectedIDPConfigID,
		&session.PasswordVerification,
		&session.PasswordlessVerification,
		&session.ExternalLoginVerification,
		&session.SecondFactorVerification,
		&session.SecondFactorVerificationType,
		&session.MultiFactorVerification,
		&session.MultiFactorVerificationType,
		&session.Sequence,
		&session.InstanceID,
		&session.ID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-NGBs1", "Errors.UserSession.NotFound")
	}
	return session, err
}

func scanUserSessions(rows *sql.Rows) ([]*model.UserSessionView, error) {
	sessions := make([]*model.UserSessionView, 0)
	for rows.Next() {
		session := new(model.UserSessionView)
		err := rows.Scan(
			&session.CreationDate,
			&session.ChangeDate,
			&session.ResourceOwner,
			&session.State,
			&session.UserAgentID,
			&session.UserID,
			&session.UserName,
			&session.LoginName,
			&session.DisplayName,
			&session.AvatarKey,
			&session.SelectedIDPConfigID,
			&session.PasswordVerification,
			&session.PasswordlessVerification,
			&session.ExternalLoginVerification,
			&session.SecondFactorVerification,
			&session.SecondFactorVerificationType,
			&session.MultiFactorVerification,
			&session.MultiFactorVerificationType,
			&session.Sequence,
			&session.InstanceID,
			&session.ID,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "VIEW-FSF3g", "Errors.Query.CloseRows")
	}
	return sessions, nil
}
