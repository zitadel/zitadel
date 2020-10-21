package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"testing"
)

func TestLoginPolicyChanges(t *testing.T) {
	type args struct {
		existing *LoginPolicy
		new      *LoginPolicy
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
				existing: &LoginPolicy{AllowUsernamePassword: false, AllowRegister: false, AllowExternalIdp: false},
				new:      &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
			},
			res: res{
				changesLen: 3,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &LoginPolicy{AllowUsernamePassword: false, AllowRegister: false, AllowExternalIdp: false},
				new:      &LoginPolicy{AllowUsernamePassword: false, AllowRegister: false, AllowExternalIdp: false},
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

func TestAppendAddLoginPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add login policy event",
			args: args{
				iam:    new(IAM),
				policy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultLoginPolicy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
		})
	}
}

func TestAppendChangeLoginPolicyEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change login policy event",
			args: args{
				iam: &IAM{DefaultLoginPolicy: &LoginPolicy{
					AllowExternalIdp:      false,
					AllowRegister:         false,
					AllowUsernamePassword: false,
				}},
				policy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
		})
	}
}

func TestAppendAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		iam      *IAM
		provider *IDPProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				iam:      &IAM{DefaultLoginPolicy: &LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				provider: &IDPProvider{Type: int32(model.IDPProviderTypeSystem), IDPConfigID: "IDPConfigID"},
				event:    &es_models.Event{},
			},
			result: &IAM{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IDPProviders: []*IDPProvider{
					{IDPConfigID: "IDPConfigID", Type: int32(model.IDPProviderTypeSystem)},
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddIDPProviderToLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.DefaultLoginPolicy.IDPProviders) != len(tt.args.iam.DefaultLoginPolicy.IDPProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.DefaultLoginPolicy.IDPProviders), len(tt.args.iam.DefaultLoginPolicy.IDPProviders))
			}
			if tt.result.DefaultLoginPolicy.IDPProviders[0].Type != tt.args.provider.Type {
				t.Errorf("got wrong idp provider type: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.IDPProviders[0].Type, tt.args.provider.Type)
			}
			if tt.result.DefaultLoginPolicy.IDPProviders[0].IDPConfigID != tt.args.provider.IDPConfigID {
				t.Errorf("got wrong idp provider idpconfigid: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.IDPProviders[0].IDPConfigID, tt.args.provider.IDPConfigID)
			}
		})
	}
}

func TestRemoveAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		iam      *IAM
		provider *IDPProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				iam: &IAM{
					DefaultLoginPolicy: &LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						IDPProviders: []*IDPProvider{
							{IDPConfigID: "IDPConfigID", Type: int32(model.IDPProviderTypeSystem)},
						}}},
				provider: &IDPProvider{Type: int32(model.IDPProviderTypeSystem), IDPConfigID: "IDPConfigID"},
				event:    &es_models.Event{},
			},
			result: &IAM{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IDPProviders:          []*IDPProvider{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveIDPProviderFromLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.DefaultLoginPolicy.IDPProviders) != len(tt.args.iam.DefaultLoginPolicy.IDPProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.DefaultLoginPolicy.IDPProviders), len(tt.args.iam.DefaultLoginPolicy.IDPProviders))
			}
		})
	}
}

func TestAppendAddSoftwareMFAToPolicyEvent(t *testing.T) {
	type args struct {
		iam   *IAM
		mfa   *MFA
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add software mfa to login policy event",
			args: args{
				iam:   &IAM{DefaultLoginPolicy: &LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				mfa:   &MFA{MfaType: int32(model.SoftwareMFATypeOTP)},
				event: &es_models.Event{},
			},
			result: &IAM{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				SoftwareMFAs: []int32{
					int32(model.SoftwareMFATypeOTP),
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddSoftwareMFAToLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.DefaultLoginPolicy.SoftwareMFAs) != len(tt.args.iam.DefaultLoginPolicy.SoftwareMFAs) {
				t.Errorf("got wrong software mfas len: expected: %v, actual: %v ", len(tt.result.DefaultLoginPolicy.SoftwareMFAs), len(tt.args.iam.DefaultLoginPolicy.SoftwareMFAs))
			}
			if tt.result.DefaultLoginPolicy.SoftwareMFAs[0] != tt.args.mfa.MfaType {
				t.Errorf("got wrong software mfa: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.SoftwareMFAs[0], tt.args.mfa)
			}
		})
	}
}
