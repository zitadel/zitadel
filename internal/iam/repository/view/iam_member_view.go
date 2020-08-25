package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func IamMemberByIDs(db *gorm.DB, table, orgID, userID string) (*model.IamMemberView, error) {
	member := new(model.IamMemberView)

	iamIDQuery := &model.IamMemberSearchQuery{Key: iam_model.IamMemberSearchKeyIamID, Value: orgID, Method: global_model.SearchMethodEquals}
	userIDQuery := &model.IamMemberSearchQuery{Key: iam_model.IamMemberSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, iamIDQuery, userIDQuery)
	err := query(db, member)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Ahq2s", "Errors.Iam.MemberNotExisting")
	}
	return member, err
}

func SearchIamMembers(db *gorm.DB, table string, req *iam_model.IamMemberSearchRequest) ([]*model.IamMemberView, uint64, error) {
	members := make([]*model.IamMemberView, 0)
	query := repository.PrepareSearchQuery(table, model.IamMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}
func IamMembersByUserID(db *gorm.DB, table string, userID string) ([]*model.IamMemberView, error) {
	members := make([]*model.IamMemberView, 0)
	queries := []*iam_model.IamMemberSearchQuery{
		{
			Key:    iam_model.IamMemberSearchKeyUserID,
			Value:  userID,
			Method: global_model.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IamMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutIamMember(db *gorm.DB, table string, role *model.IamMemberView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func PutIamMembers(db *gorm.DB, table string, members ...*model.IamMemberView) error {
	save := repository.PrepareBulkSave(table)
	m := make([]interface{}, len(members))
	for i, member := range members {
		m[i] = member
	}
	return save(db, m...)
}

func DeleteIamMember(db *gorm.DB, table, orgID, userID string) error {
	member, err := IamMemberByIDs(db, table, orgID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, member)
	return delete(db)
}
