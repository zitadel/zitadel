package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func UserByID(db *gorm.DB, table, userID, instanceID string) (*model.UserView, error) {
	user := new(model.UserView)
	userIDQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	ownerRemovedQuery := &model.UserSearchQuery{
		Key:    usr_model.UserSearchOwnerRemoved,
		Method: domain.SearchMethodEquals,
		Value:  false,
	}
	query := repository.PrepareGetByQuery(table, userIDQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, user)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-sj8Sw", "Errors.User.NotFound")
	}
	user.SetEmptyUserType()
	return user, err
}

func UsersByOrgID(db *gorm.DB, table, orgID, instanceID string) ([]*model.UserView, error) {
	users := make([]*model.UserView, 0)
	orgIDQuery := &usr_model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  orgID,
	}
	instanceIDQuery := &usr_model.UserSearchQuery{
		Key:    usr_model.UserSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	ownerRemovedQuery := &usr_model.UserSearchQuery{
		Key:    usr_model.UserSearchOwnerRemoved,
		Method: domain.SearchMethodEquals,
		Value:  false,
	}
	query := repository.PrepareSearchQuery(table, model.UserSearchRequest{
		Queries: []*usr_model.UserSearchQuery{orgIDQuery, instanceIDQuery, ownerRemovedQuery},
	})
	_, err := query(db, &users)
	return users, err
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

func DeleteUser(db *gorm.DB, table, userID, instanceID string) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{model.UserSearchKey(usr_model.UserSearchKeyUserID), userID},
		repository.Key{model.UserSearchKey(usr_model.UserSearchKeyInstanceID), instanceID},
	)
	return delete(db)
}

func DeleteInstanceUsers(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserSearchKey(usr_model.UserSearchKeyInstanceID), instanceID)
	return delete(db)
}

func UpdateOrgOwnerRemovedUsers(db *gorm.DB, table, instanceID, aggID string) error {
	update := repository.PrepareUpdateByKeys(table,
		model.UserSearchKey(usr_model.UserSearchOwnerRemoved),
		true,
		repository.Key{Key: model.UserSearchKey(usr_model.UserSearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: model.UserSearchKey(usr_model.UserSearchKeyResourceOwner), Value: aggID},
	)
	return update(db)
}
