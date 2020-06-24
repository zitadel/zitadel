package view

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view"
	"github.com/jinzhu/gorm"
)

func IamMemberByIDs(db *gorm.DB, table, orgID, userID string) (*model.IamMemberView, error) {
	member := new(model.IamMemberView)

	orgIDQuery := &model.IamMemberSearchQuery{Key: iam_model.IamMemberSearchKeyIamID, Value: orgID, Method: global_model.SearchMethodEquals}
	userIDQuery := &model.IamMemberSearchQuery{Key: iam_model.IamMemberSearchKeyUserID, Value: userID, Method: global_model.SearchMethodEquals}
	query := view.PrepareGetByQuery(table, orgIDQuery, userIDQuery)
	err := query(db, member)
	return member, err
}

func SearchIamMembers(db *gorm.DB, table string, req *iam_model.IamMemberSearchRequest) ([]*model.IamMemberView, int, error) {
	members := make([]*model.IamMemberView, 0)
	query := view.PrepareSearchQuery(table, model.IamMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
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
	query := view.PrepareSearchQuery(table, model.IamMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutIamMember(db *gorm.DB, table string, role *model.IamMemberView) error {
	save := view.PrepareSave(table)
	return save(db, role)
}

func DeleteIamMember(db *gorm.DB, table, orgID, userID string) error {
	member, err := IamMemberByIDs(db, table, orgID, userID)
	if err != nil {
		return err
	}
	delete := view.PrepareDeleteByObject(table, member)
	return delete(db)
}
