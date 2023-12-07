package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-NGBs1", "Errors.UserSession.NotFound")
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

func UserSessionsByOrgID(db *gorm.DB, table, orgID, instanceID string) ([]*model.UserSessionView, error) {
	userSessions := make([]*model.UserSessionView, 0)
	userAgentQuery := &usr_model.UserSessionSearchQuery{
		Key:    usr_model.UserSessionSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  orgID,
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
