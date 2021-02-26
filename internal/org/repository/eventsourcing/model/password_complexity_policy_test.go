package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddPasswordComplexityPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordComplexityPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add password complexity policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.PasswordComplexityPolicy{MinLength: 10},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordComplexityPolicy: &iam_es_model.PasswordComplexityPolicy{MinLength: 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddPasswordComplexityPolicyEvent(tt.args.event)
			if tt.result.PasswordComplexityPolicy.MinLength != tt.args.org.PasswordComplexityPolicy.MinLength {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordComplexityPolicy.MinLength, tt.args.org.PasswordComplexityPolicy.MinLength)
			}
		})
	}
}

func TestAppendChangePasswordComplexityPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.PasswordComplexityPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change password complexity policy event",
			args: args{
				org: &Org{PasswordComplexityPolicy: &iam_es_model.PasswordComplexityPolicy{
					MinLength: 10,
				}},
				policy: &iam_es_model.PasswordComplexityPolicy{MinLength: 5, HasLowercase: true},
				event:  &es_models.Event{},
			},
			result: &Org{PasswordComplexityPolicy: &iam_es_model.PasswordComplexityPolicy{
				MinLength:    5,
				HasLowercase: true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangePasswordComplexityPolicyEvent(tt.args.event)
			if tt.result.PasswordComplexityPolicy.MinLength != tt.args.org.PasswordComplexityPolicy.MinLength {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordComplexityPolicy.MinLength, tt.args.org.PasswordComplexityPolicy.MinLength)
			}
			if tt.result.PasswordComplexityPolicy.HasLowercase != tt.args.org.PasswordComplexityPolicy.HasLowercase {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PasswordComplexityPolicy.HasLowercase, tt.args.org.PasswordComplexityPolicy.HasLowercase)
			}
		})
	}
}
