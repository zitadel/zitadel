package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	idpQuery = `SELECT projections.idps3.id,` +
		` projections.idps3.resource_owner,` +
		` projections.idps3.creation_date,` +
		` projections.idps3.change_date,` +
		` projections.idps3.sequence,` +
		` projections.idps3.state,` +
		` projections.idps3.name,` +
		` projections.idps3.styling_type,` +
		` projections.idps3.owner_type,` +
		` projections.idps3.auto_register,` +
		` projections.idps3_oidc_config.idp_id,` +
		` projections.idps3_oidc_config.client_id,` +
		` projections.idps3_oidc_config.client_secret,` +
		` projections.idps3_oidc_config.issuer,` +
		` projections.idps3_oidc_config.scopes,` +
		` projections.idps3_oidc_config.display_name_mapping,` +
		` projections.idps3_oidc_config.username_mapping,` +
		` projections.idps3_oidc_config.authorization_endpoint,` +
		` projections.idps3_oidc_config.token_endpoint,` +
		` projections.idps3_jwt_config.idp_id,` +
		` projections.idps3_jwt_config.issuer,` +
		` projections.idps3_jwt_config.keys_endpoint,` +
		` projections.idps3_jwt_config.header_name,` +
		` projections.idps3_jwt_config.endpoint` +
		` FROM projections.idps3` +
		` LEFT JOIN projections.idps3_oidc_config ON projections.idps3.id = projections.idps3_oidc_config.idp_id AND projections.idps3.instance_id = projections.idps3_oidc_config.instance_id` +
		` LEFT JOIN projections.idps3_jwt_config ON projections.idps3.id = projections.idps3_jwt_config.idp_id AND projections.idps3.instance_id = projections.idps3_jwt_config.instance_id`
	idpCols = []string{
		"id",
		"resource_owner",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"name",
		"styling_type",
		"owner_type",
		"auto_register",
		// oidc config
		"idp_id",
		"client_id",
		"client_secret",
		"issuer",
		"scopes",
		"display_name_mapping",
		"username_mapping",
		"authorization_endpoint",
		"token_endpoint",
		// jwt config
		"idp_id",
		"issuer",
		"keys_endpoint",
		"header_name",
		"endpoint",
	}
	idpsQuery = `SELECT projections.idps3.id,` +
		` projections.idps3.resource_owner,` +
		` projections.idps3.creation_date,` +
		` projections.idps3.change_date,` +
		` projections.idps3.sequence,` +
		` projections.idps3.state,` +
		` projections.idps3.name,` +
		` projections.idps3.styling_type,` +
		` projections.idps3.owner_type,` +
		` projections.idps3.auto_register,` +
		` projections.idps3_oidc_config.idp_id,` +
		` projections.idps3_oidc_config.client_id,` +
		` projections.idps3_oidc_config.client_secret,` +
		` projections.idps3_oidc_config.issuer,` +
		` projections.idps3_oidc_config.scopes,` +
		` projections.idps3_oidc_config.display_name_mapping,` +
		` projections.idps3_oidc_config.username_mapping,` +
		` projections.idps3_oidc_config.authorization_endpoint,` +
		` projections.idps3_oidc_config.token_endpoint,` +
		` projections.idps3_jwt_config.idp_id,` +
		` projections.idps3_jwt_config.issuer,` +
		` projections.idps3_jwt_config.keys_endpoint,` +
		` projections.idps3_jwt_config.header_name,` +
		` projections.idps3_jwt_config.endpoint,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idps3` +
		` LEFT JOIN projections.idps3_oidc_config ON projections.idps3.id = projections.idps3_oidc_config.idp_id AND projections.idps3.instance_id = projections.idps3_oidc_config.instance_id` +
		` LEFT JOIN projections.idps3_jwt_config ON projections.idps3.id = projections.idps3_jwt_config.idp_id AND projections.idps3.instance_id = projections.idps3_jwt_config.instance_id`
	idpsCols = []string{
		"id",
		"resource_owner",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"name",
		"styling_type",
		"owner_type",
		"auto_register",
		// oidc config
		"idp_id",
		"client_id",
		"client_secret",
		"issuer",
		"scopes",
		"display_name_mapping",
		"username_mapping",
		"authorization_endpoint",
		"token_endpoint",
		// jwt config
		"idp_id",
		"issuer",
		"keys_endpoint",
		"header_name",
		"endpoint",
		"count",
	}
)

