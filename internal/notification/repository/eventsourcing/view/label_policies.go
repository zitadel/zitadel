package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) LabelPolicyByAggregateIDAndState(aggregateID, labelPolicyTableVar string, state int32) (*model.LabelPolicyView, error) {
	return view.GetLabelPolicyByAggregateIDAndState(v.Db, labelPolicyTableVar, aggregateID, state)
}
