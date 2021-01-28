package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"testing"
)

func TestAppendAddGrantMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectGrantMember
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add grant member",
			args: args{
				project: &Project{Grants: []*ProjectGrant{
					&ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "OrgID", RoleKeys: []string{"Key"}}}},
				member: &ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}},
				event:  &es_models.Event{},
			},
			result: &Project{
				Grants: []*ProjectGrant{
					&ProjectGrant{
						GrantID:      "ProjectGrantID",
						GrantedOrgID: "OrgID",
						RoleKeys:     []string{"Key"},
						Members:      []*ProjectGrantMember{&ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddGrantMemberEvent(tt.args.event)
			if len(tt.args.project.Grants[0].Members) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.project.Grants[0].Members))
			}
			if tt.args.project.Grants[0].Members[0] == tt.result.Grants[0].Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0].Members[0], tt.args.project.Grants[0].Members[0])
			}
		})
	}
}

func TestAppendChangeGrantMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectGrantMember
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change grant member",
			args: args{
				project: &Project{
					Grants: []*ProjectGrant{
						&ProjectGrant{
							GrantID:      "ProjectGrantID",
							GrantedOrgID: "OrgID",
							RoleKeys:     []string{"Key"},
							Members:      []*ProjectGrantMember{&ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}}}}},
				},
				member: &ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"RoleChanged"}},
				event:  &es_models.Event{},
			},
			result: &Project{
				Grants: []*ProjectGrant{
					&ProjectGrant{
						GrantID:      "ProjectGrantID",
						GrantedOrgID: "OrgID",
						RoleKeys:     []string{"Key"},
						Members:      []*ProjectGrantMember{&ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"RoleChanged"}}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeGrantMemberEvent(tt.args.event)
			if len(tt.args.project.Grants[0].Members) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.project.Grants[0].Members))
			}
			if tt.args.project.Grants[0].Members[0] == tt.result.Grants[0].Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0].Members[0], tt.args.project.Grants[0].Members[0])
			}
		})
	}
}

func TestAppendRemoveGrantMemberEvent(t *testing.T) {
	type args struct {
		project *Project
		member  *ProjectGrantMember
		event   *es_models.Event
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "append remove grant member",
			args: args{
				project: &Project{
					Grants: []*ProjectGrant{
						&ProjectGrant{
							GrantID:      "ProjectGrantID",
							GrantedOrgID: "OrgID",
							RoleKeys:     []string{"Key"},
							Members:      []*ProjectGrantMember{&ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}}}}},
				},
				member: &ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"RoleChanged"}},
				event:  &es_models.Event{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.project.appendRemoveGrantMemberEvent(tt.args.event)
			if len(tt.args.project.Grants[0].Members) != 0 {
				t.Errorf("got wrong result should have no members actual: %v ", len(tt.args.project.Grants[0].Members))
			}
		})
	}
}
