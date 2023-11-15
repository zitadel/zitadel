package view

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
	iam_model "github.com/zitadel/zitadel/v2/internal/iam/model"
	"github.com/zitadel/zitadel/v2/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/v2/internal/view/repository"

	"github.com/jinzhu/gorm"
)

func GetStylingByAggregateIDAndState(db *gorm.DB, table, aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	policy := new(model.LabelPolicyView)
	aggregateIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	stateQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyState, Value: state, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	ownerRemovedQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyOwnerRemoved, Value: false, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, stateQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-68G11", "Errors.IAM.LabelPolicy.NotExisting")
	}
	return policy, err
}

func PutStyling(db *gorm.DB, table string, policy *model.LabelPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func UpdateOrgOwnerRemovedStyling(db *gorm.DB, table, instanceID, aggID string) error {
	update := repository.PrepareUpdateByKeys(table,
		model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyOwnerRemoved),
		true,
		repository.Key{Key: model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID), Value: instanceID},
		repository.Key{Key: model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyAggregateID), Value: aggID},
	)
	return update(db)
}

func DeleteInstanceStyling(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table, model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID), instanceID)
	return delete(db)
}
