package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"
)

const (
	AggregateType = "iam"
)

type LoginPolicyIDPProviderWriteModel struct {
	eventstore.WriteModel
	idpprovider.IDPProviderWriteModel

	idpConfigID string
	iamID       string

	IsRemoved bool
}

func NewLoginPolicyIDPProviderWriteModel(iamID, idpConfigID string) *LoginPolicyIDPProviderWriteModel {
	return &LoginPolicyIDPProviderWriteModel{
		iamID:       iamID,
		idpConfigID: idpConfigID,
	}
}

func (wm *LoginPolicyIDPProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyIDPProviderAddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IDPProviderWriteModel.AppendEvents(&e.IDPProviderAddedEvent)
		}
	}
}

func (wm *LoginPolicyIDPProviderWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *LoginPolicyIDPProviderAddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = false
		case *LoginPolicyIDPProviderRemovedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = true
		}
	}
	if err := wm.IDPProviderWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *LoginPolicyIDPProviderWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
