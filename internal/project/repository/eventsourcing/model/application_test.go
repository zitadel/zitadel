package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
)

func TestApplicationChanges(t *testing.T) {
	type args struct {
		existingProject *Application
		newProject      *Application
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
			name: "application name changes",
			args: args{
				existingProject: &Application{AppID: "AppID", Name: "Name"},
				newProject:      &Application{AppID: "AppID", Name: "NameChanged"},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existingProject: &Application{AppID: "AppID", Name: "Name"},
				newProject:      &Application{AppID: "AppID", Name: "Name"},
			},
			res: res{
				changesLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingProject.Changes(tt.args.newProject)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
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
			result: &Project{
				Applications: []*Application{
					{Name: "Application"},
				},
			},
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
				project: &Project{
					Applications: []*Application{
						{Name: "Application"},
					},
				},
				app:   &Application{Name: "Application Change"},
				event: &es_models.Event{},
			},
			result: &Project{
				Applications: []*Application{
					{Name: "Application Change"},
				},
			},
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
				project: &Project{
					Applications: []*Application{
						{AppID: "AppID", Name: "Application"},
					},
				},
				app:   &Application{AppID: "AppID", Name: "Application"},
				event: &es_models.Event{},
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
				project: &Project{
					Applications: []*Application{
						{AppID: "AppID", Name: "Application", State: int32(model.AppStateActive)},
					},
				},
				app:   &ApplicationID{AppID: "AppID"},
				event: &es_models.Event{},
				state: model.AppStateInactive,
			},
			result: &Project{
				Applications: []*Application{
					{AppID: "AppID", Name: "Application", State: int32(model.AppStateInactive)},
				},
			},
		},
		{
			name: "append reactivate application event",
			args: args{
				project: &Project{
					Applications: []*Application{
						{AppID: "AppID", Name: "Application", State: int32(model.AppStateInactive)},
					},
				},
				app:   &ApplicationID{AppID: "AppID"},
				event: &es_models.Event{},
				state: model.AppStateActive,
			},
			result: &Project{
				Applications: []*Application{
					{AppID: "AppID", Name: "Application", State: int32(model.AppStateActive)},
				},
			},
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
