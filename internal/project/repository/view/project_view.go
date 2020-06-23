package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func ProjectByID(db *gorm.DB, table, projectID string) (*model.ProjectView, error) {
	project := new(model.ProjectView)

	projectIDQuery := model.ProjectSearchQuery{Key: proj_model.ProjectViewSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	query := view.PrepareGetByQuery(table, projectIDQuery)
	err := query(db, project)
	return project, err
}

func ProjectsByResourceOwner(db *gorm.DB, table, orgID string) ([]*model.ProjectView, error) {
	projects := make([]*model.ProjectView, 0)
	queries := []*proj_model.ProjectViewSearchQuery{
		&proj_model.ProjectViewSearchQuery{Key: proj_model.ProjectViewSearchKeyResourceOwner, Value: orgID, Method: global_model.SearchMethodEquals},
	}
	query := view.PrepareSearchQuery(table, model.ProjectSearchRequest{Queries: queries})
	_, err := query(db, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func SearchProjects(db *gorm.DB, table string, req *proj_model.ProjectViewSearchRequest) ([]*model.ProjectView, int, error) {
	projects := make([]*model.ProjectView, 0)
	query := view.PrepareSearchQuery(table, model.ProjectSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &projects)
	if err != nil {
		return nil, 0, err
	}
	return projects, count, nil
}

func PutProject(db *gorm.DB, table string, project *model.ProjectView) error {
	save := view.PrepareSave(table)
	return save(db, project)
}

func DeleteProject(db *gorm.DB, table, projectID string) error {
	delete := view.PrepareDeleteByKey(table, model.ProjectSearchKey(proj_model.ProjectViewSearchKeyProjectID), projectID)
	return delete(db)
}
