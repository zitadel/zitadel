package view

import (
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

func UserSessionByIDs(db *database.DB, agentID, userID, instanceID string) (userSession *model.UserSessionView, err error) {
	err = db.QueryRow(
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
func UserSessionsByAgentID(db *database.DB, agentID, instanceID string) (userSessions []*model.UserSessionView, err error) {
	err = db.Query(
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

func scanUserSession(row *sql.Row) (*model.UserSessionView, error) {
	session := new(model.UserSessionView)
	var userName, loginName, displayName, avatarKey sql.NullString
	err := row.Scan(
		&session.CreationDate,
		&session.ChangeDate,
		&session.ResourceOwner,
		&session.State,
		&session.UserAgentID,
		&session.UserID,
		&userName,
		&loginName,
		&displayName,
		&avatarKey,
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
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-NGBs1", "Errors.UserSession.NotFound")
	}
	session.UserName = userName.String
	session.LoginName = loginName.String
	session.DisplayName = displayName.String
	session.AvatarKey = avatarKey.String
	return session, err
}

func scanUserSessions(rows *sql.Rows) ([]*model.UserSessionView, error) {
	sessions := make([]*model.UserSessionView, 0)
	for rows.Next() {
		session := new(model.UserSessionView)
		var userName, loginName, displayName, avatarKey sql.NullString
		err := rows.Scan(
			&session.CreationDate,
			&session.ChangeDate,
			&session.ResourceOwner,
			&session.State,
			&session.UserAgentID,
			&session.UserID,
			&userName,
			&loginName,
			&displayName,
			&avatarKey,
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
		)
		if err != nil {
			return nil, err
		}
		session.UserName = userName.String
		session.LoginName = loginName.String
		session.DisplayName = displayName.String
		session.AvatarKey = avatarKey.String
		sessions = append(sessions, session)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "VIEW-FSF3g", "Errors.Query.CloseRows")
	}
	return sessions, nil
}
