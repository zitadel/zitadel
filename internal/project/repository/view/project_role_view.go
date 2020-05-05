package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func ProjectRoleByIDs(db *gorm.DB, table, projectID, orgID, key string) (*model.ProjectRoleView, error) {
	role := new(model.ProjectRoleView)

	projectIDQuery := model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS}
	grantIDQuery := model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_ORGID, Value: orgID, Method: global_model.SEARCHMETHOD_EQUALS}
	keyQuery := model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_KEY, Value: orgID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, projectIDQuery, grantIDQuery, keyQuery)
	err := query(db, role)
	return role, err
}

func ResourceOwnerProjectRolesByKey(db *gorm.DB, table, projectID, resourceOwner, key string) ([]*model.ProjectRoleView, error) {
	roles := make([]*model.ProjectRoleView, 0)
	queries := []*proj_model.ProjectRoleSearchQuery{
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_RESOURCEOWNER, Value: resourceOwner, Method: global_model.SEARCHMETHOD_EQUALS},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_KEY, Value: key, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Queries: queries})
	_, err := query(db, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func ResourceOwnerProjectRoles(db *gorm.DB, table, projectID, resourceOwner string) ([]*model.ProjectRoleView, error) {
	roles := make([]*model.ProjectRoleView, 0)
	queries := []*proj_model.ProjectRoleSearchQuery{
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.PROJECTROLESEARCHKEY_RESOURCEOWNER, Value: resourceOwner, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Queries: queries})
	_, err := query(db, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func SearchProjectRoles(db *gorm.DB, table string, req *proj_model.ProjectRoleSearchRequest) ([]*model.ProjectRoleView, int, error) {
	roles := make([]*model.ProjectRoleView, 0)
	query := view.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func PutProjectRole(db *gorm.DB, table string, role *model.ProjectRoleView) error {
	save := view.PrepareSave(table)
	return save(db, role)
}

func DeleteProjectRole(db *gorm.DB, table, projectID, orgID, key string) error {
	role, err := ProjectRoleByIDs(db, table, projectID, orgID, key)
	if err != nil {
		return err
	}
	delete := view.PrepareDeleteByObject(table, role)
	return delete(db)
}
