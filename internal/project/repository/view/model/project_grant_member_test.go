package model

import (
	"encoding/json"
	"reflect"
	"testing"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	es_model "github.com/zitadel/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func mockProjectGrantMemberData(member *es_model.ProjectGrantMember) []byte {
	data, _ := json.Marshal(member)
	return data
}

func TestGrantedProjectMemberAppendEvent(t *testing.T) {
	type args struct {
		event  *es_models.Event
		member *ProjectGrantMemberView
	}
	tests := []struct {
		name   string
		args   args
		result *ProjectGrantMemberView
	}{
		{
			name: "append added member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantMemberAddedType, ResourceOwner: "OrgID", Data: mockProjectGrantMemberData(&es_model.ProjectGrantMember{GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}})},
				member: &ProjectGrantMemberView{},
			},
			result: &ProjectGrantMemberView{ProjectID: "AggregateID", UserID: "UserID", GrantID: "ProjectGrantID", Roles: []string{"Role"}},
		},
		{
			name: "append changed member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantMemberAddedType, ResourceOwner: "OrgID", Data: mockProjectGrantMemberData(&es_model.ProjectGrantMember{GrantID: "ProjectGrantID", Roles: []string{"RoleChanged"}})},
				member: &ProjectGrantMemberView{ProjectID: "AggregateID", UserID: "UserID", GrantID: "ProjectGrantID", Roles: []string{"Role"}},
			},
			result: &ProjectGrantMemberView{ProjectID: "AggregateID", UserID: "UserID", GrantID: "ProjectGrantID", Roles: []string{"RoleChanged"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.member.AppendEvent(tt.args.event)
			if tt.args.member.ProjectID != tt.result.ProjectID {
				t.Errorf("got wrong result projectID: expected: %v, actual: %v ", tt.result.ProjectID, tt.args.member.ProjectID)
			}
			if tt.args.member.UserID != tt.result.UserID {
				t.Errorf("got wrong result userID: expected: %v, actual: %v ", tt.result.UserID, tt.args.member.UserID)
			}
			if tt.args.member.GrantID != tt.result.GrantID {
				t.Errorf("got wrong result ProjectGrantID: expected: %v, actual: %v ", tt.result.GrantID, tt.args.member.GrantID)
			}
			if !reflect.DeepEqual(tt.args.member.Roles, tt.result.Roles) {
				t.Errorf("got wrong result Roles: expected: %v, actual: %v ", tt.result.Roles, tt.args.member.Roles)
			}
		})
	}
}
