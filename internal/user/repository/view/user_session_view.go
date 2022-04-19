package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func UserSessionByIDs(db *gorm.DB, table, agentID, userID, instanceID string) (*model.UserSessionView, error) {
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
	instanceIDQuery := &model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareGetByQuery(table, userAgentQuery, userQuery, instanceIDQuery)
	err := query(db, userSession)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-NGBs1", "Errors.UserSession.NotFound")
	}
	return userSession, err
}

func UserSessionsByUserID(db *gorm.DB, table, userID, instanceID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery, instanceIDQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func UserSessionsByAgentID(db *gorm.DB, table, agentID, instanceID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyUserAgentID,
		Method: domain.SearchMethodEquals,
		Value:  agentID,
	}
	instanceIDQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{userAgentQuery, instanceIDQuery},
	})
	_, err := query(db, &userSessions)
	return userSessions, err
}

func ActiveUserSessions(db *gorm.DB, table string) (uint64, error) {
	activeQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyState,
		Method: domain.SearchMethodEquals,
		Value:  domain.UserSessionStateActive,
	}
	query := repository.PrepareSearchQuery(table, model.UserSessionSearchRequest{
		Queries: []*usr_model.UserSessionSearchQuery{activeQuery},
	})
	return query(db, nil)
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

func DeleteUserSessions(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.UserSessionSearchKey(usr_model.UserSessionSearchKeyUserID), userID},
		repository.Key{model.UserSessionSearchKey(usr_model.UserSessionSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}
