package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddLoginPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add login policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddLoginPolicyEvent(tt.args.event)
			if tt.result.LoginPolicy.AllowUsernamePassword != tt.args.org.LoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowUsernamePassword, tt.args.org.LoginPolicy.AllowUsernamePassword)
			}
			if tt.result.LoginPolicy.AllowRegister != tt.args.org.LoginPolicy.AllowRegister {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowRegister, tt.args.org.LoginPolicy.AllowRegister)
			}
			if tt.result.LoginPolicy.AllowExternalIdp != tt.args.org.LoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
		})
	}
}

func TestAppendChangeLoginPolicyEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change login policy event",
			args: args{
				org: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
					AllowExternalIdp:      false,
					AllowRegister:         false,
					AllowUsernamePassword: false,
				}},
				policy: &iam_es_model.LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
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
			tt.args.org.appendChangeLoginPolicyEvent(tt.args.event)
			if tt.result.LoginPolicy.AllowUsernamePassword != tt.args.org.LoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowUsernamePassword, tt.args.org.LoginPolicy.AllowUsernamePassword)
			}
			if tt.result.LoginPolicy.AllowRegister != tt.args.org.LoginPolicy.AllowRegister {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowRegister, tt.args.org.LoginPolicy.AllowRegister)
			}
			if tt.result.LoginPolicy.AllowExternalIdp != tt.args.org.LoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
		})
	}
}

func TestAppendAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		org      *Org
		provider *iam_es_model.IDPProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				org:      &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				provider: &iam_es_model.IDPProvider{Type: int32(iam_model.IDPProviderTypeSystem), IDPConfigID: "IDPConfigID"},
				event:    &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IDPProviders: []*iam_es_model.IDPProvider{
					{IDPConfigID: "IDPConfigID", Type: int32(iam_model.IDPProviderTypeSystem)},
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddIdpProviderToLoginPolicyEvent(tt.args.event)
			if tt.result.LoginPolicy.AllowUsernamePassword != tt.args.org.LoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowUsernamePassword, tt.args.org.LoginPolicy.AllowUsernamePassword)
			}
			if tt.result.LoginPolicy.AllowRegister != tt.args.org.LoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowRegister, tt.args.org.LoginPolicy.AllowRegister)
			}
			if tt.result.LoginPolicy.AllowExternalIdp != tt.args.org.LoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.LoginPolicy.IDPProviders) != len(tt.args.org.LoginPolicy.IDPProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.IDPProviders), len(tt.args.org.LoginPolicy.IDPProviders))
			}
			if tt.result.LoginPolicy.IDPProviders[0].Type != tt.args.provider.Type {
				t.Errorf("got wrong idp provider type: expected: %v, actual: %v ", tt.result.LoginPolicy.IDPProviders[0].Type, tt.args.provider.Type)
			}
			if tt.result.LoginPolicy.IDPProviders[0].IDPConfigID != tt.args.provider.IDPConfigID {
				t.Errorf("got wrong idp provider idpconfigid: expected: %v, actual: %v ", tt.result.LoginPolicy.IDPProviders[0].IDPConfigID, tt.args.provider.IDPConfigID)
			}
		})
	}
}

func TestRemoveAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		org      *Org
		provider *iam_es_model.IDPProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				org: &Org{
					LoginPolicy: &iam_es_model.LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						IDPProviders: []*iam_es_model.IDPProvider{
							{IDPConfigID: "IDPConfigID", Type: int32(iam_model.IDPProviderTypeSystem)},
						}}},
				provider: &iam_es_model.IDPProvider{Type: int32(iam_model.IDPProviderTypeSystem), IDPConfigID: "IDPConfigID"},
				event:    &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IDPProviders:          []*iam_es_model.IDPProvider{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveIdpProviderFromLoginPolicyEvent(tt.args.event)
			if tt.result.LoginPolicy.AllowUsernamePassword != tt.args.org.LoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowUsernamePassword, tt.args.org.LoginPolicy.AllowUsernamePassword)
			}
			if tt.result.LoginPolicy.AllowRegister != tt.args.org.LoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowRegister, tt.args.org.LoginPolicy.AllowRegister)
			}
			if tt.result.LoginPolicy.AllowExternalIdp != tt.args.org.LoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.LoginPolicy.IDPProviders) != len(tt.args.org.LoginPolicy.IDPProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.IDPProviders), len(tt.args.org.LoginPolicy.IDPProviders))
			}
		})
	}
}

