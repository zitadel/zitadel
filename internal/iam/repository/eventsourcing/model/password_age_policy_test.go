package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestPasswordAgePolicyChanges(t *testing.T) {
	type args struct {
		existing *PasswordAgePolicy
		new      *PasswordAgePolicy
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
			name: "age policy all attributes change",
			args: args{
				existing: &PasswordAgePolicy{MaxAgeDays: 365, ExpireWarnDays: 5},
				new:      &PasswordAgePolicy{MaxAgeDays: 730, ExpireWarnDays: 10},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10},
				new:      &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10},
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

func TestAppendAddPasswordAgePolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordAgePolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add password age policy event",
			args: args{
				iam:    new(IAM),
				policy: &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultPasswordAgePolicy: &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddPasswordAgePolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordAgePolicy.MaxAgeDays != tt.args.iam.DefaultPasswordAgePolicy.MaxAgeDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordAgePolicy.MaxAgeDays, tt.args.iam.DefaultPasswordAgePolicy.MaxAgeDays)
			}
			if tt.result.DefaultPasswordAgePolicy.ExpireWarnDays != tt.args.iam.DefaultPasswordAgePolicy.ExpireWarnDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordAgePolicy.ExpireWarnDays, tt.args.iam.DefaultPasswordAgePolicy.ExpireWarnDays)
			}
		})
	}
}

func TestAppendChangePasswordAgePolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordAgePolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change password age policy event",
			args: args{
				iam: &IAM{DefaultPasswordAgePolicy: &PasswordAgePolicy{
					MaxAgeDays: 10,
				}},
				policy: &PasswordAgePolicy{MaxAgeDays: 5},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultPasswordAgePolicy: &PasswordAgePolicy{
				MaxAgeDays: 5,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangePasswordAgePolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordAgePolicy.MaxAgeDays != tt.args.iam.DefaultPasswordAgePolicy.MaxAgeDays {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordAgePolicy.MaxAgeDays, tt.args.iam.DefaultPasswordAgePolicy.MaxAgeDays)
			}
		})
	}
}
