package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

func UserMembershipByIDs(db *gorm.DB, table, userID, aggregateID, objectID string, membertype usr_model.MemberType) (*model.UserMembershipView, error) {
	memberships := new(model.UserMembershipView)
	userIDQuery := &model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	aggregateIDQuery := &model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	objectIDQuery := &model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyObjectID, Value: objectID, Method: global_model.SearchMethodEquals}
	memberTypeQuery := &model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyMemberType, Value: int32(membertype), Method: global_model.SearchMethodEquals}

	query := repository.PrepareGetByQuery(table, userIDQuery, aggregateIDQuery, objectIDQuery, memberTypeQuery)
	err := query(db, memberships)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-sj8Sw", "Errors.UserMembership.NotFound")
	}
	return memberships, err
}

func UserMembershipsByAggregateID(db *gorm.DB, table, aggregateID string) ([]*model.UserMembershipView, error) {
	memberships := make([]*model.UserMembershipView, 0)
	aggregateIDQuery := &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareSearchQuery(table, model.UserMembershipSearchRequest{
		Queries: []*usr_model.UserMembershipSearchQuery{aggregateIDQuery},
	})
	_, err := query(db, &memberships)
	return memberships, err
}

func UserMembershipsByResourceOwner(db *gorm.DB, table, resourceOwner string) ([]*model.UserMembershipView, error) {
	memberships := make([]*model.UserMembershipView, 0)
	aggregateIDQuery := &usr_model.UserMembershipSearchQuery{Key: usr_model.UserMembershipSearchKeyResourceOwner, Value: resourceOwner, Method: global_model.SearchMethodEquals}
	query := repository.PrepareSearchQuery(table, model.UserMembershipSearchRequest{
		Queries: []*usr_model.UserMembershipSearchQuery{aggregateIDQuery},
	})
	_, err := query(db, &memberships)
	return memberships, err
}

func SearchUserMemberships(db *gorm.DB, table string, req *usr_model.UserMembershipSearchRequest) ([]*model.UserMembershipView, uint64, error) {
	users := make([]*model.UserMembershipView, 0)
	query := repository.PrepareSearchQuery(table, model.UserMembershipSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func PutUserMemberships(db *gorm.DB, table string, users ...*model.UserMembershipView) error {
	save := repository.PrepareBulkSave(table)
	u := make([]interface{}, len(users))
	for i, user := range users {
		u[i] = user
	}
	return save(db, u...)
}

func PutUserMembership(db *gorm.DB, table string, user *model.UserMembershipView) error {
	save := repository.PrepareSave(table)
	return save(db, user)
}

func DeleteUserMembership(db *gorm.DB, table, userID, aggregateID, objectID string, membertype usr_model.MemberType) error {
	delete := repository.PrepareDeleteByKeys(table,
		repository.Key{Key: model.UserMembershipSearchKey(usr_model.UserMembershipSearchKeyUserID), Value: userID},
		repository.Key{Key: model.UserMembershipSearchKey(usr_model.UserMembershipSearchKeyAggregateID), Value: aggregateID},
		repository.Key{Key: model.UserMembershipSearchKey(usr_model.UserMembershipSearchKeyObjectID), Value: objectID},
		repository.Key{Key: model.UserMembershipSearchKey(usr_model.UserMembershipSearchKeyMemberType), Value: membertype},
	)
	return delete(db)
}
