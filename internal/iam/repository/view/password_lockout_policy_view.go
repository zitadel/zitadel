package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetLockoutPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.LockoutPolicyView, error) {
	policy := new(model.LockoutPolicyView)
	aggregateIDQuery := &model.LockoutPolicySearchQuery{Key: iam_model.LockoutPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-M9fsf", "Errors.IAM.LockoutPolicy.NotExisting")
	}
	return policy, err
}

func PutLockoutPolicy(db *gorm.DB, table string, policy *model.LockoutPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeleteLockoutPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.LockoutPolicySearchKey(iam_model.LockoutPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
