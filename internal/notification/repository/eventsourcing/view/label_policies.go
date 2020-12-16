package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) LabelPolicyByAggregateID(aggregateID string, labelPolicyTableVar string) (*model.LabelPolicyView, error) {
	return view.GetLabelPolicyByAggregateID(v.Db, labelPolicyTableVar, aggregateID)
}
