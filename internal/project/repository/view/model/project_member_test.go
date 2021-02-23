package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/lib/pq"
	"reflect"
	"testing"
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
				event:  &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectMemberAdded, ResourceOwner: "OrgID", Data: mockProjectMemberData(&es_model.ProjectMember{UserID: "UserID", Roles: pq.StringArray{"Role"}})},
				member: &ProjectMemberView{},
			},
			result: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: pq.StringArray{"Role"}},
		},
		{
			name: "append changed member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectMemberAdded, ResourceOwner: "OrgID", Data: mockProjectMemberData(&es_model.ProjectMember{UserID: "UserID", Roles: pq.StringArray{"RoleChanged"}})},
				member: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: pq.StringArray{"Role"}},
			},
			result: &ProjectMemberView{ProjectID: "AggregateID", UserID: "UserID", Roles: pq.StringArray{"RoleChanged"}},
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
