package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

var (
	idpTemplateQuery = `SELECT projections.idp_templates.id,` +
		` projections.idp_templates.resource_owner,` +
		` projections.idp_templates.creation_date,` +
		` projections.idp_templates.change_date,` +
		` projections.idp_templates.sequence,` +
		` projections.idp_templates.state,` +
		` projections.idp_templates.name,` +
		` projections.idp_templates.type,` +
		` projections.idp_templates.owner_type,` +
		` projections.idp_templates.is_creation_allowed,` +
		` projections.idp_templates.is_linking_allowed,` +
		` projections.idp_templates.is_auto_creation,` +
		` projections.idp_templates.is_auto_update,` +
		` projections.idp_templates_google.idp_id,` +
		` projections.idp_templates_google.client_id,` +
		` projections.idp_templates_google.client_secret,` +
		` projections.idp_templates_google.scopes,` +
		` projections.idp_templates_ldap.idp_id,` +
		` projections.idp_templates_ldap.host,` +
		` projections.idp_templates_ldap.port,` +
		` projections.idp_templates_ldap.tls,` +
		` projections.idp_templates_ldap.base_dn,` +
		` projections.idp_templates_ldap.user_object_class,` +
		` projections.idp_templates_ldap.user_unique_attribute,` +
		` projections.idp_templates_ldap.admin,` +
		` projections.idp_templates_ldap.password,` +
		` projections.idp_templates_ldap.id_attribute,` +
		` projections.idp_templates_ldap.first_name_attribute,` +
		` projections.idp_templates_ldap.last_name_attribute,` +
		` projections.idp_templates_ldap.display_name_attribute,` +
		` projections.idp_templates_ldap.nick_name_attribute,` +
		` projections.idp_templates_ldap.preferred_username_attribute,` +
		` projections.idp_templates_ldap.email_attribute,` +
		` projections.idp_templates_ldap.email_verified,` +
		` projections.idp_templates_ldap.phone_attribute,` +
		` projections.idp_templates_ldap.phone_verified_attribute,` +
		` projections.idp_templates_ldap.preferred_language_attribute,` +
		` projections.idp_templates_ldap.avatar_url_attribute,` +
		` projections.idp_templates_ldap.profile_attribute` +
		` FROM projections.idp_templates` +
		` LEFT JOIN projections.idp_templates_google ON projections.idp_templates.id = projections.idp_templates_google.idp_id AND projections.idp_templates.instance_id = projections.idp_templates_google.instance_id` +
		` LEFT JOIN projections.idp_templates_ldap ON projections.idp_templates.id = projections.idp_templates_ldap.idp_id AND projections.idp_templates.instance_id = projections.idp_templates_ldap.instance_id`
	idpTemplateCols = []string{
		"id",
		"resource_owner",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"name",
		"type",
		"owner_type",
		"is_creation_allowed",
		"is_linking_allowed",
		"is_auto_creation",
		"is_auto_update",
		// google config
		"idp_id",
		"client_id",
		"client_secret",
		"scopes",
		// ldap config
		"idp_id",
		"host",
		"port",
		"tls",
		"base_dn",
		"user_object_class",
		"user_unique_attribute",
		"admin",
		"password",
		"id_attribute",
		"first_name_attribute",
		"last_name_attribute",
		"display_name_attribute",
		"nick_name_attribute",
		"preferred_username_attribute",
		"email_attribute",
		"email_verified",
		"phone_attribute",
		"phone_verified_attribute",
		"preferred_language_attribute",
		"avatar_url_attribute",
		"profile_attribute",
	}
	idpTemplatesQuery = `SELECT projections.idp_templates.id,` +
		` projections.idp_templates.resource_owner,` +
		` projections.idp_templates.creation_date,` +
		` projections.idp_templates.change_date,` +
		` projections.idp_templates.sequence,` +
		` projections.idp_templates.state,` +
		` projections.idp_templates.name,` +
		` projections.idp_templates.type,` +
		` projections.idp_templates.owner_type,` +
		` projections.idp_templates.is_creation_allowed,` +
		` projections.idp_templates.is_linking_allowed,` +
		` projections.idp_templates.is_auto_creation,` +
		` projections.idp_templates.is_auto_update,` +
		` projections.idp_templates_google.idp_id,` +
		` projections.idp_templates_google.client_id,` +
		` projections.idp_templates_google.client_secret,` +
		` projections.idp_templates_google.scopes,` +
		` projections.idp_templates_ldap.idp_id,` +
		` projections.idp_templates_ldap.host,` +
		` projections.idp_templates_ldap.port,` +
		` projections.idp_templates_ldap.tls,` +
		` projections.idp_templates_ldap.base_dn,` +
		` projections.idp_templates_ldap.user_object_class,` +
		` projections.idp_templates_ldap.user_unique_attribute,` +
		` projections.idp_templates_ldap.admin,` +
		` projections.idp_templates_ldap.password,` +
		` projections.idp_templates_ldap.id_attribute,` +
		` projections.idp_templates_ldap.first_name_attribute,` +
		` projections.idp_templates_ldap.last_name_attribute,` +
		` projections.idp_templates_ldap.display_name_attribute,` +
		` projections.idp_templates_ldap.nick_name_attribute,` +
		` projections.idp_templates_ldap.preferred_username_attribute,` +
		` projections.idp_templates_ldap.email_attribute,` +
		` projections.idp_templates_ldap.email_verified,` +
		` projections.idp_templates_ldap.phone_attribute,` +
		` projections.idp_templates_ldap.phone_verified_attribute,` +
		` projections.idp_templates_ldap.preferred_language_attribute,` +
		` projections.idp_templates_ldap.avatar_url_attribute,` +
		` projections.idp_templates_ldap.profile_attribute,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idp_templates` +
		` LEFT JOIN projections.idp_templates_google ON projections.idp_templates.id = projections.idp_templates_google.idp_id AND projections.idp_templates.instance_id = projections.idp_templates_google.instance_id` +
		` LEFT JOIN projections.idp_templates_ldap ON projections.idp_templates.id = projections.idp_templates_ldap.idp_id AND projections.idp_templates.instance_id = projections.idp_templates_ldap.instance_id`
	idpTemplatesCols = []string{
		"id",
		"resource_owner",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"name",
		"type",
		"owner_type",
		"is_creation_allowed",
		"is_linking_allowed",
		"is_auto_creation",
		"is_auto_update",
		// google config
		"idp_id",
		"client_id",
		"client_secret",
		"scopes",
		// ldap config
		"idp_id",
		"host",
		"port",
		"tls",
		"base_dn",
		"user_object_class",
		"user_unique_attribute",
		"admin",
		"password",
		"id_attribute",
		"first_name_attribute",
		"last_name_attribute",
		"display_name_attribute",
		"nick_name_attribute",
		"preferred_username_attribute",
		"email_attribute",
		"email_verified",
		"phone_attribute",
		"phone_verified_attribute",
		"preferred_language_attribute",
		"avatar_url_attribute",
		"profile_attribute",
		"count",
	}
)

