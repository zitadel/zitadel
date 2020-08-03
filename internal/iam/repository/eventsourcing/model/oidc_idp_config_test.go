package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"testing"
)

func TestOIDCIdpConfigChanges(t *testing.T) {
	type args struct {
		existing *OidcIdpConfig
		new      *OidcIdpConfig
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
			name: "all possible values change",
			args: args{
				existing: &OidcIdpConfig{
					IdpConfigID:  "IdpConfigID",
					ClientID:     "ClientID",
					ClientSecret: &crypto.CryptoValue{KeyID: "KeyID"},
					Issuer:       "Issuer",
					Scopes:       []string{"scope1"},
				},
				new: &OidcIdpConfig{
					IdpConfigID:  "IdpConfigID",
					ClientID:     "ClientID2",
					ClientSecret: &crypto.CryptoValue{KeyID: "KeyID2"},
					Issuer:       "Issuer2",
					Scopes:       []string{"scope1", "scope2"},
				},
			},
			res: res{
				changesLen: 5,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &OidcIdpConfig{
					IdpConfigID: "IdpConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope1"},
				},
				new: &OidcIdpConfig{
					IdpConfigID: "IdpConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope1"},
				},
			},
			res: res{
				changesLen: 1,
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

func TestAppendAddOIDCIdpConfigEvent(t *testing.T) {
	type args struct {
		iam    *Iam
		config *OidcIdpConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add oidc idp config event",
			args: args{
				iam:    &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID"}}},
				config: &OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"},
				event:  &es_models.Event{},
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddOidcIdpConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0].OIDCIDPConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.iam.IDPs[0].OIDCIDPConfig)
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}

func TestAppendChangeOIDCIdpConfigEvent(t *testing.T) {
	type args struct {
		iam    *Iam
		config *OidcIdpConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append change oidc idp config event",
			args: args{
				iam:    &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"}}}},
				config: &OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID Changed"},
				event:  &es_models.Event{},
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID Changed"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeOidcIdpConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0].OIDCIDPConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.iam.IDPs[0].OIDCIDPConfig)
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}