func TestAppendAddSoftwareMFAToPolicyEvent(t *testing.T) {
	type args struct {
		org   *Org
		mfa   *iam_es_model.MFA
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add software mfa to login policy event",
			args: args{
				org:   &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				mfa:   &iam_es_model.MFA{MfaType: int32(iam_model.SoftwareMFATypeOTP)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				SoftwareMFAs: []int32{
					int32(iam_model.SoftwareMFATypeOTP),
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddSoftwareMFAToLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.SoftwareMFAs) != len(tt.args.org.LoginPolicy.SoftwareMFAs) {
				t.Errorf("got wrong software mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.SoftwareMFAs), len(tt.args.org.LoginPolicy.SoftwareMFAs))
			}
			if tt.result.LoginPolicy.SoftwareMFAs[0] != tt.args.mfa.MfaType {
				t.Errorf("got wrong software mfa: expected: %v, actual: %v ", tt.result.LoginPolicy.SoftwareMFAs[0], tt.args.mfa)
			}
		})
	}
}

func TestRemoveSoftwareMFAFromPolicyEvent(t *testing.T) {
	type args struct {
		org   *Org
		mfa   *iam_es_model.MFA
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append remove software mfa from login policy event",
			args: args{
				org: &Org{
					LoginPolicy: &iam_es_model.LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						SoftwareMFAs: []int32{
							int32(iam_model.SoftwareMFATypeOTP),
						}}},
				mfa:   &iam_es_model.MFA{MfaType: int32(iam_model.SoftwareMFATypeOTP)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				SoftwareMFAs:          []int32{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveSoftwareMFAFromLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.SoftwareMFAs) != len(tt.args.org.LoginPolicy.SoftwareMFAs) {
				t.Errorf("got wrong idp mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.SoftwareMFAs), len(tt.args.org.LoginPolicy.SoftwareMFAs))
			}
		})
	}
}

func TestAppendAddHardwareMFAToPolicyEvent(t *testing.T) {
	type args struct {
		org   *Org
		mfa   *iam_es_model.MFA
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add hardware mfa to login policy event",
			args: args{
				org:   &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				mfa:   &iam_es_model.MFA{MfaType: int32(iam_model.HardwareMFATypeU2F)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				HardwareMFAs: []int32{
					int32(iam_model.HardwareMFATypeU2F),
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddHardwareMFAToLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.HardwareMFAs) != len(tt.args.org.LoginPolicy.HardwareMFAs) {
				t.Errorf("got wrong software mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.HardwareMFAs), len(tt.args.org.LoginPolicy.HardwareMFAs))
			}
			if tt.result.LoginPolicy.HardwareMFAs[0] != tt.args.mfa.MfaType {
				t.Errorf("got wrong software mfa: expected: %v, actual: %v ", tt.result.LoginPolicy.HardwareMFAs[0], tt.args.mfa)
			}
		})
	}
}

func TestRemoveHardwareMFAFromPolicyEvent(t *testing.T) {
	type args struct {
		org   *Org
		mfa   *iam_es_model.MFA
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append remove hardware mfa from login policy event",
			args: args{
				org: &Org{
					LoginPolicy: &iam_es_model.LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						HardwareMFAs: []int32{
							int32(iam_model.HardwareMFATypeU2F),
						}}},
				mfa:   &iam_es_model.MFA{MfaType: int32(iam_model.HardwareMFATypeU2F)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				HardwareMFAs:          []int32{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveHardwareMFAFromLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.HardwareMFAs) != len(tt.args.org.LoginPolicy.HardwareMFAs) {
				t.Errorf("got wrong idp mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.HardwareMFAs), len(tt.args.org.LoginPolicy.HardwareMFAs))
			}
		})
	}
}
