package org_iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
)

const (
	AggregateType = "iam"
)

type OrgIAMPolicyWriteModel struct {
	eventstore.WriteModel
	Policy org_iam.OrgIAMPolicyWriteModel

	iamID string
}

func NewOrgIAMPolicyWriteModel(iamID string) *OrgIAMPolicyWriteModel {
	return &OrgIAMPolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *OrgIAMPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			wm.Policy.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *OrgIAMPolicyChangedEvent:
			wm.Policy.AppendEvents(&e.OrgIAMPolicyChangedEvent)
		}
	}
}

func (wm *OrgIAMPolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgIAMPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
