package eventsourcing

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

func TestProjectFromEvents(t *testing.T) {
	type args struct {
		event   []*es_models.Event
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "project from events, ok",
			args: args{
				event: []*es_models.Event{
					&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_ACTIVE), Name: "ProjectName"},
		},
		{
			name: "project from events, nil project",
			args: args{
				event: []*es_models.Event{
					&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				},
				project: nil,
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_ACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.project != nil {
				data, _ := json.Marshal(tt.args.project)
				tt.args.event[0].Data = data
			}
			result, _ := ProjectFromEvents(tt.args.project, tt.args.event...)
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
		})
	}
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event   *es_models.Event
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append added event",
			args: args{
				event:   &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_ACTIVE), Name: "ProjectName"},
		},
		{
			name: "append change event",
			args: args{
				event:   &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectChanged},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_ACTIVE), Name: "ProjectName"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectDeactivated},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_INACTIVE)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectReactivated},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.PROJECTSTATE_ACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.project != nil {
				data, _ := json.Marshal(tt.args.project)
				tt.args.event.Data = data
			}
			result := new(Project)
			result.AppendEvent(tt.args.event)
			if result.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, result.State)
			}
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
			if result.ObjectRoot.ID != tt.result.ObjectRoot.ID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.result.ObjectRoot.ID, result.ObjectRoot.ID)
			}
		})
	}
}

func TestAppendDeactivatedEvent(t *testing.T) {
	type args struct {
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append reactivate event",
			args: args{
				project: &Project{},
			},
			result: &Project{State: int32(model.PROJECTSTATE_INACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.project.appendDeactivatedEvent()
			if tt.args.project.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.project)
			}
		})
	}
}

func TestAppendReactivatedEvent(t *testing.T) {
	type args struct {
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append reactivate event",
			args: args{
				project: &Project{},
			},
			result: &Project{State: int32(model.PROJECTSTATE_ACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.project.appendReactivatedEvent()
			if tt.args.project.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.project)
			}
		})
	}
}

func TestChanges(t *testing.T) {
	type args struct {
		existing *Project
		new      *Project
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project name changes",
			args: args{
				existing: &Project{Name: "Name"},
				new:      &Project{Name: "NameChanged"},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &Project{Name: "Name"},
				new:      &Project{Name: "Name"},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendAddMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectMember
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add member event",
			args: args{
				project: &Project{},
				member:  &ProjectMember{UserID: "UserID", Roles: []string{"Role"}},
				event:   &es_models.Event{},
			},
			result: &Project{Members: []*ProjectMember{&ProjectMember{UserID: "UserID", Roles: []string{"Role"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddMemberEvent(tt.args.event)
			if len(tt.args.project.Members) != 1 {
				t.Errorf("got wrong result should have one member actual: %v ", len(tt.args.project.Members))
			}
			if tt.args.project.Members[0] == tt.result.Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Members[0], tt.args.project.Members[0])
			}
		})
	}
}

func TestAppendChangeMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectMember
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change member event",
			args: args{
				project: &Project{Members: []*ProjectMember{&ProjectMember{UserID: "UserID", Roles: []string{"Role"}}}},
				member:  &ProjectMember{UserID: "UserID", Roles: []string{"ChangedRole"}},
				event:   &es_models.Event{},
			},
			result: &Project{Members: []*ProjectMember{&ProjectMember{UserID: "UserID", Roles: []string{"ChangedRole"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeMemberEvent(tt.args.event)
			if len(tt.args.project.Members) != 1 {
				t.Errorf("got wrong result should have one member actual: %v ", len(tt.args.project.Members))
			}
			if tt.args.project.Members[0] == tt.result.Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Members[0], tt.args.project.Members[0])
			}
		})
	}
}

func TestAppendRemoveMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectMember
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append remove member event",
			args: args{
				project: &Project{Members: []*ProjectMember{&ProjectMember{UserID: "UserID", Roles: []string{"Role"}}}},
				member:  &ProjectMember{UserID: "UserID"},
				event:   &es_models.Event{},
			},
			result: &Project{Members: []*ProjectMember{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendRemoveMemberEvent(tt.args.event)
			if len(tt.args.project.Members) != 0 {
				t.Errorf("got wrong result should have no member actual: %v ", len(tt.args.project.Members))
			}
		})
	}
}

func TestAppendAddRoleEvent(t *testing.T) {
	type args struct {
		project *Project
		role    *ProjectRole
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add role event",
			args: args{
				project: &Project{},
				role:    &ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"},
				event:   &es_models.Event{},
			},
			result: &Project{Roles: []*ProjectRole{&ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.role != nil {
				data, _ := json.Marshal(tt.args.role)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddRoleEvent(tt.args.event)
			if len(tt.args.project.Roles) != 1 {
				t.Errorf("got wrong result should have one role actual: %v ", len(tt.args.project.Roles))
			}
			if tt.args.project.Roles[0] == tt.result.Roles[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Roles[0], tt.args.project.Roles[0])
			}
		})
	}
}

