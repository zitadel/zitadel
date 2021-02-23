package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"testing"
)

func mockProjectRoleData(member *es_model.ProjectRole) []byte {
	data, _ := json.Marshal(member)
	return data
}

func TestProjectRoleAppendEvent(t *testing.T) {
	type args struct {
		event  *es_models.Event
		member *ProjectRoleView
	}
	tests := []struct {
		name   string
		args   args
		result *ProjectRoleView
	}{
		{
			name: "append added member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectRoleAdded, ResourceOwner: "OrgID", Data: mockProjectRoleData(&es_model.ProjectRole{Key: "Key", DisplayName: "DisplayName", Group: "Group"})},
				member: &ProjectRoleView{},
			},
			result: &ProjectRoleView{OrgID: "OrgID", ResourceOwner: "OrgID", ProjectID: "AggregateID", Key: "Key", DisplayName: "DisplayName", Group: "Group"},
		},
		{
			name: "append added member event",
			args: args{
				event:  &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectRoleAdded, ResourceOwner: "OrgID", Data: mockProjectRoleData(&es_model.ProjectRole{Key: "Key", DisplayName: "DisplayNameChanged", Group: "GroupChanged"})},
				member: &ProjectRoleView{OrgID: "OrgID", ResourceOwner: "OrgID", ProjectID: "AggregateID", Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			result: &ProjectRoleView{OrgID: "OrgID", ResourceOwner: "OrgID", ProjectID: "AggregateID", Key: "Key", DisplayName: "DisplayNameChanged", Group: "GroupChanged"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.member.AppendEvent(tt.args.event)
			if tt.args.member.ProjectID != tt.result.ProjectID {
				t.Errorf("got wrong result projectID: expected: %v, actual: %v ", tt.result.ProjectID, tt.args.member.ProjectID)
			}
			if tt.args.member.OrgID != tt.result.OrgID {
				t.Errorf("got wrong result orgID: expected: %v, actual: %v ", tt.result.OrgID, tt.args.member.OrgID)
			}
			if tt.args.member.ResourceOwner != tt.result.ResourceOwner {
				t.Errorf("got wrong result ResourceOwner: expected: %v, actual: %v ", tt.result.ResourceOwner, tt.args.member.ResourceOwner)
			}
			if tt.args.member.Key != tt.result.Key {
				t.Errorf("got wrong result Key: expected: %v, actual: %v ", tt.result.Key, tt.args.member.Key)
			}
			if tt.args.member.DisplayName != tt.result.DisplayName {
				t.Errorf("got wrong result DisplayName: expected: %v, actual: %v ", tt.result.DisplayName, tt.args.member.DisplayName)
			}
			if tt.args.member.Group != tt.result.Group {
				t.Errorf("got wrong result Group: expected: %v, actual: %v ", tt.result.Group, tt.args.member.Group)
			}
		})
	}
}
