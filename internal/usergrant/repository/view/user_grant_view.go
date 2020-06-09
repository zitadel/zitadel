package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func UserGrantByID(db *gorm.DB, table, grantID string) (*model.UserGrantView, error) {
	user := new(model.UserGrantView)
	query := view.PrepareGetByKey(table, model.UserGrantSearchKey(grant_model.USERGRANTSEARCHKEY_GRANT_ID), grantID)
	err := query(db, user)
	return user, err
}

func UserGrantByIDs(db *gorm.DB, table, resourceOwnerID, projectID, userID string) (*model.UserGrantView, error) {
	user := new(model.UserGrantView)

	resourceOwnerIDQuery := model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER, Value: resourceOwnerID, Method: global_model.SEARCHMETHOD_EQUALS}
	projectIDQuery := model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_PROJECT_ID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS}
	userIDQuery := model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_USER_ID, Value: userID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, resourceOwnerIDQuery, projectIDQuery, userIDQuery)
	err := query(db, user)
	return user, err
}

func SearchUserGrants(db *gorm.DB, table string, req *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, int, error) {
	users := make([]*model.UserGrantView, 0)
	query := view.PrepareSearchQuery(table, model.UserGrantSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func UserGrantsByUserID(db *gorm.DB, table, userID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		&grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_USER_ID, Value: userID, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByProjectID(db *gorm.DB, table, projectID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		&grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_PROJECT_ID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserGrantsByOrgID(db *gorm.DB, table, orgID string) ([]*model.UserGrantView, error) {
	users := make([]*model.UserGrantView, 0)
	queries := []*grant_model.UserGrantSearchQuery{
		&grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER, Value: orgID, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.UserGrantSearchRequest{Queries: queries})
	_, err := query(db, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func PutUserGrant(db *gorm.DB, table string, grant *model.UserGrantView) error {
	save := view.PrepareSave(table)
	return save(db, grant)
}

func DeleteUserGrant(db *gorm.DB, table, grantID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserGrantSearchKey(grant_model.USERGRANTSEARCHKEY_GRANT_ID), grantID)
	return delete(db)
}
