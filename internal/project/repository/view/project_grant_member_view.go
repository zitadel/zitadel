package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func ProjectGrantMemberByIDs(db *gorm.DB, table, grantID, userID string) (*model.ProjectGrantMemberView, error) {
	role := new(model.ProjectGrantMemberView)

	grantIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.PROJECTGRANTMEMBERSEARCHKEY_GRANT_ID, Value: grantID, Method: global_model.SEARCHMETHOD_EQUALS}
	userIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.PROJECTGRANTMEMBERSEARCHKEY_USER_ID, Value: userID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, grantIDQuery, userIDQuery)
	err := query(db, role)
	return role, err
}

func SearchProjectGrantMembers(db *gorm.DB, table string, req *proj_model.ProjectGrantMemberSearchRequest) ([]*model.ProjectGrantMemberView, int, error) {
	roles := make([]*model.ProjectGrantMemberView, 0)
	query := view.PrepareSearchQuery(table, model.ProjectGrantMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func PutProjectGrantMember(db *gorm.DB, table string, role *model.ProjectGrantMemberView) error {
	save := view.PrepareSave(table)
	return save(db, role)
}

func DeleteProjectGrantMember(db *gorm.DB, table, grantID, userID string) error {
	role, err := ProjectGrantMemberByIDs(db, table, grantID, userID)
	if err != nil {
		return err
	}
	delete := view.PrepareDeleteByObject(table, role)
	return delete(db)
}
