package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func OrgMemberByIDs(db *gorm.DB, table, orgID, userID string) (*model.OrgMemberView, error) {
	member := new(model.OrgMemberView)

	orgIDQuery := &model.OrgMemberSearchQuery{Key: org_model.OrgMemberSearchKeyOrgID, Value: orgID, Method: global_model.SearchMethodEquals}
	userIDQuery := &model.OrgMemberSearchQuery{Key: org_model.OrgMemberSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, orgIDQuery, userIDQuery)
	err := query(db, member)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-gIaTM", "Errors.Org.MemberNotFound")
	}
	return member, err
}

func SearchOrgMembers(db *gorm.DB, table string, req *org_model.OrgMemberSearchRequest) ([]*model.OrgMemberView, uint64, error) {
	members := make([]*model.OrgMemberView, 0)
	query := repository.PrepareSearchQuery(table, model.OrgMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}
func OrgMembersByUserID(db *gorm.DB, table string, userID string) ([]*model.OrgMemberView, error) {
	members := make([]*model.OrgMemberView, 0)
	queries := []*org_model.OrgMemberSearchQuery{
		{
			Key:    org_model.OrgMemberSearchKeyUserID,
			Value:  userID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.OrgMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutOrgMember(db *gorm.DB, table string, role *model.OrgMemberView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func PutOrgMembers(db *gorm.DB, table string, members ...*model.OrgMemberView) error {
	save := repository.PrepareBulkSave(table)
	m := make([]interface{}, len(members))
	for i, member := range members {
		m[i] = member
	}
	return save(db, m...)
}

func DeleteOrgMember(db *gorm.DB, table, orgID, userID string) error {
	member, err := OrgMemberByIDs(db, table, orgID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, member)
	return delete(db)
}

func DeleteOrgMembersByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgMemberSearchKey(org_model.OrgMemberSearchKeyUserID), userID)
	return delete(db)
}
