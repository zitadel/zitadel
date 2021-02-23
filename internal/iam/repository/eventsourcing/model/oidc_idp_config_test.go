package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestOIDCIdpConfigChanges(t *testing.T) {
	type args struct {
		existing *OIDCIDPConfig
		new      *OIDCIDPConfig
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
				existing: &OIDCIDPConfig{
					IDPConfigID:  "IDPConfigID",
					ClientID:     "ClientID",
					ClientSecret: &crypto.CryptoValue{KeyID: "KeyID"},
					Issuer:       "Issuer",
					Scopes:       []string{"scope1"},
				},
				new: &OIDCIDPConfig{
					IDPConfigID:  "IDPConfigID",
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
				existing: &OIDCIDPConfig{
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope1"},
				},
				new: &OIDCIDPConfig{
					IDPConfigID: "IDPConfigID",
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
		iam    *IAM
		config *OIDCIDPConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add oidc idp config event",
			args: args{
				iam:    &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID"}}},
				config: &OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"},
				event:  &es_models.Event{},
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddOIDCIDPConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.iam.IDPs))
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
		iam    *IAM
		config *OIDCIDPConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change oidc idp config event",
			args: args{
				iam:    &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}}}},
				config: &OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID Changed"},
				event:  &es_models.Event{},
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID Changed"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeOIDCIDPConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.iam.IDPs))
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
