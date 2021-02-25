package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectMemberByIDs(db *gorm.DB, table, projectID, userID string) (*model.ProjectMemberView, error) {
	role := new(model.ProjectMemberView)

	projectIDQuery := model.ProjectMemberSearchQuery{Key: proj_model.ProjectMemberSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals}
	userIDQuery := model.ProjectMemberSearchQuery{Key: proj_model.ProjectMemberSearchKeyUserID, Value: userID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, userIDQuery)
	err := query(db, role)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-EgWQ2", "Errors.Project.Member.NotExisting")
	}
	return role, err
}

func ProjectMembersByProjectID(db *gorm.DB, table, projectID string) ([]*model.ProjectMemberView, error) {
	members := make([]*model.ProjectMemberView, 0)
	queries := []*proj_model.ProjectMemberSearchQuery{
		&proj_model.ProjectMemberSearchQuery{Key: proj_model.ProjectMemberSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func SearchProjectMembers(db *gorm.DB, table string, req *proj_model.ProjectMemberSearchRequest) ([]*model.ProjectMemberView, uint64, error) {
	roles := make([]*model.ProjectMemberView, 0)
	query := repository.PrepareSearchQuery(table, model.ProjectMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}
func ProjectMembersByUserID(db *gorm.DB, table string, userID string) ([]*model.ProjectMemberView, error) {
	members := make([]*model.ProjectMemberView, 0)
	queries := []*proj_model.ProjectMemberSearchQuery{
		&proj_model.ProjectMemberSearchQuery{Key: proj_model.ProjectMemberSearchKeyUserID, Value: userID, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutProjectMember(db *gorm.DB, table string, role *model.ProjectMemberView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func PutProjectMembers(db *gorm.DB, table string, members ...*model.ProjectMemberView) error {
	save := repository.PrepareBulkSave(table)
	m := make([]interface{}, len(members))
	for i, member := range members {
		m[i] = member
	}
	return save(db, m...)
}

func DeleteProjectMember(db *gorm.DB, table, projectID, userID string) error {
	role, err := ProjectMemberByIDs(db, table, projectID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, role)
	return delete(db)
}

func DeleteProjectMembersByProjectID(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectMemberSearchKey(proj_model.ProjectMemberSearchKeyProjectID), projectID)
	return delete(db)
}
