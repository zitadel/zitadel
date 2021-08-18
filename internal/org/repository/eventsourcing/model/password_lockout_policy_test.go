package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddLockoutPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LockoutPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add lockout policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.LockoutPolicy{MaxPasswordAttempts: 10},
				event:  &es_models.Event{},
			},
			result: &Org{LockoutPolicy: &iam_es_model.LockoutPolicy{MaxPasswordAttempts: 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddLockoutPolicyEvent(tt.args.event)
			if tt.result.LockoutPolicy.MaxPasswordAttempts != tt.args.org.LockoutPolicy.MaxPasswordAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LockoutPolicy.MaxPasswordAttempts, tt.args.org.LockoutPolicy.MaxPasswordAttempts)
			}
		})
	}
}

func TestAppendChangeLockoutPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LockoutPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change lockout policy event",
			args: args{
				org: &Org{LockoutPolicy: &iam_es_model.LockoutPolicy{
					MaxPasswordAttempts: 10,
				}},
				policy: &iam_es_model.LockoutPolicy{MaxPasswordAttempts: 5, ShowLockOutFailures: true},
				event:  &es_models.Event{},
			},
			result: &Org{LockoutPolicy: &iam_es_model.LockoutPolicy{
				MaxPasswordAttempts: 5,
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
			tt.args.org.appendChangeLockoutPolicyEvent(tt.args.event)
			if tt.result.LockoutPolicy.MaxPasswordAttempts != tt.args.org.LockoutPolicy.MaxPasswordAttempts {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LockoutPolicy.MaxPasswordAttempts, tt.args.org.LockoutPolicy.MaxPasswordAttempts)
			}
			if tt.result.LockoutPolicy.ShowLockOutFailures != tt.args.org.LockoutPolicy.ShowLockOutFailures {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LockoutPolicy.ShowLockOutFailures, tt.args.org.LockoutPolicy.ShowLockOutFailures)
			}
		})
	}
}
