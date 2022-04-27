package view

import (
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
)

func (v *View) StylingByAggregateIDAndState(aggregateID, labelPolicyTableVar string, state int32) (*model.LabelPolicyView, error) {
	return view.GetStylingByAggregateIDAndState(v.Db, labelPolicyTableVar, aggregateID, state)
}
