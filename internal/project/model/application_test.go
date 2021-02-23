package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestApplicationValid(t *testing.T) {
	type args struct {
		app *Application
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "valid oidc application: responsetype code",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype code",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code & id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
					},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDTokenToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
					},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & id_token & token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       AppTypeOIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
						GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
					},
				},
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.app.IsValid(true)
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
