package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddPasswordAgePolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordAgePolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add password age policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.PasswordAgePolicy{MaxAgeDays: 10},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordAgePolicy: &iam_es_model.PasswordAgePolicy{MaxAgeDays: 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddPasswordAgePolicyEvent(tt.args.event)
			if tt.result.PasswordAgePolicy.MaxAgeDays != tt.args.org.PasswordAgePolicy.MaxAgeDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordAgePolicy.MaxAgeDays, tt.args.org.PasswordAgePolicy.MaxAgeDays)
			}
		})
	}
}

func TestAppendChangePasswordAgePolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordAgePolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change password age policy event",
			args: args{
				org: &Org{PasswordAgePolicy: &iam_es_model.PasswordAgePolicy{
					MaxAgeDays: 10,
				}},
				policy: &iam_es_model.PasswordAgePolicy{MaxAgeDays: 5, ExpireWarnDays: 10},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordAgePolicy: &iam_es_model.PasswordAgePolicy{
				MaxAgeDays:     5,
				ExpireWarnDays: 10,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangePasswordAgePolicyEvent(tt.args.event)
			if tt.result.PasswordAgePolicy.MaxAgeDays != tt.args.org.PasswordAgePolicy.MaxAgeDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordAgePolicy.MaxAgeDays, tt.args.org.PasswordAgePolicy.MaxAgeDays)
			}
			if tt.result.PasswordAgePolicy.ExpireWarnDays != tt.args.org.PasswordAgePolicy.ExpireWarnDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordAgePolicy.ExpireWarnDays, tt.args.org.PasswordAgePolicy.ExpireWarnDays)
			}
		})
	}
}
