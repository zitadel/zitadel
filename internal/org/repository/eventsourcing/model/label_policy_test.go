package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestAppendAddLabelPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LabelPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add label policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"},
				event:  &es_models.Event{},
			},
			result: &Org{LabelPolicy: &iam_es_model.LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddLabelPolicyEvent(tt.args.event)
			if tt.result.LabelPolicy.PrimaryColor != tt.args.org.LabelPolicy.PrimaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LabelPolicy.PrimaryColor, tt.args.org.LabelPolicy.PrimaryColor)
			}
			if tt.result.LabelPolicy.BackgroundColor != tt.args.org.LabelPolicy.BackgroundColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LabelPolicy.BackgroundColor, tt.args.org.LabelPolicy.BackgroundColor)
			}
		})
	}
}

func TestAppendChangeLabelPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LabelPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change label policy event",
			args: args{
				org: &Org{LabelPolicy: &iam_es_model.LabelPolicy{
					BackgroundColor: "FFFFF0",
					PrimaryColor:    "000001",
				}},
				policy: &iam_es_model.LabelPolicy{PrimaryColor: "000000", BackgroundColor: "FFFFFF"},
				event:  &es_models.Event{},
			},
			result: &Org{LabelPolicy: &iam_es_model.LabelPolicy{
				BackgroundColor: "FFFFFF",
				PrimaryColor:    "000000",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeLabelPolicyEvent(tt.args.event)
			if tt.result.LabelPolicy.PrimaryColor != tt.args.org.LabelPolicy.PrimaryColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LabelPolicy.PrimaryColor, tt.args.org.LabelPolicy.PrimaryColor)
			}
			if tt.result.LabelPolicy.BackgroundColor != tt.args.org.LabelPolicy.BackgroundColor {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LabelPolicy.BackgroundColor, tt.args.org.LabelPolicy.BackgroundColor)
			}
		})
	}
}
