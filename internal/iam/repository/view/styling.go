package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func GetStylingByAggregateIDAndState(db *gorm.DB, table, aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	policy := new(model.LabelPolicyView)
	aggregateIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	stateQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyState, Value: state, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	ownerRemovedQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyOwnerRemoved, Value: false, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, stateQuery, instanceIDQuery, ownerRemovedQuery)
	err := query(db, policy)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-68G11", "Errors.IAM.LabelPolicy.NotExisting")
	}
	return policy, err
}
