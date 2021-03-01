package view

import (
	auth_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"

	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func UserSessionByIDs(db *gorm.DB, table, agentID, userID string) (*model.UserSessionView, error) {
	userSession := new(model.UserSessionView)
	userAgentQuery := model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserAgentID,
		Method: domain.SearchMethodEquals,
		Value:  agentID,
	}
	userQuery := model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	query := repository.PrepareGetByQuery(table, userAgentQuery, userQuery)
	err := query(db, userSession)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-NGBs1", "Errors.UserSession.NotFound")
	}
	return userSession, err
}

func UserSessionsByUserID(db *gorm.DB, table, userID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserID,
		Method: domain.SearchMethodEquals,
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
		Method: domain.SearchMethodEquals,
		Value:  agentID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func ActiveUserSessions(db *gorm.DB, table string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	activeQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyState,
		Method: domain.SearchMethodEquals,
		Value:  auth_model.UserSessionStateActive,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{activeQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func PutUserSession(db *gorm.DB, table string, session *model.UserSessionView) error {
	save := repository.PrepareSave(table)
	return save(db, session)
}

func PutUserSessions(db *gorm.DB, table string, sessions ...*model.UserSessionView) error {
	save := repository.PrepareBulkSave(table)
	s := make([]interface{}, len(sessions))
	for i, session := range sessions {
		s[i] = session
	}
	return save(db, s...)
}

func DeleteUserSessions(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSessionSearchKey(usr_model.UserSessionSearchKeyUserID), userID)
	return delete(db)
}
