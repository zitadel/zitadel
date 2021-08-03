package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestPasswordLockoutPolicyChanges(t *testing.T) {
	type args struct {
		existing *LockoutPolicy
		new      *LockoutPolicy
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
				existing: &LockoutPolicy{MaxPasswordAttempts: 365, ShowLockOutFailures: true},
				new:      &LockoutPolicy{MaxPasswordAttempts: 730, ShowLockOutFailures: false},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true},
				new:      &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true},
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
		policy *LockoutPolicy
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
				policy: &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultLockoutPolicy: &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddLockoutPolicyEvent(tt.args.event)
			if tt.result.DefaultLockoutPolicy.MaxPasswordAttempts != tt.args.iam.DefaultLockoutPolicy.MaxPasswordAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLockoutPolicy.MaxPasswordAttempts, tt.args.iam.DefaultLockoutPolicy.MaxPasswordAttempts)
			}
			if tt.result.DefaultLockoutPolicy.ShowLockOutFailures != tt.args.iam.DefaultLockoutPolicy.ShowLockOutFailures {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLockoutPolicy.ShowLockOutFailures, tt.args.iam.DefaultLockoutPolicy.ShowLockOutFailures)
			}
		})
	}
}

func TestAppendChangePasswordLockoutPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *LockoutPolicy
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
				iam: &IAM{DefaultLockoutPolicy: &LockoutPolicy{
					MaxPasswordAttempts: 10,
				}},
				policy: &LockoutPolicy{MaxPasswordAttempts: 5},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultLockoutPolicy: &LockoutPolicy{
				MaxPasswordAttempts: 5,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeLockoutPolicyEvent(tt.args.event)
			if tt.result.DefaultLockoutPolicy.MaxPasswordAttempts != tt.args.iam.DefaultLockoutPolicy.MaxPasswordAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLockoutPolicy.MaxPasswordAttempts, tt.args.iam.DefaultLockoutPolicy.MaxPasswordAttempts)
			}
		})
	}
}
