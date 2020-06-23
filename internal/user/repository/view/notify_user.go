package view

import (
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func NotifyUserByID(db *gorm.DB, table, userID string) (*model.NotifyUser, error) {
	user := new(model.NotifyUser)
	query := view.PrepareGetByKey(table, model.UserSearchKey(usr_model.NotifyUserSearchKeyUserID), userID)
	err := query(db, user)
	return user, err
}

func PutNotifyUser(db *gorm.DB, table string, project *model.NotifyUser) error {
	save := view.PrepareSave(table)
	return save(db, project)
}

func DeleteNotifyUser(db *gorm.DB, table, userID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.NotifyUserSearchKeyUserID), userID)
	return delete(db)
}
