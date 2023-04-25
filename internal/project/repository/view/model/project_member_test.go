package model

import (
	"encoding/json"
	"reflect"
	"testing"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	es_model "github.com/zitadel/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func mockProjectMemberData(member *es_model.ProjectMember) []byte {
	data, _ := json.Marshal(member)
	return data
}

func TestProjectMemberAppendEvent(t *testing.T) {
	type args struct {
		event  *es_models.Event
		member *ProjectMemberView
	}
	tests := []struct {
		name   string
		args   args
		result *ProjectMemberView
	}{
		{
			name: "append added member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.MemberAddedType, ResourceOwner: "OrgID", Data: mockProjectMemberData(&es_model.ProjectMember{UserID: "UserID", Roles: []string{"Role"}})},
				member: &ProjectMemberView{},
			},
			result: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: []string{"Role"}},
		},
		{
			name: "append changed member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.MemberAddedType, ResourceOwner: "OrgID", Data: mockProjectMemberData(&es_model.ProjectMember{UserID: "UserID", Roles: []string{"RoleChanged"}})},
				member: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: []string{"Role"}},
			},
			result: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: []string{"RoleChanged"}},
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
			if !reflect.DeepEqual(tt.args.member.Roles, tt.result.Roles) {
				t.Errorf("got wrong result Roles: expected: %v, actual: %v ", tt.result.Roles, tt.args.member.Roles)
			}
		})
	}
}
