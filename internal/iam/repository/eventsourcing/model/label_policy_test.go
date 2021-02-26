package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestLabelPolicyChanges(t *testing.T) {
	type args struct {
		existing *LabelPolicy
		new      *LabelPolicy
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "labelpolicy all attributes change",
			args: args{
				existing: &LabelPolicy{PrimaryColor: "000001", SecondaryColor: "FFFFFA"},
				new:      &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"},
				new:      &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendAddLabelPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *LabelPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add label policy event",
			args: args{
				iam:    new(IAM),
				policy: &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultLabelPolicy: &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddLabelPolicyEvent(tt.args.event)
			if tt.result.DefaultLabelPolicy.PrimaryColor != tt.args.iam.DefaultLabelPolicy.PrimaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.PrimaryColor, tt.args.iam.DefaultLabelPolicy.PrimaryColor)
			}
			if tt.result.DefaultLabelPolicy.SecondaryColor != tt.args.iam.DefaultLabelPolicy.SecondaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.SecondaryColor, tt.args.iam.DefaultLabelPolicy.SecondaryColor)
			}
		})
	}
}

func TestAppendChangeLabelPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *LabelPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change label policy event",
			args: args{
				iam: &IAM{DefaultLabelPolicy: &LabelPolicy{
					PrimaryColor: "000001", SecondaryColor: "FFFFF0",
				}},
				policy: &LabelPolicy{PrimaryColor: "000000", SecondaryColor: "FFFFFF"},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultLabelPolicy: &LabelPolicy{
				PrimaryColor: "000000", SecondaryColor: "FFFFFF",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeLabelPolicyEvent(tt.args.event)
			if tt.result.DefaultLabelPolicy.PrimaryColor != tt.args.iam.DefaultLabelPolicy.PrimaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.PrimaryColor, tt.args.iam.DefaultLabelPolicy.PrimaryColor)
			}
			if tt.result.DefaultLabelPolicy.SecondaryColor != tt.args.iam.DefaultLabelPolicy.SecondaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.SecondaryColor, tt.args.iam.DefaultLabelPolicy.SecondaryColor)
			}
		})
	}
}