func TestAppendChangeRoleEvent(t *testing.T) {
	type args struct {
		project *Project
		role    *ProjectRole
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change role event",
			args: args{
				project: &Project{Roles: []*ProjectRole{&ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"}}},
				role:    &ProjectRole{Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"},
				event:   &es_models.Event{},
			},
			result: &Project{Roles: []*ProjectRole{&ProjectRole{Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.role != nil {
				data, _ := json.Marshal(tt.args.role)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeRoleEvent(tt.args.event)
			if len(tt.args.project.Roles) != 1 {
				t.Errorf("got wrong result should have one role actual: %v ", len(tt.args.project.Roles))
			}
			if tt.args.project.Roles[0] == tt.result.Roles[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Roles[0], tt.args.project.Roles[0])
			}
		})
	}
}

func TestAppendRemoveRoleEvent(t *testing.T) {
	type args struct {
		project *Project
		role    *ProjectRole
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append remove role event",
			args: args{
				project: &Project{Roles: []*ProjectRole{&ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"}}},
				role:    &ProjectRole{Key: "Key"},
				event:   &es_models.Event{},
			},
			result: &Project{Roles: []*ProjectRole{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.role != nil {
				data, _ := json.Marshal(tt.args.role)
				tt.args.event.Data = data
			}
			tt.args.project.appendRemoveRoleEvent(tt.args.event)
			if len(tt.args.project.Roles) != 0 {
				t.Errorf("got wrong result should have no role actual: %v ", len(tt.args.project.Roles))
			}
		})
	}
}

func TestAppendAddAppEvent(t *testing.T) {
	type args struct {
		project *Project
		app     *Application
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add application event",
			args: args{
				project: &Project{},
				app:     &Application{Name: "Application"},
				event:   &es_models.Event{},
			},
			result: &Project{Applications: []*Application{&Application{Name: "Application"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.app != nil {
				data, _ := json.Marshal(tt.args.app)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddAppEvent(tt.args.event)
			if len(tt.args.project.Applications) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.Applications))
			}
			if tt.args.project.Applications[0] == tt.result.Applications[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Applications[0], tt.args.project.Applications[0])
			}
		})
	}
}

func TestAppendChangeAppEvent(t *testing.T) {
	type args struct {
		project *Project
		app     *Application
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change application event",
			args: args{
				project: &Project{Applications: []*Application{&Application{Name: "Application"}}},
				app:     &Application{Name: "Application Change"},
				event:   &es_models.Event{},
			},
			result: &Project{Applications: []*Application{&Application{Name: "Application Change"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.app != nil {
				data, _ := json.Marshal(tt.args.app)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeAppEvent(tt.args.event)
			if len(tt.args.project.Applications) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.Applications))
			}
			if tt.args.project.Applications[0] == tt.result.Applications[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Applications[0], tt.args.project.Applications[0])
			}
		})
	}
}

func TestAppendRemoveAppEvent(t *testing.T) {
	type args struct {
		project *Project
		app     *Application
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append remove application event",
			args: args{
				project: &Project{Applications: []*Application{&Application{AppID: "AppID", Name: "Application"}}},
				app:     &Application{AppID: "AppID", Name: "Application"},
				event:   &es_models.Event{},
			},
			result: &Project{Applications: []*Application{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.app != nil {
				data, _ := json.Marshal(tt.args.app)
				tt.args.event.Data = data
			}
			tt.args.project.appendRemoveAppEvent(tt.args.event)
			if len(tt.args.project.Applications) != 0 {
				t.Errorf("got wrong result should have no apps actual: %v ", len(tt.args.project.Applications))
			}
		})
	}
}

func TestAppendAppStateEvent(t *testing.T) {
	type args struct {
		project *Project
		app     *ApplicationID
		event   *es_models.Event
		state   model.AppState
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append deactivate application event",
			args: args{
				project: &Project{Applications: []*Application{&Application{AppID: "AppID", Name: "Application", State: model.AppStateToInt(model.APPSTATE_ACTIVE)}}},
				app:     &ApplicationID{AppID: "AppID"},
				event:   &es_models.Event{},
				state:   model.APPSTATE_INACTIVE,
			},
			result: &Project{Applications: []*Application{&Application{AppID: "AppID", Name: "Application", State: model.AppStateToInt(model.APPSTATE_INACTIVE)}}},
		},
		{
			name: "append reactivate application event",
			args: args{
				project: &Project{Applications: []*Application{&Application{AppID: "AppID", Name: "Application", State: model.AppStateToInt(model.APPSTATE_INACTIVE)}}},
				app:     &ApplicationID{AppID: "AppID"},
				event:   &es_models.Event{},
				state:   model.APPSTATE_ACTIVE,
			},
			result: &Project{Applications: []*Application{&Application{AppID: "AppID", Name: "Application", State: model.AppStateToInt(model.APPSTATE_ACTIVE)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.app != nil {
				data, _ := json.Marshal(tt.args.app)
				tt.args.event.Data = data
			}
			tt.args.project.appendAppStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.project.Applications) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.Applications))
			}
			if tt.args.project.Applications[0] == tt.result.Applications[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Applications[0], tt.args.project.Applications[0])
			}
		})
	}
}
