package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func GrantedProjectByIDs(db *gorm.DB, table, projectID, orgID string) (*model.GrantedProjectView, error) {
	project := new(model.GrantedProjectView)

	projectIDQuery := model.GrantedProjectSearchQuery{Key: proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS}
	orgIDQuery := model.GrantedProjectSearchQuery{Key: proj_model.GRANTEDPROJECTSEARCHKEY_ORGID, Value: orgID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, projectIDQuery, orgIDQuery)
	err := query(db, project)
	return project, err
}

func GrantedProjectGrantByIDs(db *gorm.DB, table, projectID, grantID string) (*model.GrantedProjectView, error) {
	project := new(model.GrantedProjectView)

	projectIDQuery := model.GrantedProjectSearchQuery{Key: proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS}
	grantIDQuery := model.GrantedProjectSearchQuery{Key: proj_model.GRANTEDPROJECTSEARCHKEY_GRANTID, Value: grantID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, projectIDQuery, grantIDQuery)
	err := query(db, project)
	return project, err
}
func GrantedProjectsByID(db *gorm.DB, table, projectID string) ([]*model.GrantedProjectView, error) {
	projects := make([]*model.GrantedProjectView, 0)
	queries := []*proj_model.GrantedProjectSearchQuery{
		&proj_model.GrantedProjectSearchQuery{Key: proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID, Value: projectID, Method: global_model.SEARCHMETHOD_EQUALS},
	}
	query := view.PrepareSearchQuery(table, model.GrantedProjectSearchRequest{Queries: queries})
	_, err := query(db, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func SearchGrantedProjects(db *gorm.DB, table string, req *proj_model.GrantedProjectSearchRequest) ([]*model.GrantedProjectView, int, error) {
	projects := make([]*model.GrantedProjectView, 0)
	query := view.PrepareSearchQuery(table, model.GrantedProjectSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &projects)
	if err != nil {
		return nil, 0, err
	}
	return projects, count, nil
}

func PutGrantedProject(db *gorm.DB, table string, project *model.GrantedProjectView) error {
	save := view.PrepareSave(table)
	return save(db, project)
}

func DeleteGrantedProject(db *gorm.DB, table, projectID, orgID string) error {
	project, err := GrantedProjectByIDs(db, table, projectID, orgID)
	if err != nil {
		return err
	}
	delete := view.PrepareDeleteByObject(table, project)
	return delete(db)
}
