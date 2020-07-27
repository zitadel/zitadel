package eventsourcing

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	repo_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstoreComplexity(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *OrgEventstore {
	return &OrgEventstore{
		Eventstore: mockEs,
	}
}

func GetMockChangesOrgOK(ctrl *gomock.Controller) *OrgEventstore {
	org := model.Org{
		Name: "MusterOrg",
	}
	data, err := json.Marshal(org)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: repo_model.OrgAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesOrgNoEvents(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}