func Test_IDPTemplateTemplatesPrepares(t *testing.T) {
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
			name:    "prepareIDPTemplateByIDQuery no result",
			prepare: prepareIDPTemplateByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpTemplateQuery),
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
			object: (*IDPTemplate)(nil),
		},

		{
			name:    "prepareIDPTemplateByIDQuery google idp",
			prepare: prepareIDPTemplateByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpTemplateQuery),
					idpTemplateCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPTypeGoogle,
						domain.IdentityProviderTypeOrg,
						true,
						true,
						true,
						true,
						// google
						"idp-id",
						"client_id",
						nil,
						database.StringArray{"profile"},
						// ldap config
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
			},
			object: &IDPTemplate{
				CreationDate:      testNow,
				ChangeDate:        testNow,
				Sequence:          20211109,
				ResourceOwner:     "ro",
				ID:                "idp-id",
				State:             domain.IDPStateActive,
				Name:              "idp-name",
				Type:              domain.IDPTypeGoogle,
				OwnerType:         domain.IdentityProviderTypeOrg,
				IsCreationAllowed: true,
				IsLinkingAllowed:  true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				GoogleIDPTemplate: &GoogleIDPTemplate{
					IDPID:        "idp-id",
					ClientID:     "client_id",
					ClientSecret: nil,
					Scopes:       []string{"profile"},
				},
			},
		},
		{
			name:    "prepareIDPTemplateByIDQuery ldap idp",
			prepare: prepareIDPTemplateByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpTemplateQuery),
					idpTemplateCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPTypeLDAP,
						domain.IdentityProviderTypeOrg,
						true,
						true,
						true,
						true,
						// google
						nil,
						nil,
						nil,
						nil,
						// ldap config
						"idp-id",
						"host",
						"port",
						true,
						"base",
						"user",
						"uid",
						"admin",
						nil,
						"id",
						"first",
						"last",
						"display",
						"nickname",
						"username",
						"email",
						"emailVerified",
						"phone",
						"phoneVerified",
						"lang",
						"avatar",
						"profile",
					},
				),
			},
			object: &IDPTemplate{
				CreationDate:      testNow,
				ChangeDate:        testNow,
				Sequence:          20211109,
				ResourceOwner:     "ro",
				ID:                "idp-id",
				State:             domain.IDPStateActive,
				Name:              "idp-name",
				Type:              domain.IDPTypeLDAP,
				OwnerType:         domain.IdentityProviderTypeOrg,
				IsCreationAllowed: true,
				IsLinkingAllowed:  true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				LDAPIDPTemplate: &LDAPIDPTemplate{
					IDPID:               "idp-id",
					Host:                "host",
					Port:                "port",
					TLS:                 true,
					BaseDN:              "base",
					UserObjectClass:     "user",
					UserUniqueAttribute: "uid",
					Admin:               "admin",
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "id",
						FirstNameAttribute:         "first",
						LastNameAttribute:          "last",
						DisplayNameAttribute:       "display",
						NickNameAttribute:          "nickname",
						PreferredUsernameAttribute: "username",
						EmailAttribute:             "email",
						EmailVerifiedAttribute:     "emailVerified",
						PhoneAttribute:             "phone",
						PhoneVerifiedAttribute:     "phoneVerified",
						PreferredLanguageAttribute: "lang",
						AvatarURLAttribute:         "avatar",
						ProfileAttribute:           "profile",
					},
				},
			},
		},
		{
			name:    "prepareIDPTemplateByIDQuery no config",
			prepare: prepareIDPTemplateByIDQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(idpTemplateQuery),
					idpTemplateCols,
					[]driver.Value{
						"idp-id",
						"ro",
						testNow,
						testNow,
						uint64(20211109),
						domain.IDPConfigStateActive,
						"idp-name",
						domain.IDPTypeLDAP,
						domain.IdentityProviderTypeOrg,
						true,
						true,
						true,
						true,
						// google config
						nil,
						nil,
						nil,
						nil,
						// ldap config
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
			},
			object: &IDPTemplate{
				CreationDate:      testNow,
				ChangeDate:        testNow,
				Sequence:          20211109,
				ResourceOwner:     "ro",
				ID:                "idp-id",
				State:             domain.IDPStateActive,
				Name:              "idp-name",
				Type:              domain.IDPTypeLDAP,
				OwnerType:         domain.IdentityProviderTypeOrg,
				IsCreationAllowed: true,
				IsLinkingAllowed:  true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
			},
		},
		{
			name:    "prepareIDPTemplateByIDQuery sql err",
			prepare: prepareIDPTemplateByIDQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(idpTemplateQuery),
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
			name:    "prepareIDPTemplatesQuery no result",
			prepare: prepareIDPTemplatesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpTemplatesQuery),
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
			object: &IDPTemplates{Templates: []*IDPTemplate{}},
		},
		{
			name:    "prepareIDPTemplatesQuery ldap idp",
			prepare: prepareIDPTemplatesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpTemplatesQuery),
					idpTemplatesCols,
					[][]driver.Value{
						{
							"idp-id",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPTypeLDAP,
							domain.IdentityProviderTypeOrg,
							true,
							true,
							true,
							true,
							// google config
							nil,
							nil,
							nil,
							nil,
							// ldap config
							"idp-id",
							"host",
							"port",
							true,
							"base",
							"user",
							"uid",
							"admin",
							nil,
							"id",
							"first",
							"last",
							"display",
							"nickname",
							"username",
							"email",
							"emailVerified",
							"phone",
							"phoneVerified",
							"lang",
							"avatar",
							"profile",
						},
					},
				),
			},
			object: &IDPTemplates{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Templates: []*IDPTemplate{
					{
						CreationDate:      testNow,
						ChangeDate:        testNow,
						Sequence:          20211109,
						ResourceOwner:     "ro",
						ID:                "idp-id",
						State:             domain.IDPStateActive,
						Name:              "idp-name",
						Type:              domain.IDPTypeLDAP,
						OwnerType:         domain.IdentityProviderTypeOrg,
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						LDAPIDPTemplate: &LDAPIDPTemplate{
							IDPID:               "idp-id",
							Host:                "host",
							Port:                "port",
							TLS:                 true,
							BaseDN:              "base",
							UserObjectClass:     "user",
							UserUniqueAttribute: "uid",
							Admin:               "admin",
							LDAPAttributes: idp.LDAPAttributes{
								IDAttribute:                "id",
								FirstNameAttribute:         "first",
								LastNameAttribute:          "last",
								DisplayNameAttribute:       "display",
								NickNameAttribute:          "nickname",
								PreferredUsernameAttribute: "username",
								EmailAttribute:             "email",
								EmailVerifiedAttribute:     "emailVerified",
								PhoneAttribute:             "phone",
								PhoneVerifiedAttribute:     "phoneVerified",
								PreferredLanguageAttribute: "lang",
								AvatarURLAttribute:         "avatar",
								ProfileAttribute:           "profile",
							},
						},
					},
				},
			},
		},
		{
			name:    "prepareIDPTemplatesQuery no config",
			prepare: prepareIDPTemplatesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpTemplatesQuery),
					idpTemplatesCols,
					[][]driver.Value{
						{
							"idp-id",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPTypeLDAP,
							domain.IdentityProviderTypeOrg,
							true,
							true,
							true,
							true,
							// google config
							nil,
							nil,
							nil,
							nil,
							// ldap config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &IDPTemplates{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Templates: []*IDPTemplate{
					{
						CreationDate:      testNow,
						ChangeDate:        testNow,
						Sequence:          20211109,
						ResourceOwner:     "ro",
						ID:                "idp-id",
						State:             domain.IDPStateActive,
						Name:              "idp-name",
						Type:              domain.IDPTypeLDAP,
						OwnerType:         domain.IdentityProviderTypeOrg,
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
		},
		{
			name:    "prepareIDPTemplatesQuery all config types",
			prepare: prepareIDPTemplatesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(idpTemplatesQuery),
					idpTemplatesCols,
					[][]driver.Value{
						{
							"idp-id-ldap",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPTypeLDAP,
							domain.IdentityProviderTypeOrg,
							true,
							true,
							true,
							true,
							// google config
							nil,
							nil,
							nil,
							nil,
							// ldap config
							"idp-id-ldap",
							"host",
							"port",
							true,
							"base",
							"user",
							"uid",
							"admin",
							nil,
							"id",
							"first",
							"last",
							"display",
							"nickname",
							"username",
							"email",
							"emailVerified",
							"phone",
							"phoneVerified",
							"lang",
							"avatar",
							"profile",
						},
						{
							"idp-id-google",
							"ro",
							testNow,
							testNow,
							uint64(20211109),
							domain.IDPConfigStateActive,
							"idp-name",
							domain.IDPTypeGoogle,
							domain.IdentityProviderTypeOrg,
							true,
							true,
							true,
							true,
							// google
							"idp-id-google",
							"client_id",
							nil,
							database.StringArray{"profile"},
							// ldap config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &IDPTemplates{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Templates: []*IDPTemplate{
					{
						CreationDate:      testNow,
						ChangeDate:        testNow,
						Sequence:          20211109,
						ResourceOwner:     "ro",
						ID:                "idp-id-ldap",
						State:             domain.IDPStateActive,
						Name:              "idp-name",
						Type:              domain.IDPTypeLDAP,
						OwnerType:         domain.IdentityProviderTypeOrg,
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						LDAPIDPTemplate: &LDAPIDPTemplate{
							IDPID:               "idp-id-ldap",
							Host:                "host",
							Port:                "port",
							TLS:                 true,
							BaseDN:              "base",
							UserObjectClass:     "user",
							UserUniqueAttribute: "uid",
							Admin:               "admin",
							LDAPAttributes: idp.LDAPAttributes{
								IDAttribute:                "id",
								FirstNameAttribute:         "first",
								LastNameAttribute:          "last",
								DisplayNameAttribute:       "display",
								NickNameAttribute:          "nickname",
								PreferredUsernameAttribute: "username",
								EmailAttribute:             "email",
								EmailVerifiedAttribute:     "emailVerified",
								PhoneAttribute:             "phone",
								PhoneVerifiedAttribute:     "phoneVerified",
								PreferredLanguageAttribute: "lang",
								AvatarURLAttribute:         "avatar",
								ProfileAttribute:           "profile",
							},
						},
					},
					{
						CreationDate:      testNow,
						ChangeDate:        testNow,
						Sequence:          20211109,
						ResourceOwner:     "ro",
						ID:                "idp-id-google",
						State:             domain.IDPStateActive,
						Name:              "idp-name",
						Type:              domain.IDPTypeGoogle,
						OwnerType:         domain.IdentityProviderTypeOrg,
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						GoogleIDPTemplate: &GoogleIDPTemplate{
							IDPID:        "idp-id-google",
							ClientID:     "client_id",
							ClientSecret: nil,
							Scopes:       []string{"profile"},
						},
					},
				},
			},
		},
		{
			name:    "prepareIDPTemplatesQuery sql err",
			prepare: prepareIDPTemplatesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(idpTemplatesQuery),
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
