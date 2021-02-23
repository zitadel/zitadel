package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	"testing"
)

func TestAppendAddGrantEvent(t *testing.T) {
	type args struct {
		project *Project
		role    *ProjectGrant
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add grant event",
			args: args{
				project: &Project{},
				role:    &ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}},
				event:   &es_models.Event{},
			},
			result: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.role != nil {
				data, _ := json.Marshal(tt.args.role)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddGrantEvent(tt.args.event)
			if len(tt.args.project.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.project.Grants))
			}
			if tt.args.project.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.project.Grants[0])
			}
		})
	}
}

func TestAppendChangeGrantEvent(t *testing.T) {
	type args struct {
		project *Project
		grant   *ProjectGrant
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change grant event",
			args: args{
				project: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
				grant:   &ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"KeyChanged"}},
				event:   &es_models.Event{},
			},
			result: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"KeyChanged"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeGrantEvent(tt.args.event)
			if len(tt.args.project.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.project.Grants))
			}
			if tt.args.project.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.project.Grants[0])
			}
		})
	}
}

func TestAppendRemoveGrantEvent(t *testing.T) {
	type args struct {
		project *Project
		grant   *ProjectGrant
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
				project: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
				grant:   &ProjectGrant{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
			},
			result: &Project{Grants: []*ProjectGrant{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.project.appendRemoveGrantEvent(tt.args.event)
			if len(tt.args.project.Grants) != 0 {
				t.Errorf("got wrong result should have no grant actual: %v ", len(tt.args.project.Grants))
			}
		})
	}
}

func TestAppendGrantStateEvent(t *testing.T) {
	type args struct {
		project *Project
		grant   *ProjectGrantID
		event   *es_models.Event
		state   model.ProjectGrantState
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append deactivate grant event",
			args: args{
				project: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
				grant:   &ProjectGrantID{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
				state:   model.ProjectGrantStateInactive,
			},
			result: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}, State: int32(model.ProjectGrantStateInactive)}}},
		},
		{
			name: "append reactivate grant event",
			args: args{
				project: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
				grant:   &ProjectGrantID{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
				state:   model.ProjectGrantStateActive,
			},
			result: &Project{Grants: []*ProjectGrant{&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}, State: int32(model.ProjectGrantStateActive)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.project.appendGrantStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.project.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.project.Grants))
			}
			if tt.args.project.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.project.Grants[0])
			}
		})
	}
}
