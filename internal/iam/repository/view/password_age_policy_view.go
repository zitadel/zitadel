package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetPasswordAgePolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.PasswordAgePolicyView, error) {
	policy := new(model.PasswordAgePolicyView)
	aggregateIDQuery := &model.PasswordAgePolicySearchQuery{Key: iam_model.PasswordAgePolicySearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Lso0cs", "Errors.IAM.PasswordAgePolicy.NotExisting")
	}
	return policy, err
}

func PutPasswordAgePolicy(db *gorm.DB, table string, policy *model.PasswordAgePolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeletePasswordAgePolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.PasswordAgePolicySearchKey(iam_model.PasswordAgePolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
