package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectGrantMemberByIDs(db *gorm.DB, table, grantID, userID string) (*model.ProjectGrantMemberView, error) {
	grant := new(model.ProjectGrantMemberView)

	grantIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyGrantID, Value: grantID, Method: domain.SearchMethodEquals}
	userIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyUserID, Value: userID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, grantIDQuery, userIDQuery)
	err := query(db, grant)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Sgr32", "Errors.Project.Member.NotExisting")
	}
	return grant, err
}

func ProjectGrantMembersByProjectID(db *gorm.DB, table, projectID string) ([]*model.ProjectGrantMemberView, error) {
	members := make([]*model.ProjectGrantMemberView, 0)
	queries := []*proj_model.ProjectGrantMemberSearchQuery{
		{Key: proj_model.ProjectGrantMemberSearchKeyProjectID, Value: projectID, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectGrantMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func SearchProjectGrantMembers(db *gorm.DB, table string, req *proj_model.ProjectGrantMemberSearchRequest) ([]*model.ProjectGrantMemberView, uint64, error) {
	roles := make([]*model.ProjectGrantMemberView, 0)
	query := repository.PrepareSearchQuery(table, model.ProjectGrantMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func ProjectGrantMembersByUserID(db *gorm.DB, table, userID string) ([]*model.ProjectGrantMemberView, error) {
	members := make([]*model.ProjectGrantMemberView, 0)
	queries := []*proj_model.ProjectGrantMemberSearchQuery{
		{Key: proj_model.ProjectGrantMemberSearchKeyUserID, Value: userID, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectGrantMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutProjectGrantMember(db *gorm.DB, table string, role *model.ProjectGrantMemberView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func PutProjectGrantMembers(db *gorm.DB, table string, members ...*model.ProjectGrantMemberView) error {
	save := repository.PrepareBulkSave(table)
	m := make([]interface{}, len(members))
	for i, member := range members {
		m[i] = member
	}
	return save(db, m...)
}

func DeleteProjectGrantMember(db *gorm.DB, table, grantID, userID string) error {
	grant, err := ProjectGrantMemberByIDs(db, table, grantID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, grant)
	return delete(db)
}

func DeleteProjectGrantMembersByProjectID(db *gorm.DB, table, projectID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectGrantMemberSearchKey(proj_model.ProjectGrantMemberSearchKeyProjectID), projectID)
	return delete(db)
}

func DeleteProjectGrantMembersByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.ProjectGrantMemberSearchKey(proj_model.ProjectGrantMemberSearchKeyUserID), userID)
	return delete(db)
}
