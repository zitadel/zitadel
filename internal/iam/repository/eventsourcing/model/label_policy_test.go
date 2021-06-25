package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

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
				policy: &LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultLabelPolicy: &LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"}},
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
			if tt.result.DefaultLabelPolicy.BackgroundColor != tt.args.iam.DefaultLabelPolicy.BackgroundColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.BackgroundColor, tt.args.iam.DefaultLabelPolicy.BackgroundColor)
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
					PrimaryColor: "000001", BackgroundColor: "FFFFF0",
				}},
				policy: &LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultLabelPolicy: &LabelPolicy{
				PrimaryColor: "000000", BackgroundColor: "FFFFFF",
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
			if tt.result.DefaultLabelPolicy.BackgroundColor != tt.args.iam.DefaultLabelPolicy.BackgroundColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLabelPolicy.BackgroundColor, tt.args.iam.DefaultLabelPolicy.BackgroundColor)
			}
		})
	}
}
