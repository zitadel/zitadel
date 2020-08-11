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
		iam    *Iam
		policy *LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add login policy event",
			args: args{
				iam:    &Iam{},
				policy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  &es_models.Event{},
			},
			result: &Iam{DefaultLoginPolicy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true}},
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
		iam    *Iam
		policy *LoginPolicy
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append change login policy event",
			args: args{
				iam: &Iam{DefaultLoginPolicy: &LoginPolicy{
					AllowExternalIdp:      false,
					AllowRegister:         false,
					AllowUsernamePassword: false,
				}},
				policy: &LoginPolicy{AllowUsernamePassword: true, AllowRegister: true, AllowExternalIdp: true},
				event:  &es_models.Event{},
			},
			result: &Iam{DefaultLoginPolicy: &LoginPolicy{
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
		iam      *Iam
		provider *IdpProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				iam:      &Iam{DefaultLoginPolicy: &LoginPolicy{AllowExternalIdp: true, AllowRegister: true, AllowUsernamePassword: true}},
				provider: &IdpProvider{Type: int32(model.IdpProviderTypeSystem), IdpConfigID: "IdpConfigID"},
				event:    &es_models.Event{},
			},
			result: &Iam{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IdpProviders: []*IdpProvider{
					{IdpConfigID: "IdpConfigID", Type: int32(model.IdpProviderTypeSystem)},
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddIdpProviderToLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIdp: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.DefaultLoginPolicy.IdpProviders) != len(tt.args.iam.DefaultLoginPolicy.IdpProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.DefaultLoginPolicy.IdpProviders), len(tt.args.iam.DefaultLoginPolicy.IdpProviders))
			}
			if tt.result.DefaultLoginPolicy.IdpProviders[0].Type != tt.args.provider.Type {
				t.Errorf("got wrong idp provider type: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.IdpProviders[0].Type, tt.args.provider.Type)
			}
			if tt.result.DefaultLoginPolicy.IdpProviders[0].IdpConfigID != tt.args.provider.IdpConfigID {
				t.Errorf("got wrong idp provider idpconfigid: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.IdpProviders[0].IdpConfigID, tt.args.provider.IdpConfigID)
			}
		})
	}
}

func TestRemoveAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		iam      *Iam
		provider *IdpProvider
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add idp to login policy event",
			args: args{
				iam: &Iam{
					DefaultLoginPolicy: &LoginPolicy{
						AllowExternalIdp:      true,
						AllowRegister:         true,
						AllowUsernamePassword: true,
						IdpProviders: []*IdpProvider{
							{IdpConfigID: "IdpConfigID", Type: int32(model.IdpProviderTypeSystem)},
						}}},
				provider: &IdpProvider{Type: int32(model.IdpProviderTypeSystem), IdpConfigID: "IdpConfigID"},
				event:    &es_models.Event{},
			},
			result: &Iam{DefaultLoginPolicy: &LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IdpProviders:          []*IdpProvider{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.provider != nil {
				data, _ := json.Marshal(tt.args.provider)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveIdpProviderFromLoginPolicyEvent(tt.args.event)
			if tt.result.DefaultLoginPolicy.AllowUsernamePassword != tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowUsernamePassword, tt.args.iam.DefaultLoginPolicy.AllowUsernamePassword)
			}
			if tt.result.DefaultLoginPolicy.AllowRegister != tt.args.iam.DefaultLoginPolicy.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowRegister, tt.args.iam.DefaultLoginPolicy.AllowRegister)
			}
			if tt.result.DefaultLoginPolicy.AllowExternalIdp != tt.args.iam.DefaultLoginPolicy.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIdp: expected: %v, actual: %v ", tt.result.DefaultLoginPolicy.AllowExternalIdp, tt.args.iam.DefaultLoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.DefaultLoginPolicy.IdpProviders) != len(tt.args.iam.DefaultLoginPolicy.IdpProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.DefaultLoginPolicy.IdpProviders), len(tt.args.iam.DefaultLoginPolicy.IdpProviders))
			}
		})
	}
}
