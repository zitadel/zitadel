package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func OrgMemberByIDs(db *gorm.DB, table, orgID, userID string) (*model.OrgMemberView, error) {
	member := new(model.OrgMemberView)

	orgIDQuery := &model.OrgMemberSearchQuery{Key: org_model.ORGMEMBERSEARCHKEY_ORG_ID, Value: orgID, Method: global_model.SEARCHMETHOD_EQUALS}
	userIDQuery := &model.OrgMemberSearchQuery{Key: org_model.ORGMEMBERSEARCHKEY_USER_ID, Value: userID, Method: global_model.SEARCHMETHOD_EQUALS}
	query := view.PrepareGetByQuery(table, orgIDQuery, userIDQuery)
	err := query(db, member)
	return member, err
}

func SearchOrgMembers(db *gorm.DB, table string, req *org_model.OrgMemberSearchRequest) ([]*model.OrgMemberView, int, error) {
	members := make([]*model.OrgMemberView, 0)
	query := view.PrepareSearchQuery(table, model.OrgMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
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
			Key:    org_model.ORGMEMBERSEARCHKEY_USER_ID,
			Value:  userID,
			Method: global_model.SEARCHMETHOD_EQUALS,
		},
	}
	query := view.PrepareSearchQuery(table, model.OrgMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutOrgMember(db *gorm.DB, table string, role *model.OrgMemberView) error {
	save := view.PrepareSave(table)
	return save(db, role)
}

func DeleteOrgMember(db *gorm.DB, table, orgID, userID string) error {
	member, err := OrgMemberByIDs(db, table, orgID, userID)
	if err != nil {
		return err
	}
	delete := view.PrepareDeleteByObject(table, member)
	return delete(db)
}
