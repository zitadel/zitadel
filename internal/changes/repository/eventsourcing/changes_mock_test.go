package eventsourcing

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/changes/model"
	chg_type "github.com/caos/zitadel/internal/changes/types"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstoreComplexity(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *ChangesEventstore {
	return &ChangesEventstore{
		Eventstore: mockEs,
	}
}

func GetMockChangesUserOK(ctrl *gomock.Controller) *ChangesEventstore {
	user := chg_type.User{
		FirstName:    "Hans",
		LastName:     "Muster",
		EMailAddress: "a@b.ch",
		Phone:        "+41 12 345 67 89",
		Language:     "D",
		UserName:     "HansMuster",
	}
	data, err := json.Marshal(user)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, AggregateType: model.User, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesUserNoEvents(ctrl *gomock.Controller) *ChangesEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesProjectOK(ctrl *gomock.Controller) *ChangesEventstore {
	project := chg_type.Project{
		Name: "MusterProject",
	}
	data, err := json.Marshal(project)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDProject", Sequence: 1, AggregateType: model.User, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesApplicationOK(ctrl *gomock.Controller) *ChangesEventstore {
	app := chg_type.App{
		Name:     "MusterApp",
		AppId:    "AppId",
		AppType:  3,
		ClientId: "MyClient",
	}
	data, err := json.Marshal(app)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: model.User, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesOrgOK(ctrl *gomock.Controller) *ChangesEventstore {
	org := chg_type.Org{
		Name:   "MusterOrg",
		Domain: "myDomain",
		UserId: "myUserId",
	}
	data, err := json.Marshal(org)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: model.User, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}
