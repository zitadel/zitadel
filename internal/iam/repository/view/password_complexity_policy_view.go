package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetPasswordComplexityPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.PasswordComplexityPolicyView, error) {
	policy := new(model.PasswordComplexityPolicyView)
	aggregateIDQuery := &model.PasswordComplexityPolicySearchQuery{Key: iam_model.PasswordComplexityPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Lso0cs", "Errors.IAM.PasswordComplexityPolicy.NotExisting")
	}
	return policy, err
}

func PutPasswordComplexityPolicy(db *gorm.DB, table string, policy *model.PasswordComplexityPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeletePasswordComplexityPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.PasswordComplexityPolicySearchKey(iam_model.PasswordComplexityPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
