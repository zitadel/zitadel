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
		provider *iam_es_model.IdpProvider
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
				provider: &iam_es_model.IdpProvider{Type: int32(iam_model.IdpProviderTypeSystem), IdpConfigID: "IdpConfigID"},
				event:    &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IdpProviders: []*iam_es_model.IdpProvider{
					{IdpConfigID: "IdpConfigID", Type: int32(iam_model.IdpProviderTypeSystem)},
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
				t.Errorf("got wrong result AllowExternalIdp: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.LoginPolicy.IdpProviders) != len(tt.args.org.LoginPolicy.IdpProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.IdpProviders), len(tt.args.org.LoginPolicy.IdpProviders))
			}
			if tt.result.LoginPolicy.IdpProviders[0].Type != tt.args.provider.Type {
				t.Errorf("got wrong idp provider type: expected: %v, actual: %v ", tt.result.LoginPolicy.IdpProviders[0].Type, tt.args.provider.Type)
			}
			if tt.result.LoginPolicy.IdpProviders[0].IdpConfigID != tt.args.provider.IdpConfigID {
				t.Errorf("got wrong idp provider idpconfigid: expected: %v, actual: %v ", tt.result.LoginPolicy.IdpProviders[0].IdpConfigID, tt.args.provider.IdpConfigID)
			}
		})
	}
}

func TestRemoveAddIdpToPolicyEvent(t *testing.T) {
	type args struct {
		org      *Org
		provider *iam_es_model.IdpProvider
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
						IdpProviders: []*iam_es_model.IdpProvider{
							{IdpConfigID: "IdpConfigID", Type: int32(iam_model.IdpProviderTypeSystem)},
						}}},
				provider: &iam_es_model.IdpProvider{Type: int32(iam_model.IdpProviderTypeSystem), IdpConfigID: "IdpConfigID"},
				event:    &es_models.Event{},
			},
			result: &Org{LoginPolicy: &iam_es_model.LoginPolicy{
				AllowExternalIdp:      true,
				AllowRegister:         true,
				AllowUsernamePassword: true,
				IdpProviders:          []*iam_es_model.IdpProvider{}}},
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
				t.Errorf("got wrong result AllowExternalIdp: expected: %v, actual: %v ", tt.result.LoginPolicy.AllowExternalIdp, tt.args.org.LoginPolicy.AllowExternalIdp)
			}
			if len(tt.result.LoginPolicy.IdpProviders) != len(tt.args.org.LoginPolicy.IdpProviders) {
				t.Errorf("got wrong idp provider len: expected: %v, actual: %v ", len(tt.result.LoginPolicy.IdpProviders), len(tt.args.org.LoginPolicy.IdpProviders))
			}
		})
	}
}
