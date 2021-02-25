package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetDefaultLoginPolicies(db *gorm.DB, table string) ([]*model.LoginPolicyView, error) {
	loginPolicies := make([]*model.LoginPolicyView, 0)
	queries := []*iam_model.LoginPolicySearchQuery{
		{Key: iam_model.LoginPolicySearchKeyDefault, Value: true, Method: domain.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.LoginPolicySearchRequest{Queries: queries})
	_, err := query(db, &loginPolicies)
	if err != nil {
		return nil, err
	}
	return loginPolicies, nil
}

func GetLoginPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.LoginPolicyView, error) {
	policy := new(model.LoginPolicyView)
	aggregateIDQuery := &model.LoginPolicySearchQuery{Key: iam_model.LoginPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
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

func PutLoginPolicies(db *gorm.DB, table string, policies ...*model.LoginPolicyView) error {
	save := repository.PrepareBulkSave(table)
	u := make([]interface{}, len(policies))
	for i, user := range policies {
		u[i] = user
	}
	return save(db, u...)
}

func DeleteLoginPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.LoginPolicySearchKey(iam_model.LoginPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
