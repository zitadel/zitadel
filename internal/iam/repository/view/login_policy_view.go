package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetLoginPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.LoginPolicyView, error) {
	policy := new(model.LoginPolicyView)
	aggregateIDQuery := &model.LoginPolicySearchQuery{Key: iam_model.LoginPolicySearchKeyAggregateID, Value: aggregateID, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-Lso0cs", "Errors.IAM.LoginPolicy.NotExisting")
	}
	return policy, err
}

func PutLoginPolicy(db *gorm.DB, table string, policy *model.LoginPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeleteLoginPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.LoginPolicySearchKey(iam_model.LoginPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
