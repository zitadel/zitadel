package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestOrgIAMPolicyChanges(t *testing.T) {
	type args struct {
		existing *OrgIAMPolicy
		new      *OrgIAMPolicy
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
			name: "org iam policy all attributes change",
			args: args{
				existing: &OrgIAMPolicy{UserLoginMustBeDomain: true},
				new:      &OrgIAMPolicy{UserLoginMustBeDomain: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &OrgIAMPolicy{UserLoginMustBeDomain: true},
				new:      &OrgIAMPolicy{UserLoginMustBeDomain: true},
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

func TestAppendAddOrgIAMPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *OrgIAMPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add org iam policy event",
			args: args{
				iam:    new(IAM),
				policy: &OrgIAMPolicy{UserLoginMustBeDomain: true},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultOrgIAMPolicy: &OrgIAMPolicy{UserLoginMustBeDomain: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddOrgIAMPolicyEvent(tt.args.event)
			if tt.result.DefaultOrgIAMPolicy.UserLoginMustBeDomain != tt.args.iam.DefaultOrgIAMPolicy.UserLoginMustBeDomain {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultOrgIAMPolicy.UserLoginMustBeDomain, tt.args.iam.DefaultOrgIAMPolicy.UserLoginMustBeDomain)
			}
		})
	}
}

func TestAppendChangeOrgIAMPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *OrgIAMPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change org iam policy event",
			args: args{
				iam: &IAM{DefaultOrgIAMPolicy: &OrgIAMPolicy{
					UserLoginMustBeDomain: true,
				}},
				policy: &OrgIAMPolicy{UserLoginMustBeDomain: false},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultOrgIAMPolicy: &OrgIAMPolicy{
				UserLoginMustBeDomain: false,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeOrgIAMPolicyEvent(tt.args.event)
			if tt.result.DefaultOrgIAMPolicy.UserLoginMustBeDomain != tt.args.iam.DefaultOrgIAMPolicy.UserLoginMustBeDomain {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultOrgIAMPolicy.UserLoginMustBeDomain, tt.args.iam.DefaultOrgIAMPolicy.UserLoginMustBeDomain)
			}
		})
	}
}
