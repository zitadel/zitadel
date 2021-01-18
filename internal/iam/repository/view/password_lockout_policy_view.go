package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetPasswordLockoutPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.PasswordLockoutPolicyView, error) {
	policy := new(model.PasswordLockoutPolicyView)
	aggregateIDQuery := &model.PasswordLockoutPolicySearchQuery{Key: iam_model.PasswordLockoutPolicySearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Lso0cs", "Errors.IAM.PasswordLockoutPolicy.NotExisting")
	}
	return policy, err
}

func PutPasswordLockoutPolicy(db *gorm.DB, table string, policy *model.PasswordLockoutPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeletePasswordLockoutPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.PasswordLockoutPolicySearchKey(iam_model.PasswordLockoutPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
