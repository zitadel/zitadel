package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"testing"
)

func mockIamData(iam *Iam) []byte {
	data, _ := json.Marshal(iam)
	return data
}

func TestProjectRoleAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		iam   *Iam
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append set up start event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IamSetupStarted, ResourceOwner: "OrgID"},
				iam:   &Iam{},
			},
			result: &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true},
		},
		{
			name: "append set up done event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IamSetupDone, ResourceOwner: "OrgID"},
				iam:   &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true},
			},
			result: &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true, SetUpDone: true},
		},
		{
			name: "append globalorg event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: GlobalOrgSet, ResourceOwner: "OrgID", Data: mockIamData(&Iam{GlobalOrgID: "GlobalOrg"})},
				iam:   &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true},
			},
			result: &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true, GlobalOrgID: "GlobalOrg"},
		},
		{
			name: "append iamproject event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: IamProjectSet, ResourceOwner: "OrgID", Data: mockIamData(&Iam{IamProjectID: "IamProject"})},
				iam:   &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true},
			},
			result: &Iam{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SetUpStarted: true, IamProjectID: "IamProject"},
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
			if tt.args.iam.IamProjectID != tt.result.IamProjectID {
				t.Errorf("got wrong result IamProjectID: expected: %v, actual: %v ", tt.result.IamProjectID, tt.args.iam.IamProjectID)
			}
		})
	}
}
