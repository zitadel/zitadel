package view

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/iam/repository/view"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
)

const (
	stylingTyble = "adminapi.styling2"
)

func (v *View) StylingByAggregateIDAndState(aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	return view.GetStylingByAggregateIDAndState(v.Db, stylingTyble, aggregateID, instanceID, state)
}

func (v *View) PutStyling(policy *model.LabelPolicyView, event eventstore.Event) error {
	return view.PutStyling(v.Db, stylingTyble, policy)
}

func (v *View) DeleteInstanceStyling(event eventstore.Event) error {
	return view.DeleteInstanceStyling(v.Db, stylingTyble, event.Aggregate().InstanceID)
}

func (v *View) UpdateOrgOwnerRemovedStyling(event eventstore.Event) error {
	return view.UpdateOrgOwnerRemovedStyling(v.Db, stylingTyble, event.Aggregate().InstanceID, event.Aggregate().ID)
}
