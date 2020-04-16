package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
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
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_CODE},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype code",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_CODE},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_IMPLICIT},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_IMPLICIT},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_TOKEN_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_IMPLICIT},
					},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_TOKEN_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code & id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_CODE, OIDCRESPONSETYPE_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE, OIDCGRANTTYPE_IMPLICIT},
					},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_CODE, OIDCRESPONSETYPE_TOKEN_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE, OIDCGRANTTYPE_IMPLICIT},
					},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & id_token & token_id_token",
			args: args{
				app: &Application{
					ObjectRoot: models.ObjectRoot{ID: "ID"},
					AppID:      "AppID",
					Name:       "Name",
					Type:       APPTYPE_OIDC,
					OIDCConfig: &OIDCConfig{
						ResponseTypes: []OIDCResponseType{OIDCRESPONSETYPE_CODE, OIDCRESPONSETYPE_ID_TOKEN, OIDCRESPONSETYPE_TOKEN_ID_TOKEN},
						GrantTypes:    []OIDCGrantType{OIDCGRANTTYPE_AUTHORIZATION_CODE, OIDCGRANTTYPE_IMPLICIT},
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
