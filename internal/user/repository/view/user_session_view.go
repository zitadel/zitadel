package view

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zitadel/logging"

	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed user_session_by_id.sql
var userSessionByIDQuery string

//go:embed user_sessions_by_user_agent.sql
var userSessionsByUserAgentQuery string

func UserSessionByIDs(db *gorm.DB, agentID, userID, instanceID string) (*model.UserSessionView, error) {
	tx, err := db.DB().BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			logging.WithError(err).Info("commit failed")
			err = tx.Rollback()
			logging.OnError(err).Info("rollback failed")
		}
	}()
	row := tx.QueryRow(
		userSessionByIDQuery,
		agentID,
		userID,
		instanceID,
	)
	userSession, err := scanUserSession(row)
	return userSession, err
}
func UserSessionsByAgentID(db *gorm.DB, agentID, instanceID string) ([]*model.UserSessionView, error) {
	tx, err := db.DB().BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			logging.WithError(err).Info("commit failed")
			err = tx.Rollback()
			logging.OnError(err).Info("rollback failed")
		}
	}()
	rows, err := tx.Query(
		userSessionsByUserAgentQuery,
		instanceID,
		agentID,
	)
	if err != nil {
		return nil, err
	}
	userSessions, err := scanUserSessions(rows)
	return userSessions, err
}

func PutUserSession(db *gorm.DB, table string, session *model.UserSessionView) error {
	save := repository.PrepareSave(table)
	return save(db, session)
}

func DeleteUserSessions(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.UserSessionSearchKey(usr_model.UserSessionSearchKeyUserID), Value: userID},
		repository.Key{Key: model.UserSessionSearchKey(usr_model.UserSessionSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
}

func DeleteInstanceUserSessions(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table,
		model.UserSessionSearchKey(usr_model.UserSessionSearchKeyInstanceID),
		instanceID,
	)
	return delete(db)
}

func DeleteOrgUserSessions(db *gorm.DB, table, instanceID, orgID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.UserSessionSearchKey(usr_model.UserSessionSearchKeyResourceOwner), Value: orgID},
		repository.Key{Key: model.UserSessionSearchKey(usr_model.UserSessionSearchKeyInstanceID), Value: instanceID},
	)
	return delete(db)
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
