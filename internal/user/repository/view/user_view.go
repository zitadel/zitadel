package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func UserByID(db *gorm.DB, table, userID string) (*model.UserView, error) {
	user := new(model.UserView)
	query := view.PrepareGetByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_USER_ID), userID)
	err := query(db, user)
	return user, err
}

func UserByUserName(db *gorm.DB, table, userName string) (*model.UserView, error) {
	user := new(model.UserView)
	query := view.PrepareGetByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_USER_NAME), userName)
	err := query(db, user)
	return user, err
}

func SearchUsers(db *gorm.DB, table string, req *usr_model.UserSearchRequest) ([]*model.UserView, int, error) {
	users := make([]*model.UserView, 0)
	query := view.PrepareSearchQuery(table, model.UserSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func GetGlobalUserByEmail(db *gorm.DB, table, email string) (*model.UserView, error) {
	user := new(model.UserView)
	query := view.PrepareGetByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_EMAIL), email)
	err := query(db, user)
	return user, err
}

func IsUserUnique(db *gorm.DB, table, userName, email string) (bool, error) {
	user := new(model.UserView)
	query := view.PrepareGetByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_EMAIL), email)
	err := query(db, user)
	if err != nil && !caos_errs.IsNotFound(err) {
		return false, err
	}
	if user != nil {
		return false, nil
	}
	query = view.PrepareGetByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_USER_NAME), email)
	err = query(db, user)
	if err != nil && !caos_errs.IsNotFound(err) {
		return false, err
	}
	return user == nil, nil
}

func UserMfas(db *gorm.DB, table, userID string) ([]*usr_model.MultiFactor, error) {
	user, err := UserByID(db, table, userID)
	if err != nil {
		return nil, err
	}
	if user.OTPState == int32(usr_model.MFASTATE_UNSPECIFIED) {
		return []*usr_model.MultiFactor{}, nil
	}
	return []*usr_model.MultiFactor{&usr_model.MultiFactor{Type: usr_model.MFATYPE_OTP, State: usr_model.MfaState(user.OTPState)}}, nil
}

func PutUser(db *gorm.DB, table string, project *model.UserView) error {
	save := view.PrepareSave(table)
	return save(db, project)
}

func DeleteUser(db *gorm.DB, table, userID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.USERSEARCHKEY_USER_ID), userID)
	return delete(db)
}
