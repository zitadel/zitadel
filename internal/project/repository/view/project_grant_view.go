package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectGrantByProjectAndOrg(db *gorm.DB, table, projectID, orgID string) (*model.ProjectGrantView, error) {
	projectGrant := new(model.ProjectGrantView)

	projectIDQuery := model.ProjectGrantSearchQuery{Key: proj_model.GrantedProjectSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	orgIDQuery := model.ProjectGrantSearchQuery{Key: proj_model.GrantedProjectSearchKeyOrgID, Value: orgID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, orgIDQuery)
	err := query(db, projectGrant)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-WR3z2", "Errors.Project.Grant.NotExisting")
	}
	return projectGrant, err
}

func ProjectGrantByID(db *gorm.DB, table, grantID string) (*model.ProjectGrantView, error) {
	projectGrant := new(model.ProjectGrantView)
	grantIDQuery := model.ProjectGrantSearchQuery{Key: proj_model.GrantedProjectSearchKeyGrantID, Value: grantID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, grantIDQuery)
	err := query(db, projectGrant)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-EGdh4", "Errors.Project.Grant.NotFound")
	}
	return projectGrant, err
}

func ProjectGrantsByProjectID(db *gorm.DB, table, projectID string) ([]*model.ProjectGrantView, error) {
	projectGrants := make([]*model.ProjectGrantView, 0)
	queries := []*proj_model.ProjectGrantViewSearchQuery{
		{Key: proj_model.GrantedProjectSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectGrantSearchRequest{Queries: queries})
	_, err := query(db, &projectGrants)
	return projectGrants, err
}

func ProjectGrantsByProjectIDAndRoleKey(db *gorm.DB, table, projectID, roleKey string) ([]*model.ProjectGrantView, error) {
	projectGrants := make([]*model.ProjectGrantView, 0)
	queries := []*proj_model.ProjectGrantViewSearchQuery{
		{Key: proj_model.GrantedProjectSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
		{Key: proj_model.GrantedProjectSearchKeyRoleKeys, Value: roleKey, Method: global_model.SearchMethodListContains},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectGrantSearchRequest{Queries: queries})
	_, err := query(db, &projectGrants)

	return projectGrants, err
}

func SearchProjectGrants(db *gorm.DB, table string, req *proj_model.ProjectGrantViewSearchRequest) ([]*model.ProjectGrantView, uint64, error) {
	projectGrants := make([]*model.ProjectGrantView, 0)
	query := repository.PrepareSearchQuery(table, model.ProjectGrantSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &projectGrants)

	return projectGrants, count, err
}

func PutProjectGrant(db *gorm.DB, table string, grant *model.ProjectGrantView) error {
	save := repository.PrepareSave(table)
	return save(db, grant)
}

func PutProjectGrants(db *gorm.DB, table string, grants ...*model.ProjectGrantView) error {
	save := repository.PrepareBulkSave(table)
	g := make([]interface{}, len(grants))
	for i, grant := range grants {
		g[i] = grant
	}
	return save(db, g...)
}

func DeleteProjectGrant(db *gorm.DB, table, grantID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectGrantSearchKey(proj_model.GrantedProjectSearchKeyGrantID), grantID)
	return delete(db)
}

func DeleteProjectGrantsByProjectID(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectGrantSearchKey(proj_model.GrantedProjectSearchKeyProjectID), projectID)
	return delete(db)
}
