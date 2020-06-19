package view

import (
	"github.com/jinzhu/gorm"

	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view"
)

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

func UserSessionsByUserID(db *gorm.DB, table, userID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.USERSESSIONSEARCHKEY_USER_ID,
		Method: global_model.SEARCHMETHOD_EQUALS,
		Value:  userID,
	}
	query := view.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func UserSessionsByAgentID(db *gorm.DB, table, agentID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.USERSESSIONSEARCHKEY_USER_AGENT_ID,
		Method: global_model.SEARCHMETHOD_EQUALS,
		Value:  agentID,
	}
	query := view.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func PutUserSession(db *gorm.DB, table string, session *model.UserSessionView) error {
	save := view.PrepareSave(table)
	return save(db, session)
}

func DeleteUserSessions(db *gorm.DB, table, userID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserSessionSearchKey(usr_model.USERSESSIONSEARCHKEY_USER_ID), userID)
	return delete(db)
}
