package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
	"time"
)

func TestApplicationValid(t *testing.T) {
	type args struct {
		app *OIDCApp
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "no app name",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "invalid clock skew",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "AppName",
					ClockSkew:     time.Minute * 1,
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype code",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code & id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & id_token & token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.app.IsValid()
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
