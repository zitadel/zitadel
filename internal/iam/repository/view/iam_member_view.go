package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func IAMMemberByIDs(db *gorm.DB, table, orgID, userID string) (*model.IAMMemberView, error) {
	member := new(model.IAMMemberView)

	iamIDQuery := &model.IAMMemberSearchQuery{Key: iam_model.IAMMemberSearchKeyIamID, Value: orgID, Method: domain.SearchMethodEquals}
	userIDQuery := &model.IAMMemberSearchQuery{Key: iam_model.IAMMemberSearchKeyUserID, Value: userID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, iamIDQuery, userIDQuery)
	err := query(db, member)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Ahq2s", "Errors.IAM.MemberNotExisting")
	}
	return member, err
}

func SearchIAMMembers(db *gorm.DB, table string, req *iam_model.IAMMemberSearchRequest) ([]*model.IAMMemberView, uint64, error) {
	members := make([]*model.IAMMemberView, 0)
	query := repository.PrepareSearchQuery(table, model.IAMMemberSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &members)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}
func IAMMembersByUserID(db *gorm.DB, table string, userID string) ([]*model.IAMMemberView, error) {
	members := make([]*model.IAMMemberView, 0)
	queries := []*iam_model.IAMMemberSearchQuery{
		{
			Key:    iam_model.IAMMemberSearchKeyUserID,
			Value:  userID,
			Method: domain.SearchMethodEquals,
		},
	}
	query := repository.PrepareSearchQuery(table, model.IAMMemberSearchRequest{Queries: queries})
	_, err := query(db, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func PutIAMMember(db *gorm.DB, table string, role *model.IAMMemberView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func PutIAMMembers(db *gorm.DB, table string, members ...*model.IAMMemberView) error {
	save := repository.PrepareBulkSave(table)
	m := make([]interface{}, len(members))
	for i, member := range members {
		m[i] = member
	}
	return save(db, m...)
}

func DeleteIAMMember(db *gorm.DB, table, orgID, userID string) error {
	member, err := IAMMemberByIDs(db, table, orgID, userID)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, member)
	return delete(db)
}

func DeleteIAMMembersByUserID(db *gorm.DB, table, userID string) error {
	delete := repository.PrepareDeleteByKey(table, model.IAMMemberSearchKey(iam_model.IAMMemberSearchKeyUserID), userID)
	return delete(db)
}
