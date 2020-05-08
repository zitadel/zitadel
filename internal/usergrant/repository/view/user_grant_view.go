package view

import (
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

func SearchUserGrants(db *gorm.DB, table string, req *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, int, error) {
	users := make([]*model.UserGrantView, 0)
	query := view.PrepareSearchQuery(table, model.UserGrantSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &users)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func PutUserGrant(db *gorm.DB, table string, grant *model.UserGrantView) error {
	save := view.PrepareSave(table)
	return save(db, grant)
}

func DeleteUserGrant(db *gorm.DB, table, grantID string) error {
	delete := view.PrepareDeleteByKey(table, model.UserGrantSearchKey(grant_model.USERGRANTSEARCHKEY_GRANT_ID), grantID)
	return delete(db)
}
