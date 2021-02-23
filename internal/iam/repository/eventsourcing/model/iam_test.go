package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func mockIamData(iam *IAM) []byte {
	data, _ := json.Marshal(iam)
	return data
}

func TestProjectRoleAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		iam   *IAM
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append set up start event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IAMSetupStarted, ResourceOwner: "OrgID"},
				iam:   &IAM{},
			},
			result: &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1},
		},
		{
			name: "append set up done event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IAMSetupDone, ResourceOwner: "OrgID"},
				iam:   &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1},
			},
			result: &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1, SetUpDone: Step1},
		},
		{
			name: "append globalorg event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: GlobalOrgSet, ResourceOwner: "OrgID", Data: mockIamData(&IAM{GlobalOrgID: "GlobalOrg"})},
				iam:   &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1},
			},
			result: &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1, GlobalOrgID: "GlobalOrg"},
		},
		{
			name: "append iamproject event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IAMProjectSet, ResourceOwner: "OrgID", Data: mockIamData(&IAM{IAMProjectID: "IamProject"})},
				iam:   &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1},
			},
			result: &IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: Step1, IAMProjectID: "IamProject"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.iam.AppendEvent(tt.args.event)
			if tt.args.iam.AggregateID != tt.result.AggregateID {
				t.Errorf("got wrong result AggregateID: expected: %v, actual: %v ", tt.result.AggregateID, tt.args.iam.AggregateID)
			}
			if tt.args.iam.SetUpDone != tt.result.SetUpDone {
				t.Errorf("got wrong result SetUpDone: expected: %v, actual: %v ", tt.result.SetUpDone, tt.args.iam.SetUpDone)
			}
			if tt.args.iam.GlobalOrgID != tt.result.GlobalOrgID {
				t.Errorf("got wrong result GlobalOrgID: expected: %v, actual: %v ", tt.result.GlobalOrgID, tt.args.iam.GlobalOrgID)
			}
			if tt.args.iam.IAMProjectID != tt.result.IAMProjectID {
				t.Errorf("got wrong result IAMProjectID: expected: %v, actual: %v ", tt.result.IAMProjectID, tt.args.iam.IAMProjectID)
			}
		})
	}
}
