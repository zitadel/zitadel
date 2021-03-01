package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetOrgIAMPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.OrgIAMPolicyView, error) {
	policy := new(model.OrgIAMPolicyView)
	aggregateIDQuery := &model.OrgIAMPolicySearchQuery{Key: iam_model.OrgIAMPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-5fi9s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return policy, err
}

func PutOrgIAMPolicy(db *gorm.DB, table string, policy *model.OrgIAMPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeleteOrgIAMPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.OrgIAMPolicySearchKey(iam_model.OrgIAMPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
