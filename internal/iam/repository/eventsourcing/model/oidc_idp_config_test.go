package model

import (
	"testing"

	"github.com/caos/zitadel/internal/crypto"
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
