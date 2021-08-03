package view

import (
	"github.com/caos/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func UserLockByID(db *gorm.DB, table, userID string) (*model.UserLockView, error) {
	user := new(model.UserLockView)
	query := repository.PrepareGetByKey(table, model.UserSearchKey(usr_model.UserLockSearchKeyID), userID)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-sj8Sw", "Errors.User.NotFound")
	}
	return user, err
}

func PutUserLock(db *gorm.DB, table string, userLock *model.UserLockView) error {
	save := repository.PrepareSave(table)
	return save(db, userLock)
}

func DeleteUserLock(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.UserLockSearchKeyID), userID)
	return delete(db)
}
