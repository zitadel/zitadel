package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func UserGrantByID(db *gorm.DB, table, grantID string) (*model.UserGrantView, error) {
	user := new(model.UserGrantView)
	query := repository.PrepareGetByKey(table, model.UserGrantSearchKey(grant_model.UserGrantSearchKeyGrantID), grantID)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Nqwf1", "Errors.Token.NotFound")
	}
	return user, err
}

func UserGrantByIDs(db *gorm.DB, table, resourceOwnerID, projectID, userID string) (*model.UserGrantView, error) {
	user := new(model.UserGrantView)

	resourceOwnerIDQuery := model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyResourceOwner, Value: resourceOwnerID, Method: global_model.SearchMethodEquals}
	projectIDQuery := model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	userIDQuery := model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, resourceOwnerIDQuery, projectIDQuery, userIDQuery)
	err := query(db, user)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Q1tq2", "Errors.UserGrant.NotFound")
	}
	return user, err
}

func SearchUserGrants(db *gorm.DB, table string, req *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, uint64, error) {
	users := make([]*model.UserGrantView, 0)
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func UserGrantsByUserID(db *gorm.DB, table, userID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByProjectID(db *gorm.DB, table, projectID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByProjectAndUserID(db *gorm.DB, table, projectID, userID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
		{Key: grant_model.UserGrantSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByProjectIDAndRole(db *gorm.DB, table, projectID, roleKey string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
		{Key: grant_model.UserGrantSearchKeyRoleKey, Value: roleKey, Method: global_model.SearchMethodListContains},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByOrgIDAndProjectID(db *gorm.DB, table, orgID, projectID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyResourceOwner, Value: orgID, Method: global_model.SearchMethodEquals},
		{Key: grant_model.UserGrantSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByOrgID(db *gorm.DB, table, orgID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		{Key: grant_model.UserGrantSearchKeyResourceOwner, Value: orgID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func PutUserGrant(db *gorm.DB, table string, grant *model.UserGrantView) error {
	save := repository.PrepareSave(table)
	return save(db, grant)
}

func DeleteUserGrant(db *gorm.DB, table, grantID string) error {
	delete := repository.PrepareDeleteByKey(table, model.UserGrantSearchKey(grant_model.UserGrantSearchKeyGrantID), grantID)
	return delete(db)
}
