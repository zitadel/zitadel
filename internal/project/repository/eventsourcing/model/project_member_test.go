package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

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
