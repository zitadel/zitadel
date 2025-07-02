package resources

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/test"
)

func TestOperationCollection_Apply(t *testing.T) {
	tests := []struct {
		name              string
		op                *patch.Operation
		prepare           func(user *ScimUser)
		want              *ScimUser
		wantFn            func(user *ScimUser)
		wantModifications []string
		wantErr           bool
	}{
		{
			name: "add unknown path",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("fooBar")),
				Value:     json.RawMessage(`{ "userName": "hans.muster", "nickname": "hansi" }`),
			},
			wantErr: true,
		},
		{
			name: "add without path",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Value:     json.RawMessage(`{ "userName": "hans.muster", "nickname": "hansi" }`),
			},
			want: &ScimUser{
				UserName: "hans.muster",
				NickName: "hansi",
			},
		},
		{
			name: "add without path (sample from rfc)",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Value: json.RawMessage(`
					{
						"emails":[
							{
							  "value":"babs@jensen.org",
							  "type":"home",
                              "primary": true
							}
						],
						"nickname":"Babs"
					}`),
			},
			want: &ScimUser{
				NickName: "Babs",
				Emails: []*ScimEmail{
					{
						Value:   "jeanie.pendleton@example.com",
						Primary: false,
					},
					{
						Value:   "babs@jensen.org",
						Primary: true,
					},
				},
			},
		},
		{
			name: "add with path value which is nil",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("externalid")),
				Value:     json.RawMessage(`"externalid-1"`),
			},
			want: &ScimUser{
				ExternalID: "externalid-1",
			},
			wantModifications: []string{"rep:externalid"},
		},
		{
			name: "add complex attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("name")),
				Value: json.RawMessage(`{
					"Formatted": "added-formatted",
					"FamilyName": "added-family-name",
					"GivenName": "added-given-name",
					"MiddleName": "added-middle-name",
					"HonorificPrefix": "added-honorific-prefix",
					"HonorificSuffix": "added-honorific-suffix"
				}`),
			},
			want: &ScimUser{
				Name: &ScimUserName{
					Formatted:       "added-formatted",
					FamilyName:      "added-family-name",
					GivenName:       "added-given-name",
					MiddleName:      "added-middle-name",
					HonorificPrefix: "added-honorific-prefix",
					HonorificSuffix: "added-honorific-suffix",
				},
			},
			wantModifications: []string{"rep:name"},
		},
		{
			name: "add complex attribute value",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("name.formatted")),
				Value:     json.RawMessage(`"added-formatted"`),
			},
			want: &ScimUser{
				Name: &ScimUserName{
					Formatted: "added-formatted",
				},
			},
			wantModifications: []string{"rep:name.formatted"},
		},
		{
			name: "add single to multi valued empty attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`{
					"value": "added-entitlement"
				}`),
			},
			prepare: func(user *ScimUser) {
				user.Entitlements = nil
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "added-entitlement",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add single to multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`{
					"value": "added-entitlement"
				}`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
					{
						Value:   "added-entitlement",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add single primary to multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`{
					"value": "added-entitlement",
					"primary": true
				}`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: false,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
					{
						Value:   "added-entitlement",
						Primary: true,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add unique valued item in multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`
					{
						"value": "my-entitlement-1",
						"display": "entitlement-1-patched",
						"primary": true
					}`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "entitlement-1-patched",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add unique valued item and additional item in multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`[
						{
							"value": "my-entitlement-1",
							"display": "entitlement-1-patched",
							"primary": true
						},
						{
							"value": "my-entitlement-3",
							"display": "entitlement-3",
							"primary": false
						}
					]`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "entitlement-1-patched",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
					{
						Value:   "my-entitlement-3",
						Display: "entitlement-3",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add unique valued items in multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`[
						{
							"value": "my-entitlement-1",
							"display": "entitlement-1-patched",
							"primary": true
						},
						{
							"value": "my-entitlement-2",
							"display": "entitlement-2-patched",
							"primary": false
						}
					]`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "entitlement-1-patched",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "entitlement-2-patched",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add multiple to multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(` [
					{
						"value": "added-entitlement",
						"primary": true
					},
					{
						"value": "added-entitlement-2"
					}
				]`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: false,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
					{
						Value:   "added-entitlement",
						Primary: true,
					},
					{
						Value:   "added-entitlement-2",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"add:entitlements"},
		},
		{
			name: "add multiple primaries to multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeAdd,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(` [
					{
						"value": "added-entitlement",
						"primary": true
					},
					{
						"value": "added-entitlement-2",
						"primary": true
					}
				]`),
			},
			wantErr: true,
		},
		{
			name: "remove unknown path",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("fooBar")),
			},
			wantErr: true,
		},
		{
			name: "remove without a path",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
			},
			wantErr: true,
		},
		{
			name: "remove single valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath("nickname")),
			},
			wantFn: func(user *ScimUser) {
				assert.Equal(t, user.NickName, "")
			},
			wantModifications: []string{"rem:nickname"},
		},
		{
			name: "remove multi valued attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath("entitlements")),
			},
			wantFn: func(user *ScimUser) {
				assert.Len(t, user.Entitlements, 0)
			},
			wantModifications: []string{"rem:entitlements"},
		},
		{
			name: "remove multi valued attribute with filter",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath(`entitlements[display ew "1"]`)),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"rem:entitlements"},
		},
		{
			name: "remove multi valued attribute with filter matches all",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath(`entitlements[value pr]`)),
			},
			wantFn: func(user *ScimUser) {
				assert.Len(t, user.Entitlements, 0)
			},
			wantModifications: []string{"rem:entitlements"},
		},
		{
			name: "remove attribute of multi valued attribute with filter",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath(`entitlements[display ew "1"].display`)),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "",
						Type:    "main-entitlement",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			},
			wantModifications: []string{"rem:entitlements.display"},
		},
		{
			// this should not fail
			// according to the rfc only the replace operation fails without a matching target
			name: "remove multi valued attribute with filter and no matches",
			op: &patch.Operation{
				Operation: patch.OperationTypeRemove,
				Path:      test.Must(filter.ParsePath(`entitlements[display eq "FOOBAR"]`)),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			},
		},
		{
			name: "replace unknown path",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("fooBar")),
				Value:     json.RawMessage(`"fooBar"`),
			},
			wantErr: true,
		},
		{
			name: "replace without path",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Value:     json.RawMessage(`{ "userName": "hans.muster", "nickname": "hansi" }`),
			},
			want: &ScimUser{
				UserName: "hans.muster",
				NickName: "hansi",
			},
		},
		{
			name: "replace single valued attribute with path",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("nickname")),
				Value:     json.RawMessage(`"fooBar"`),
			},
			want: &ScimUser{
				NickName: "fooBar",
			},
		},
		{
			name: "replace with path value which is nil",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("externalid")),
				Value:     json.RawMessage(`"externalid-1"`),
			},
			want: &ScimUser{
				ExternalID: "externalid-1",
			},
			wantModifications: []string{"rep:externalid"},
		},
		{
			name: "replace complex attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("name")),
				Value: json.RawMessage(`{
					"Formatted": "added-formatted",
					"FamilyName": "added-family-name",
					"GivenName": "added-given-name",
					"MiddleName": "added-middle-name",
					"HonorificPrefix": "added-honorific-prefix",
					"HonorificSuffix": "added-honorific-suffix"
				}`),
			},
			want: &ScimUser{
				Name: &ScimUserName{
					Formatted:       "added-formatted",
					FamilyName:      "added-family-name",
					GivenName:       "added-given-name",
					MiddleName:      "added-middle-name",
					HonorificPrefix: "added-honorific-prefix",
					HonorificSuffix: "added-honorific-suffix",
				},
			},
			wantModifications: []string{"rep:name"},
		},
		{
			name: "replace complex multi attribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`{
					"value": "entitlement patched"
				}`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value: "entitlement patched",
					},
				},
			},
			wantModifications: []string{"rep:entitlements"},
		},
		{
			name: "replace complex multi attribute with multiple",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`[
					{
						"value": "entitlement patched"
					},
					{
						"value": "entitlement patched2"
					}
				]`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value: "entitlement patched",
					},
					{
						Value: "entitlement patched2",
					},
				},
			},
			wantModifications: []string{"rep:entitlements"},
		},
		{
			name: "replace complex multi attribute with multiple primary",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath("entitlements")),
				Value: json.RawMessage(`[
					{
						"value": "entitlement patched",
						"primary": true
					},
					{
						"value": "entitlement patched2",
						"primary": true
					}
				]`),
			},
			wantErr: true,
		},
		{
			name: "replace filter no match",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`entitlements[value eq "foobar"]`)),
			},
			wantErr: true,
		},
		{
			name: "replace filter complex subattribute",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`entitlements[value eq "my-entitlement-1"].display`)),
				Value:     json.RawMessage(`"updated display"`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "updated display",
						Type:    "main-entitlement",
						Primary: true,
					},
				},
			},
			wantModifications: []string{"rep:entitlements.display"},
		},
		{
			name: "replace filter complex subattribute primary",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`entitlements[value eq "my-entitlement-2"].primary`)),
				Value:     json.RawMessage(`true`),
			},
			want: &ScimUser{
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: false,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: true,
					},
				},
			},
			wantModifications: []string{"rep:entitlements.primary"},
		},
		{
			name: "replace filter complex subattribute multiple primary",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`roles[primary ne true].primary`)),
				Value:     json.RawMessage(`true`),
			},
			wantErr: true,
		},
		{
			name: "replace filter complex subattribute multiple emails primary value",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`emails[primary eq true].value`)),
				Value:     json.RawMessage(`"jeanie.rebecca.pendleton@example.com"`),
			},
			want: &ScimUser{
				Emails: []*ScimEmail{
					{
						Value:   "jeanie.rebecca.pendleton@example.com",
						Primary: true,
					},
				},
			},
		},
		{
			name: "replace filter complex subattribute multiple emails type value",
			op: &patch.Operation{
				Operation: patch.OperationTypeReplace,
				Path:      test.Must(filter.ParsePath(`emails[type eq "work"].value`)),
				Value:     json.RawMessage(`"jeanie.rebecca.pendleton@example.com"`),
			},
			want: &ScimUser{
				Emails: []*ScimEmail{
					{
						Value:   "jeanie.rebecca.pendleton@example.com",
						Primary: true,
						Type:    "work",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &ScimUser{
				ID:       "1",
				UserName: "username-1",
				Name: &ScimUserName{
					Formatted:       "Dr. Jeanie R. Pendleton III",
					FamilyName:      "Pendleton",
					GivenName:       "Jeanie",
					MiddleName:      "Rebecca",
					HonorificPrefix: "Dr.",
					HonorificSuffix: "III",
				},
				DisplayName:       "Jeanie Pendleton",
				NickName:          "Jenny",
				Title:             "Ms.",
				ProfileUrl:        test.Must(schemas.ParseHTTPURL("https://example.com/profile.gif")),
				PreferredLanguage: language.MustParse("en-US"),
				Locale:            "en-US",
				Timezone:          "America/New_York",
				Active:            schemas.NewRelaxedBool(true),
				Emails: []*ScimEmail{
					{
						Value:   "jeanie.pendleton@example.com",
						Primary: true,
						Type:    "work",
					},
				},
				PhoneNumbers: []*ScimPhoneNumber{
					{
						Value:   "+1 775-599-5252",
						Primary: true,
					},
				},
				Ims: []*ScimIms{
					{
						Value: "jeeeeeny91",
						Type:  "icq",
					},
				},
				Addresses: []*ScimAddress{
					{
						Type:          "work",
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza\nHollywood, CA 91608 USA",
						Primary:       true,
					},
				},
				Photos: []*ScimPhoto{
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/F")),
						Type:  "photo",
					},
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/T")),
						Type:  "thumbnail",
					},
				},
				Roles: []*ScimRole{
					{
						Value:   "my-role-1",
						Display: "Rolle 1",
						Type:    "main-role",
						Primary: true,
					},
					{
						Value:   "my-role-2",
						Display: "Rolle 2",
						Type:    "secondary-role",
						Primary: false,
					},
					{
						Value:   "my-role-3",
						Display: "Rolle 3",
						Type:    "third-role",
						Primary: false,
					},
				},
				Entitlements: []*ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			}

			if tt.prepare != nil {
				tt.prepare(user)
			}

			patcher := new(simplePatcher)
			operations := patch.OperationCollection{tt.op}
			err := operations.Apply(patcher, user)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			if !test.PartiallyDeepEqual(tt.want, user) {
				t.Errorf("apply() got = %#v, want %#v", user, tt.want)
			}

			if tt.wantModifications != nil {
				assert.EqualValues(t, tt.wantModifications, patcher.modifications)
			}

			if tt.wantFn != nil {
				tt.wantFn(user)
			}
		})
	}
}

func Test_userPatcher_updateMetadata(t *testing.T) {
	tests := []struct {
		name                     string
		metadataPath             []string
		metadataChanges          map[metadata.Key]*domain.Metadata
		metadataKeysToRemove     map[metadata.Key]bool
		wantErr                  bool
		wantMetadataChanges      map[metadata.Key]*domain.Metadata
		wantMetadataKeysToRemove map[metadata.Key]bool
	}{
		{
			name: "empty path",
		},
		{
			name:         "unknown attribute",
			metadataPath: []string{"fooBar"},
		},
		{
			name:         "unknown nested attribute",
			metadataPath: []string{"foo", "bar"},
		},
		{
			name:         "simple attribute",
			metadataPath: []string{"title"},
			wantMetadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyTitle: {
					Key:   string(metadata.KeyTitle),
					Value: []byte("Mr."),
				},
			},
		},
		{
			name:         "simple attribute with previous deletion",
			metadataPath: []string{"title"},
			metadataKeysToRemove: map[metadata.Key]bool{
				metadata.KeyTitle: true,
			},
			wantMetadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyTitle: {
					Key:   string(metadata.KeyTitle),
					Value: []byte("Mr."),
				},
			},
		},
		{
			name:         "nested attribute",
			metadataPath: []string{"name", "middlename"},
			wantMetadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyMiddleName: {
					Key:   string(metadata.KeyMiddleName),
					Value: []byte("middle name"),
				},
			},
		},
		{
			name:         "nested attribute with previous deletion",
			metadataPath: []string{"name", "middlename"},
			metadataKeysToRemove: map[metadata.Key]bool{
				metadata.KeyMiddleName: true,
			},
			wantMetadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyMiddleName: {
					Key:   string(metadata.KeyMiddleName),
					Value: []byte("middle name"),
				},
			},
		},
		{
			name:         "complex object root",
			metadataPath: []string{"name"},
			wantMetadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyMiddleName: {
					Key:   string(metadata.KeyMiddleName),
					Value: []byte("middle name"),
				},
				metadata.KeyHonorificPrefix: {
					Key:   string(metadata.KeyHonorificPrefix),
					Value: []byte("Adel"),
				},
				metadata.KeyHonorificSuffix: {
					Key:   string(metadata.KeyHonorificSuffix),
					Value: []byte("III"),
				},
			},
		},
		{
			name:         "delete previous modified attribute",
			metadataPath: []string{"locale"},
			metadataChanges: map[metadata.Key]*domain.Metadata{
				metadata.KeyLocale: {
					Key:   string(metadata.KeyLocale),
					Value: []byte("edited locale"),
				},
			},
			wantMetadataKeysToRemove: map[metadata.Key]bool{
				metadata.KeyLocale: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &userPatcher{
				ctx: context.Background(),
				user: &ScimUser{
					Title: "Mr.",
					Name: &ScimUserName{
						MiddleName:      "middle name",
						HonorificPrefix: "Adel",
						HonorificSuffix: "III",
					},
				},
				metadataChanges:      tt.metadataChanges,
				metadataKeysToRemove: tt.metadataKeysToRemove,
			}

			if p.metadataChanges == nil {
				p.metadataChanges = make(map[metadata.Key]*domain.Metadata)
			}

			if p.metadataKeysToRemove == nil {
				p.metadataKeysToRemove = make(map[metadata.Key]bool)
			}

			err := p.updateMetadata(tt.metadataPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantMetadataChanges != nil {
				if !reflect.DeepEqual(p.metadataChanges, tt.wantMetadataChanges) {
					t.Errorf("updateMetadata() got = %#v, want %#v", p.metadataChanges, tt.wantMetadataChanges)
				}
			}

			if tt.metadataKeysToRemove != nil {
				if !reflect.DeepEqual(p.metadataKeysToRemove, tt.metadataKeysToRemove) {
					t.Errorf("updateMetadata() got = %#v, want %#v", p.metadataKeysToRemove, tt.metadataKeysToRemove)
				}
			}
		})
	}
}

type simplePatcher struct {
	modifications []string
}

func (s *simplePatcher) FilterEvaluator() *filter.Evaluator {
	return filter.NewEvaluator(schemas.IdUser)
}

func (s *simplePatcher) Added(attributePath []string) error {
	return s.modified("add", attributePath)
}

func (s *simplePatcher) Replaced(attributePath []string) error {
	return s.modified("rep", attributePath)
}

func (s *simplePatcher) Removed(attributePath []string) error {
	return s.modified("rem", attributePath)
}

func (s *simplePatcher) modified(op string, attributePath []string) error {
	s.modifications = append(s.modifications, op+":"+strings.Join(attributePath, "."))
	return nil
}
