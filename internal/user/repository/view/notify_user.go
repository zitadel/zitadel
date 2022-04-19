package view

import (
	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

func NotifyUserByID(db *gorm.DB, table, userID, instanceID string) (*model.NotifyUser, error) {
	user := new(model.NotifyUser)
	query := repository.PrepareGetByQuery(table,
		model.NotifyUserSearchQuery{Key: usr_model.NotifyUserSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID},
		model.NotifyUserSearchQuery{Key: usr_model.NotifyUserSearchKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
	)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Gad31", "Errors.User.NotFound")
	}
	return user, err
}

func NotifyUsersByOrgID(db *gorm.DB, table, orgID, instanceID string) ([]*model.NotifyUser, error) {
	users := make([]*model.NotifyUser, 0)
	orgIDQuery := &usr_model.NotifyUserSearchQuery{
		Key:    usr_model.NotifyUserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  orgID,
	}
	instanceIDQuery := &usr_model.NotifyUserSearchQuery{
		Key:    usr_model.NotifyUserSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, model.NotifyUserSearchRequest{
		Queries: []*usr_model.NotifyUserSearchQuery{orgIDQuery, instanceIDQuery},
	})
	_, err := query(db, &users)
	return users, err
}

func PutNotifyUser(db *gorm.DB, table string, project *model.NotifyUser) error {
	save := repository.PrepareSave(table)
	return save(db, project)
}

func DeleteNotifyUser(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.UserSearchKey(usr_model.NotifyUserSearchKeyUserID), userID},
		repository.Key{model.UserSearchKey(usr_model.NotifyUserSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}
