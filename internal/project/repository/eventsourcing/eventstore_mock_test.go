package eventsourcing

import (
	"encoding/json"

	"github.com/golang/mock/gomock"

	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	repo_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
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

func GetSonyFlacke() id.Generator {
	return id.SonyFlakeGenerator
}

func GetMockPwGenerator(ctrl *gomock.Controller) crypto.Generator {
	generator := crypto.NewMockGenerator(ctrl)
	generator.EXPECT().Length().Return(uint(10)).AnyTimes()
	generator.EXPECT().Runes().Return([]rune("abcdefghijklmnopqrstuvwxyz")).AnyTimes()
	generator.EXPECT().Alg().Return(crypto.NewBCrypt(10)).AnyTimes()
	return generator
}

func GetMockProjectByIDOK(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectByIDNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithPw(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs)
}

func GetMockManipulateInactiveProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 2, Type: model.ProjectDeactivated, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithMember(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	memberData, _ := json.Marshal(model.ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithRole(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	roleData, _ := json.Marshal(model.ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectRoleAdded, Data: roleData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithOIDCApp(ctrl *gomock.Controller, authMethod proj_model.OIDCAuthMethodType) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	appData, _ := json.Marshal(model.Application{AppID: "AppID", Name: "Name"})
	oidcData, _ := json.Marshal(model.OIDCConfig{
		AppID:          "AppID",
		ResponseTypes:  []int32{int32(proj_model.OIDCResponseTypeCode)},
		GrantTypes:     []int32{int32(proj_model.OIDCGrantTypeAuthorizationCode)},
		AuthMethodType: int32(authMethod),
	})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstoreWithPw(ctrl, mockEs)
}

func GetMockManipulateProjectWithSAMLApp(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	appData, _ := json.Marshal(model.Application{AppID: "AppID", Name: "Name"})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithGrant(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	grantData, _ := json.Marshal(model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: []string{"Key"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantAdded, Data: grantData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithGrantExistingRole(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	roleData, _ := json.Marshal(model.ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"})
	roleData2, _ := json.Marshal(model.ProjectRole{Key: "KeyChanged", DisplayName: "DisplayName", Group: "Group"})
	grantData, _ := json.Marshal(model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: []string{"Key"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectRoleAdded, Data: roleData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectRoleAdded, Data: roleData2},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantAdded, Data: grantData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectWithGrantMember(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(model.Project{Name: "Name"})
	grantData, _ := json.Marshal(model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: []string{"Key"}})
	memberData, _ := json.Marshal(model.ProjectGrantMember{GrantID: "GrantID", UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantAdded, Data: grantData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockManipulateProjectNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST")).AnyTimes()
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectMemberByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(model.Project{Name: "Name"})
	memberData, _ := json.Marshal(model.ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectAppsByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(model.Project{Name: "Name"})
	appData, _ := json.Marshal(model.Application{AppID: "AppID", Name: "Name"})
	oidcData, _ := json.Marshal(model.OIDCConfig{ClientID: "ClientID"})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.ApplicationAdded, Data: appData},
		&es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: model.OIDCConfigAdded, Data: oidcData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectGrantByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(model.Project{Name: "Name"})
	grantData, _ := json.Marshal(model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantID", RoleKeys: []string{"Key"}})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantAdded, Data: grantData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockProjectGrantMemberByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(model.Project{Name: "Name"})
	grantData, _ := json.Marshal(model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantID", RoleKeys: []string{"Key"}})
	memberData, _ := json.Marshal(model.ProjectGrantMember{GrantID: "GrantID", UserID: "UserID", Roles: []string{"Role"}})

	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantAdded, Data: grantData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectGrantMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesProjectOK(ctrl *gomock.Controller) *ProjectEventstore {
	project := model.Project{
		Name: "MusterProject",
	}
	data, err := json.Marshal(project)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDProject", Sequence: 1, AggregateType: repo_model.ProjectAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesProjectNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockedEventstoreComplexity(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *ProjectEventstore {
	return &ProjectEventstore{
		Eventstore: mockEs,
	}
}

func GetMockChangesApplicationOK(ctrl *gomock.Controller) *ProjectEventstore {
	app := model.Application{
		Name:  "MusterApp",
		AppID: "AppId",
		Type:  3,
	}
	data, err := json.Marshal(app)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Type: "project.application.added", Sequence: 1, AggregateType: repo_model.ProjectAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}

func GetMockChangesApplicationNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil).AnyTimes()
	return GetMockedEventstoreComplexity(ctrl, mockEs)
}
