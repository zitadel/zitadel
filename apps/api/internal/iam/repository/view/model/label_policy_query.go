package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

type LabelPolicySearchQuery iam_model.LabelPolicySearchQuery
type LabelPolicySearchKey iam_model.LabelPolicySearchKey

func (req LabelPolicySearchQuery) GetKey() repository.ColumnKey {
	return LabelPolicySearchKey(req.Key)
}

func (req LabelPolicySearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req LabelPolicySearchQuery) GetValue() interface{} {
	return req.Value
}

func (key LabelPolicySearchKey) ToColumnName() string {
	switch iam_model.LabelPolicySearchKey(key) {
	case iam_model.LabelPolicySearchKeyAggregateID:
		return LabelPolicyKeyAggregateID
	case iam_model.LabelPolicySearchKeyState:
		return LabelPolicyKeyState
	case iam_model.LabelPolicySearchKeyInstanceID:
		return LabelPolicyKeyInstanceID
	case iam_model.LabelPolicySearchKeyOwnerRemoved:
		return LabelPolicyKeyOwnerRemoved

	default:
		return ""
	}
}
