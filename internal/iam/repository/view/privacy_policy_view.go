package view

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func GetPrivacyPolicyByAggregateID(db *gorm.DB, table, aggregateID string) (*model.PrivacyPolicyView, error) {
	policy := new(model.PrivacyPolicyView)
	aggregateIDQuery := &model.PrivacyPolicySearchQuery{Key: iam_model.PrivacyPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-2N9fs", "Errors.IAM.PrivacyPolicy.NotExisting")
	}
	return policy, err
}

func PutPrivacyPolicy(db *gorm.DB, table string, policy *model.PrivacyPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeletePrivacyPolicy(db *gorm.DB, table, aggregateID string) error {
	delete := repository.PrepareDeleteByKey(table, model.PrivacyPolicySearchKey(iam_model.PrivacyPolicySearchKeyAggregateID), aggregateID)

	return delete(db)
}
