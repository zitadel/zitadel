package view

import (
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"

	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func UserSessionByIDs(db *gorm.DB, table, agentID, userID string) (*model.UserSessionView, error) {
	userSession := new(model.UserSessionView)
	userAgentQuery := model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserAgentID,
		Method: global_model.SearchMethodEquals,
		Value:  agentID,
	}
	userQuery := model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserID,
		Method: global_model.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareGetByQuery(table, userAgentQuery, userQuery)
	err := query(db, userSession)
	return userSession, err
}

func UserSessionsByUserID(db *gorm.DB, table, userID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserID,
		Method: global_model.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func UserSessionsByAgentID(db *gorm.DB, table, agentID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserAgentID,
		Method: global_model.SearchMethodEquals,
		Value:  agentID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func PutUserSession(db *gorm.DB, table string, session *model.UserSessionView) error {
	save := repository.PrepareSave(table)
	return save(db, session)
}

func DeleteUserSessions(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSessionSearchKey(usr_model.UserSessionSearchKeyUserID), userID)
	return delete(db)
}
