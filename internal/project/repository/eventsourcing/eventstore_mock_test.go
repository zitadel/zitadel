package eventsourcing

import (
	"encoding/json"
	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/golang/mock/gomock"
	"github.com/sony/sonyflake"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *ProjectEventstore {
	return &ProjectEventstore{
		Eventstore:   mockEs,
		projectCache: GetMockCache(ctrl),
		idGenerator:  GetSonyFlacke(),
	}
}

func GetMockedEventstoreWithPw(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *ProjectEventstore {
	return &ProjectEventstore{
		Eventstore:   mockEs,
		projectCache: GetMockCache(ctrl),
		idGenerator:  GetSonyFlacke(),
		pwGenerator:  GetMockPwGenerator(ctrl),
	}
}
func GetMockCache(ctrl *gomock.Controller) *ProjectCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &ProjectCache{projectCache: mockCache}
}

func GetSonyFlacke() *sonyflake.Sonyflake {
	return sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GetMockPwGenerator(ctrl *gomock.Controller) crypto.Generator {
	generator := crypto.NewMockGenerator(ctrl)
	generator.EXPECT().Length().Return(uint(10))
	generator.EXPECT().Runes().Return([]rune("abcdefghijklmnopqrstuvwxyz"))
	generator.EXPECT().Alg().Return(crypto.NewBCrypt(10))
	return generator
}

func GetMockProjectByIDOK(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectByIDNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithPw(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithPw(ctrl, mockEs)
}

func GetMockManipulateInactiveProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 2, Type: model.ProjectDeactivated, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithMember(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	memberData, _ := json.Marshal(ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithRole(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	roleData, _ := json.Marshal(ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectRoleAdded, Data: roleData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithOIDCApp(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	appData, _ := json.Marshal(Application{AppID: "AppID", Name: "Name"})
	oidcData, _ := json.Marshal(OIDCConfig{
		AppID:         "AppID",
		ResponseTypes: []int32{int32(model.OIDCRESPONSETYPE_CODE)},
		GrantTypes:    []int32{int32(model.OIDCGRANTTYPE_AUTHORIZATION_CODE)},
	})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstoreWithPw(ctrl, mockEs)
}

func GetMockManipulateProjectWithSAMLApp(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	appData, _ := json.Marshal(Application{AppID: "AppID", Name: "Name"})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectMemberByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(Project{Name: "Name"})
	memberData, _ := json.Marshal(ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectAppsByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(Project{Name: "Name"})
	appData, _ := json.Marshal(Application{AppID: "AppID", Name: "Name"})
	oidcData, _ := json.Marshal(OIDCConfig{ClientID: "ClientID"})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}
