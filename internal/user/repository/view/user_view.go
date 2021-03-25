package view

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func UserByID(db *gorm.DB, table, userID string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.PrepareGetByKey(table, model.UserSearchKey(usr_model.UserSearchKeyUserID), userID)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-sj8Sw", "Errors.User.NotFound")
	}
	user.SetEmptyUserType()
	return user, err
}

func UserByUserName(db *gorm.DB, table, userName string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.PrepareGetByKey(table, model.UserSearchKey(usr_model.UserSearchKeyUserName), userName)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Lso9s", "Errors.User.NotFound")
	}
	user.SetEmptyUserType()
	return user, err
}

func UserByLoginName(db *gorm.DB, table, loginName string) (*model.UserView, error) {
	user := new(model.UserView)
	loginNameQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyLoginNames,
		Method: domain.SearchMethodListContains,
		Value:  loginName,
	}
	query := repository.PrepareGetByQuery(table, loginNameQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-AD4qs", "Errors.User.NotFound")
	}
	user.SetEmptyUserType()
	return user, err
}

func UserByLoginNameAndResourceOwner(db *gorm.DB, table, loginName, resourceOwner string) (*model.UserView, error) {
	user := new(model.UserView)
	loginNameQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyLoginNames,
		Method: domain.SearchMethodListContains,
		Value:  loginName,
	}
	resourceOwnerQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  resourceOwner,
	}
	query := repository.PrepareGetByQuery(table, loginNameQuery, resourceOwnerQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-AD4qs", "Errors.User.NotFoundOnOrg")
	}
	user.SetEmptyUserType()
	return user, err
}

func UsersByOrgID(db *gorm.DB, table, orgID string) ([]*model.UserView, error) {
	users := make([]*model.UserView, 0)
	orgIDQuery := &usr_model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  orgID,
	}
	query := repository.PrepareSearchQuery(table, model.UserSearchRequest{
		Queries: []*usr_model.UserSearchQuery{orgIDQuery},
	})
	_, err := query(db, &users)
	return users, err
}

func UserIDsByDomain(db *gorm.DB, table, orgDomain string) ([]string, error) {
	type id struct {
		Id string
	}
	ids := make([]id, 0)
	orgIDQuery := &usr_model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyUserName,
		Method: domain.SearchMethodEndsWithIgnoreCase,
		Value:  "%" + orgDomain,
	}
	query := repository.PrepareSearchQuery(table, model.UserSearchRequest{
		Queries: []*usr_model.UserSearchQuery{orgIDQuery},
	})
	_, err := query(db, &ids)
	if err != nil {
		return nil, err
	}
	users := make([]string, len(ids))
	for i, id := range ids {
		users[i] = id.Id
	}
	return users, err
}

func SearchUsers(db *gorm.DB, table string, req *usr_model.UserSearchRequest) ([]*model.UserView, uint64, error) {
	users := make([]*model.UserView, 0)
	query := repository.PrepareSearchQuery(table, model.UserSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func GetGlobalUserByLoginName(db *gorm.DB, table, loginName string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.PrepareGetByQuery(table, &model.UserSearchQuery{Key: usr_model.UserSearchKeyLoginNames, Value: loginName, Method: domain.SearchMethodListContains})
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-8uWer", "Errors.User.NotFound")
	}
	user.SetEmptyUserType()
	return user, err
}

func IsUserUnique(db *gorm.DB, table, userName, email string) (bool, error) {
	user := new(model.UserView)
	query := repository.PrepareGetByKey(table, model.UserSearchKey(usr_model.UserSearchKeyEmail), email)
	err := query(db, user)
	if err != nil && !caos_errs.IsNotFound(err) {
		return false, err
	}
	if user.UserName != "" {
		return false, nil
	}
	query = repository.PrepareGetByKey(table, model.UserSearchKey(usr_model.UserSearchKeyUserName), userName)
	err = query(db, user)
	if err != nil && !caos_errs.IsNotFound(err) {
		return false, err
	}
	return user.UserName == "", nil
}

func UserMFAs(db *gorm.DB, table, userID string) ([]*usr_model.MultiFactor, error) {
	user, err := UserByID(db, table, userID)
	if err != nil {
		return nil, err
	}
	if user.OTPState == int32(usr_model.MFAStateUnspecified) {
		return []*usr_model.MultiFactor{}, nil
	}
	return []*usr_model.MultiFactor{{Type: usr_model.MFATypeOTP, State: usr_model.MFAState(user.OTPState)}}, nil
}

func PutUsers(db *gorm.DB, table string, users ...*model.UserView) error {
	save := repository.PrepareBulkSave(table)
	u := make([]interface{}, len(users))
	for i, user := range users {
		u[i] = user
	}
	return save(db, u...)
}

func PutUser(db *gorm.DB, table string, user *model.UserView) error {
	save := repository.PrepareSave(table)
	return save(db, user)
}

func DeleteUser(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.UserSearchKeyUserID), userID)
	return delete(db)
}
