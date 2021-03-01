package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectByID(db *gorm.DB, table, projectID string) (*model.ProjectView, error) {
	project := new(model.ProjectView)

	projectIDQuery := model.ProjectSearchQuery{Key: proj_model.ProjectViewSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery)
	err := query(db, project)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-NEO7W", "Errors.Project.NotFound")
	}
	return project, err
}

func ProjectsByResourceOwner(db *gorm.DB, table, orgID string) ([]*model.ProjectView, error) {
	projects := make([]*model.ProjectView, 0)
	queries := []*proj_model.ProjectViewSearchQuery{
		{Key: proj_model.ProjectViewSearchKeyResourceOwner, Value: orgID, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectSearchRequest{Queries: queries})
	_, err := query(db, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func SearchProjects(db *gorm.DB, table string, req *proj_model.ProjectViewSearchRequest) ([]*model.ProjectView, uint64, error) {
	projects := make([]*model.ProjectView, 0)
	query := repository.PrepareSearchQuery(table, model.ProjectSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &projects)
	if err != nil {
		return nil, 0, err
	}
	return projects, count, nil
}

func PutProject(db *gorm.DB, table string, project *model.ProjectView) error {
	save := repository.PrepareSave(table)
	return save(db, project)
}

func DeleteProject(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectSearchKey(proj_model.ProjectViewSearchKeyProjectID), projectID)
	return delete(db)
}
