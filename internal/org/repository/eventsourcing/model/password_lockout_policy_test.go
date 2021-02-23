package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddPasswordLockoutPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordLockoutPolicy
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
				policy: &iam_es_model.PasswordLockoutPolicy{MaxAttempts: 10},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordLockoutPolicy: &iam_es_model.PasswordLockoutPolicy{MaxAttempts: 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddPasswordLockoutPolicyEvent(tt.args.event)
			if tt.result.PasswordLockoutPolicy.MaxAttempts != tt.args.org.PasswordLockoutPolicy.MaxAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordLockoutPolicy.MaxAttempts, tt.args.org.PasswordLockoutPolicy.MaxAttempts)
			}
		})
	}
}

func TestAppendChangePasswordLockoutPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordLockoutPolicy
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
				org: &Org{PasswordLockoutPolicy: &iam_es_model.PasswordLockoutPolicy{
					MaxAttempts: 10,
				}},
				policy: &iam_es_model.PasswordLockoutPolicy{MaxAttempts: 5, ShowLockOutFailures: true},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordLockoutPolicy: &iam_es_model.PasswordLockoutPolicy{
				MaxAttempts:         5,
				ShowLockOutFailures: true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangePasswordLockoutPolicyEvent(tt.args.event)
			if tt.result.PasswordLockoutPolicy.MaxAttempts != tt.args.org.PasswordLockoutPolicy.MaxAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordLockoutPolicy.MaxAttempts, tt.args.org.PasswordLockoutPolicy.MaxAttempts)
			}
			if tt.result.PasswordLockoutPolicy.ShowLockOutFailures != tt.args.org.PasswordLockoutPolicy.ShowLockOutFailures {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordLockoutPolicy.ShowLockOutFailures, tt.args.org.PasswordLockoutPolicy.ShowLockOutFailures)
			}
		})
	}
}
