package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectGrantMemberByIDs(db *gorm.DB, table, grantID, userID string) (*model.ProjectGrantMemberView, error) {
	role := new(model.ProjectGrantMemberView)

	grantIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyGrantID, Value: grantID, Method: global_model.SearchMethodEquals}
	userIDQuery := model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, grantIDQuery, userIDQuery)
	err := query(db, role)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Sgr32", "Errors.Project.MemberNotExisting")
	}
	return role, err
}

func SearchProjectGrantMembers(db *gorm.DB, table string, req *proj_model.ProjectGrantMemberSearchRequest) ([]*model.ProjectGrantMemberView, int, error) {
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
		&proj_model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals},
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

func DeleteProjectGrantMember(db *gorm.DB, table, grantID, userID string) error {
	role, err := ProjectGrantMemberByIDs(db, table, grantID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, role)
	return delete(db)
}
