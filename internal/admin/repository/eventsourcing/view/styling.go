package view

import (
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
)

const (
	stylingTable = "adminapi.styling2"
)

func (v *View) StylingByAggregateIDAndState(aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	return view.GetStylingByAggregateIDAndState(v.Db, stylingTable, aggregateID, instanceID, state)
}
