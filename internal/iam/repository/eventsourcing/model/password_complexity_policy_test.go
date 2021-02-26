package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestPasswordComplexityPolicyChanges(t *testing.T) {
	type args struct {
		existing *PasswordComplexityPolicy
		new      *PasswordComplexityPolicy
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
			name: "loginpolicy all attributes change",
			args: args{
				existing: &PasswordComplexityPolicy{MinLength: 10, HasUppercase: true, HasLowercase: true, HasNumber: true, HasSymbol: true},
				new:      &PasswordComplexityPolicy{MinLength: 5, HasUppercase: false, HasLowercase: false, HasNumber: false, HasSymbol: false},
			},
			res: res{
				changesLen: 5,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &PasswordComplexityPolicy{MinLength: 10, HasUppercase: true, HasLowercase: true, HasNumber: true, HasSymbol: true},
				new:      &PasswordComplexityPolicy{MinLength: 10, HasUppercase: true, HasLowercase: true, HasNumber: true, HasSymbol: true},
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

func TestAppendAddPasswordComplexityPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordComplexityPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add password complexity policy event",
			args: args{
				iam:    new(IAM),
				policy: &PasswordComplexityPolicy{MinLength: 10, HasUppercase: true, HasLowercase: true, HasNumber: true, HasSymbol: true},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultPasswordComplexityPolicy: &PasswordComplexityPolicy{MinLength: 10, HasUppercase: true, HasLowercase: true, HasNumber: true, HasSymbol: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddPasswordComplexityPolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordComplexityPolicy.MinLength != tt.args.iam.DefaultPasswordComplexityPolicy.MinLength {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.MinLength, tt.args.iam.DefaultPasswordComplexityPolicy.MinLength)
			}
			if tt.result.DefaultPasswordComplexityPolicy.HasUppercase != tt.args.iam.DefaultPasswordComplexityPolicy.HasUppercase {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.HasUppercase, tt.args.iam.DefaultPasswordComplexityPolicy.HasUppercase)
			}
			if tt.result.DefaultPasswordComplexityPolicy.HasLowercase != tt.args.iam.DefaultPasswordComplexityPolicy.HasLowercase {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.HasLowercase, tt.args.iam.DefaultPasswordComplexityPolicy.HasLowercase)
			}
			if tt.result.DefaultPasswordComplexityPolicy.HasNumber != tt.args.iam.DefaultPasswordComplexityPolicy.HasNumber {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.HasNumber, tt.args.iam.DefaultPasswordComplexityPolicy.HasNumber)
			}
			if tt.result.DefaultPasswordComplexityPolicy.HasSymbol != tt.args.iam.DefaultPasswordComplexityPolicy.HasSymbol {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.HasSymbol, tt.args.iam.DefaultPasswordComplexityPolicy.HasSymbol)
			}
		})
	}
}

func TestAppendChangePasswordComplexityPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *PasswordComplexityPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change password complexity policy event",
			args: args{
				iam: &IAM{DefaultPasswordComplexityPolicy: &PasswordComplexityPolicy{
					MinLength: 10,
				}},
				policy: &PasswordComplexityPolicy{MinLength: 5},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultPasswordComplexityPolicy: &PasswordComplexityPolicy{
				MinLength: 5,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangePasswordComplexityPolicyEvent(tt.args.event)
			if tt.result.DefaultPasswordComplexityPolicy.MinLength != tt.args.iam.DefaultPasswordComplexityPolicy.MinLength {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultPasswordComplexityPolicy.MinLength, tt.args.iam.DefaultPasswordComplexityPolicy.MinLength)
			}
		})
	}
}
