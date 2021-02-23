package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
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

func TestAppendAddSecondFactorToPolicyEvent(t *testing.T) {
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
			name: "append add second factor to login policy event",
			args: args{
				org:   &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				mfa:   &iam_es_model.MFA{MFAType: int32(iam_model.SecondFactorTypeOTP)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				SecondFactors: []int32{
					int32(iam_model.SecondFactorTypeOTP),
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddSecondFactorToLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.SecondFactors) != len(tt.args.org.LoginPolicy.SecondFactors) {
				t.Errorf("got wrong second factor len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.SecondFactors), len(tt.args.org.LoginPolicy.SecondFactors))
			}
			if tt.result.LoginPolicy.SecondFactors[0] != tt.args.mfa.MFAType {
				t.Errorf("got wrong second factor: expected: %v, actual: %v ", tt.result.LoginPolicy.SecondFactors[0], tt.args.mfa)
			}
		})
	}
}

func TestRemoveSecondFactorFromPolicyEvent(t *testing.T) {
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
			name: "append remove second factor from login policy event",
			args: args{
				org: &Org{
					LoginPolicy: &iam_es_model.LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						SecondFactors: []int32{
							int32(iam_model.SecondFactorTypeOTP),
						}}},
				mfa:   &iam_es_model.MFA{MFAType: int32(iam_model.SecondFactorTypeOTP)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				SecondFactors:         []int32{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveSecondFactorFromLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.SecondFactors) != len(tt.args.org.LoginPolicy.SecondFactors) {
				t.Errorf("got wrong idp mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.SecondFactors), len(tt.args.org.LoginPolicy.SecondFactors))
			}
		})
	}
}

func TestAppendAddMultiFactorToPolicyEvent(t *testing.T) {
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
			name: "append add mfa to login policy event",
			args: args{
				org:   &Org{LoginPolicy: &iam_es_model.LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				mfa:   &iam_es_model.MFA{MFAType: int32(iam_model.MultiFactorTypeU2FWithPIN)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				MultiFactors: []int32{
					int32(iam_model.MultiFactorTypeU2FWithPIN),
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddMultiFactorToLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.MultiFactors) != len(tt.args.org.LoginPolicy.MultiFactors) {
				t.Errorf("got wrong second factor len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.MultiFactors), len(tt.args.org.LoginPolicy.MultiFactors))
			}
			if tt.result.LoginPolicy.MultiFactors[0] != tt.args.mfa.MFAType {
				t.Errorf("got wrong second factor: expected: %v, actual: %v ", tt.result.LoginPolicy.MultiFactors[0], tt.args.mfa)
			}
		})
	}
}

func TestRemoveMultiFactorFromPolicyEvent(t *testing.T) {
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
			name: "append remove mfa from login policy event",
			args: args{
				org: &Org{
					LoginPolicy: &iam_es_model.LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						MultiFactors: []int32{
							int32(iam_model.MultiFactorTypeU2FWithPIN),
						}}},
				mfa:   &iam_es_model.MFA{MFAType: int32(iam_model.MultiFactorTypeU2FWithPIN)},
				event: &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				MultiFactors:          []int32{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mfa != nil {
				data, _ := json.Marshal(tt.args.mfa)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveMultiFactorFromLoginPolicyEvent(tt.args.event)
			if len(tt.result.LoginPolicy.MultiFactors) != len(tt.args.org.LoginPolicy.MultiFactors) {
				t.Errorf("got wrong idp mfa len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.MultiFactors), len(tt.args.org.LoginPolicy.MultiFactors))
			}
		})
	}
}
