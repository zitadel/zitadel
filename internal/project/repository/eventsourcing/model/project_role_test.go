package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"testing"
)

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
