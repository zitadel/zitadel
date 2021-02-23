package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestPasswordLockoutPolicyChanges(t *testing.T) {
	type args struct {
		existing *PasswordLockoutPolicy
		new      *PasswordLockoutPolicy
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
			name: "lockout policy all attributes change",
			args: args{
				existing: &PasswordLockoutPolicy{MaxAttempts: 365, ShowLockOutFailures: true},
				new:      &PasswordLockoutPolicy{MaxAttempts: 730, ShowLockOutFailures: false},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &PasswordLockoutPolicy{MaxAttempts: 10, ShowLockOutFailures: true},
				new:      &PasswordLockoutPolicy{MaxAttempts: 10, ShowLockOutFailures: true},
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

func TestAppendAddPasswordLockoutPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordLockoutPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add password lockout policy event",
			args: args{
				iam:    new(IAM),
				policy: &PasswordLockoutPolicy{MaxAttempts: 10, ShowLockOutFailures: true},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultPasswordLockoutPolicy: &PasswordLockoutPolicy{MaxAttempts: 10, ShowLockOutFailures: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddPasswordLockoutPolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordLockoutPolicy.MaxAttempts != tt.args.iam.DefaultPasswordLockoutPolicy.MaxAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordLockoutPolicy.MaxAttempts, tt.args.iam.DefaultPasswordLockoutPolicy.MaxAttempts)
			}
			if tt.result.DefaultPasswordLockoutPolicy.ShowLockOutFailures != tt.args.iam.DefaultPasswordLockoutPolicy.ShowLockOutFailures {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordLockoutPolicy.ShowLockOutFailures, tt.args.iam.DefaultPasswordLockoutPolicy.ShowLockOutFailures)
			}
		})
	}
}

func TestAppendChangePasswordLockoutPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordLockoutPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change password lockout policy event",
			args: args{
				iam: &IAM{DefaultPasswordLockoutPolicy: &PasswordLockoutPolicy{
					MaxAttempts: 10,
				}},
				policy: &PasswordLockoutPolicy{MaxAttempts: 5},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultPasswordLockoutPolicy: &PasswordLockoutPolicy{
				MaxAttempts: 5,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangePasswordLockoutPolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordLockoutPolicy.MaxAttempts != tt.args.iam.DefaultPasswordLockoutPolicy.MaxAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordLockoutPolicy.MaxAttempts, tt.args.iam.DefaultPasswordLockoutPolicy.MaxAttempts)
			}
		})
	}
}