func Test_IDPPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareIDPByIDQuery no result",
			prepare: prepareIDPByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpQuery),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*IDP)(nil),
		},
		{
			name:    "prepareIDPByIDQuery oidc idp",
			prepare: prepareIDPByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpQuery),
					idpCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPConfigStylingTypeGoogle,
						domain.IdentityProviderTypeOrg,
						true,
						// oidc config
						"idp-id",
						"oidc-client-id",
						nil,
						"oidc-issuer",
						database.StringArray{"scope"},
						domain.OIDCMappingFieldEmail,
						domain.OIDCMappingFieldPreferredLoginName,
						"auth.endpoint.ch",
						"token.endpoint.ch",
						// jwt config
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
			},
			object: &IDP{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				ResourceOwner: "ro",
				ID:            "idp-id",
				State:         domain.IDPConfigStateActive,
				Name:          "idp-name",
				StylingType:   domain.IDPConfigStylingTypeGoogle,
				OwnerType:     domain.IdentityProviderTypeOrg,
				AutoRegister:  true,
				OIDCIDP: &OIDCIDP{
					IDPID:                 "idp-id",
					ClientID:              "oidc-client-id",
					ClientSecret:          &crypto.CryptoValue{},
					Issuer:                "oidc-issuer",
					Scopes:                database.StringArray{"scope"},
					DisplayNameMapping:    domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
					AuthorizationEndpoint: "auth.endpoint.ch",
					TokenEndpoint:         "token.endpoint.ch",
				},
			},
		},
		{
			name:    "prepareIDPByIDQuery jwt config",
			prepare: prepareIDPByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpQuery),
					idpCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPConfigStylingTypeGoogle,
						domain.IdentityProviderTypeOrg,
						true,
						// oidc config
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						// jwt config
						"idp-id",
						"jwt-issuer",
						"key.ch",
						"x-header-name",
						"jwt.endpoint.ch",
					},
				),
			},
			object: &IDP{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				ResourceOwner: "ro",
				ID:            "idp-id",
				State:         domain.IDPConfigStateActive,
				Name:          "idp-name",
				StylingType:   domain.IDPConfigStylingTypeGoogle,
				OwnerType:     domain.IdentityProviderTypeOrg,
				AutoRegister:  true,
				JWTIDP: &JWTIDP{
					IDPID:        "idp-id",
					Issuer:       "jwt-issuer",
					KeysEndpoint: "key.ch",
					HeaderName:   "x-header-name",
					Endpoint:     "jwt.endpoint.ch",
				},
			},
		},
		{
			name:    "prepareIDPByIDQuery no config",
			prepare: prepareIDPByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpQuery),
					idpCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPConfigStylingTypeGoogle,
						domain.IdentityProviderTypeOrg,
						true,
						// oidc config
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						// jwt config
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
			},
			object: &IDP{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				ResourceOwner: "ro",
				ID:            "idp-id",
				State:         domain.IDPConfigStateActive,
				Name:          "idp-name",
				StylingType:   domain.IDPConfigStylingTypeGoogle,
				OwnerType:     domain.IdentityProviderTypeOrg,
				AutoRegister:  true,
			},
		},
		{
			name:    "prepareIDPByIDQuery sql err",
			prepare: prepareIDPByIDQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(idpQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
		{
			name:    "prepareIDPsQuery no result",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpsQuery),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: &IDPs{IDPs: []*IDP{}},
		},
		{
			name:    "prepareIDPsQuery oidc idp",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpsQuery),
					idpsCols,
					[][]driver.Value{
						{
							"idp-id",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							"idp-id",
							"oidc-client-id",
							nil,
							"oidc-issuer",
							database.StringArray{"scope"},
							domain.OIDCMappingFieldEmail,
							domain.OIDCMappingFieldPreferredLoginName,
							"auth.endpoint.ch",
							"token.endpoint.ch",
							// jwt config
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &IDPs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				IDPs: []*IDP{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
						OIDCIDP: &OIDCIDP{
							IDPID:                 "idp-id",
							ClientID:              "oidc-client-id",
							ClientSecret:          &crypto.CryptoValue{},
							Issuer:                "oidc-issuer",
							Scopes:                database.StringArray{"scope"},
							DisplayNameMapping:    domain.OIDCMappingFieldEmail,
							UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
							AuthorizationEndpoint: "auth.endpoint.ch",
							TokenEndpoint:         "token.endpoint.ch",
						},
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery jwt config",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpsQuery),
					idpsCols,
					[][]driver.Value{
						{
							"idp-id",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							// jwt config
							"idp-id",
							"jwt-issuer",
							"key.ch",
							"x-header-name",
							"jwt.endpoint.ch",
						},
					},
				),
			},
			object: &IDPs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				IDPs: []*IDP{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
						JWTIDP: &JWTIDP{
							IDPID:        "idp-id",
							Issuer:       "jwt-issuer",
							KeysEndpoint: "key.ch",
							HeaderName:   "x-header-name",
							Endpoint:     "jwt.endpoint.ch",
						},
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery no config",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpsQuery),
					idpsCols,
					[][]driver.Value{
						{
							"idp-id",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							// jwt config
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &IDPs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				IDPs: []*IDP{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery all config types",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpsQuery),
					idpsCols,
					[][]driver.Value{
						{
							"idp-id-1",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							// jwt config
							nil,
							nil,
							nil,
							nil,
							nil,
						},
						{
							"idp-id-2",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							"idp-id",
							"oidc-client-id",
							nil,
							"oidc-issuer",
							database.StringArray{"scope"},
							domain.OIDCMappingFieldEmail,
							domain.OIDCMappingFieldPreferredLoginName,
							"auth.endpoint.ch",
							"token.endpoint.ch",
							// jwt config
							nil,
							nil,
							nil,
							nil,
							nil,
						},
						{
							"idp-id-3",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPConfigStylingTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							// jwt config
							"idp-id",
							"jwt-issuer",
							"key.ch",
							"x-header-name",
							"jwt.endpoint.ch",
						},
					},
				),
			},
			object: &IDPs{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				IDPs: []*IDP{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id-1",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id-2",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
						OIDCIDP: &OIDCIDP{
							IDPID:                 "idp-id",
							ClientID:              "oidc-client-id",
							ClientSecret:          &crypto.CryptoValue{},
							Issuer:                "oidc-issuer",
							Scopes:                database.StringArray{"scope"},
							DisplayNameMapping:    domain.OIDCMappingFieldEmail,
							UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
							AuthorizationEndpoint: "auth.endpoint.ch",
							TokenEndpoint:         "token.endpoint.ch",
						},
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						ResourceOwner: "ro",
						ID:            "idp-id-3",
						State:         domain.IDPConfigStateActive,
						Name:          "idp-name",
						StylingType:   domain.IDPConfigStylingTypeGoogle,
						OwnerType:     domain.IdentityProviderTypeOrg,
						AutoRegister:  true,
						JWTIDP: &JWTIDP{
							IDPID:        "idp-id",
							Issuer:       "jwt-issuer",
							KeysEndpoint: "key.ch",
							HeaderName:   "x-header-name",
							Endpoint:     "jwt.endpoint.ch",
						},
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery sql err",
			prepare: prepareIDPsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(idpsQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
