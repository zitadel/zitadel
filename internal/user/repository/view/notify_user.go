package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func NotifyUserByID(db *gorm.DB, table, userID string) (*model.NotifyUser, error) {
	user := new(model.NotifyUser)
	query := repository.PrepareGetByKey(table, model.NotifyUserSearchKey(usr_model.NotifyUserSearchKeyUserID), userID)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Gad31", "Errors.User.NotFound")
	}
	return user, err
}

func NotifyUsersByOrgID(db *gorm.DB, table, orgID string) ([]*model.NotifyUser, error) {
	users := make([]*model.NotifyUser, 0)
	orgIDQuery := &usr_model.NotifyUserSearchQuery{
		Key:    usr_model.NotifyUserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  orgID,
	}
	query := repository.PrepareSearchQuery(table, model.NotifyUserSearchRequest{
		Queries: []*usr_model.NotifyUserSearchQuery{orgIDQuery},
	})
	_, err := query(db, &users)
	return users, err
}

func PutNotifyUser(db *gorm.DB, table string, project *model.NotifyUser) error {
	save := repository.PrepareSave(table)
	return save(db, project)
}

func DeleteNotifyUser(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.NotifyUserSearchKeyUserID), userID)
	return delete(db)
}
