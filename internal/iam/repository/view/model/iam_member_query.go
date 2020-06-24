package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view"
)

type IamMemberSearchRequest iam_model.IamMemberSearchRequest
type IamMemberSearchQuery iam_model.IamMemberSearchQuery
type IamMemberSearchKey iam_model.IamMemberSearchKey

func (req IamMemberSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IamMemberSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IamMemberSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == iam_model.IamMemberSearchKeyUnspecified {
		return nil
	}
	return IamMemberSearchKey(req.SortingColumn)
}

func (req IamMemberSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IamMemberSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IamMemberSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IamMemberSearchQuery) GetKey() view.ColumnKey {
	return IamMemberSearchKey(req.Key)
}

func (req IamMemberSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req IamMemberSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IamMemberSearchKey) ToColumnName() string {
	switch iam_model.IamMemberSearchKey(key) {
	case iam_model.IamMemberSearchKeyEmail:
		return IamMemberKeyEmail
	case iam_model.IamMemberSearchKeyFirstName:
		return IamMemberKeyFirstName
	case iam_model.IamMemberSearchKeyLastName:
		return IamMemberKeyLastName
	case iam_model.IamMemberSearchKeyUserName:
		return IamMemberKeyUserName
	case iam_model.IamMemberSearchKeyUserID:
		return IamMemberKeyUserID
	case iam_model.IamMemberSearchKeyIamID:
		return IamMemberKeyIamID
	default:
		return ""
	}
}
