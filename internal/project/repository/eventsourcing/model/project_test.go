package model

import (
	"encoding/json"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestProjectFromEvents(t *testing.T) {
	type args struct {
		event   []eventstore.Event
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
				event: []eventstore.Event{
					&es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectAddedType},
				},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateActive), Name: "ProjectName"},
		},
		{
			name: "project from events, nil project",
			args: args{
				event: []eventstore.Event{
					&es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectAddedType},
				},
				project: nil,
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.project != nil {
				data, _ := json.Marshal(tt.args.project)
				tt.args.event[0].(*es_models.Event).Data = data
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
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectAddedType},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateActive), Name: "ProjectName"},
		},
		{
			name: "append change event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectChangedType},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateActive), Name: "ProjectName"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectDeactivatedType},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateInactive)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.ProjectReactivatedType},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.ProjectStateActive)},
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
			if result.ObjectRoot.AggregateID != tt.result.ObjectRoot.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.result.ObjectRoot.AggregateID, result.ObjectRoot.AggregateID)
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
			result: &Project{State: int32(model.ProjectStateInactive)},
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
			result: &Project{State: int32(model.ProjectStateActive)},
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
