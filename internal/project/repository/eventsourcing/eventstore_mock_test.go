package eventsourcing

import (
	"encoding/json"
	mock_cache "github.com/caos/zitadel/internal/cache/mock"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/golang/mock/gomock"
)

func GetMockCache(ctrl *gomock.Controller) *ProjectCache {
	mockCache := mock_cache.NewMockCache(ctrl)
	mockCache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &ProjectCache{projectCache: mockCache}
}

func GetMockProjectByIDOK(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockProjectByIDNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockManipulateProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockManipulateInactiveProject(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 2, Type: model.ProjectDeactivated, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockManipulateProjectWithMember(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	memberData, _ := json.Marshal(ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockManipulateProjectWithRole(ctrl *gomock.Controller) *ProjectEventstore {
	data, _ := json.Marshal(Project{Name: "Name"})
	roleData, _ := json.Marshal(ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectRoleAdded, Data: roleData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockManipulateProjectNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}

func GetMockProjectMemberByIDsOK(ctrl *gomock.Controller) *ProjectEventstore {
	projectData, _ := json.Marshal(Project{Name: "Name"})
	memberData, _ := json.Marshal(ProjectMember{UserID: "UserID", Roles: []string{"Role"}})
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: projectData},
		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectMemberAdded, Data: memberData},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return &ProjectEventstore{Eventstore: mockEs, projectCache: GetMockCache(ctrl)}
}
