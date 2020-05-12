package view

import (
	"github.com/jinzhu/gorm"

	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view"
)

func UserSessionByID(db *gorm.DB, table, sessionID string) (*model.UserSessionView, error) {
	userSession := new(model.UserSessionView)
	query := view.PrepareGetByKey(table, model.UserSessionSearchKey(usr_model.USERSESSIONSEARCHKEY_SESSION_ID), sessionID)
	err := query(db, userSession)
	return userSession, err
}

func UserSessionByIDs(db *gorm.DB, table, agentID, userID string) (*model.UserSessionView, error) {
	userSession := new(model.UserSessionView)
	userAgentQuery := model.UserSessionSearchQuery{
		Key:    usr_model.USERSESSIONSEARCHKEY_USER_AGENT_ID,
		Method: global_model.SEARCHMETHOD_EQUALS,
		Value:  agentID,
	}
	userQuery := model.UserSessionSearchQuery{
		Key:    usr_model.USERSESSIONSEARCHKEY_USER_ID,
		Method: global_model.SEARCHMETHOD_EQUALS,
		Value:  userID,
	}
	query := view.PrepareGetByQuery(table, userAgentQuery, userQuery)
	err := query(db, userSession)
	return userSession, err
}

func PutUserSession(db *gorm.DB, table string, session *model.UserSessionView) error {
	save := view.PrepareSave(table)
	return save(db, session)
}

func DeleteUserSession(db *gorm.DB, table, sessionID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserSessionSearchKey(usr_model.USERSESSIONSEARCHKEY_USER_ID), sessionID)
	return delete(db)
}
