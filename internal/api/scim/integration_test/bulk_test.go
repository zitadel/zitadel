//go:build integration

package integration_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/internal/test"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

var (
	//go:embed testdata/bulk_test_full.json
	bulkFullJson []byte

	//go:embed testdata/bulk_test_fail_on_errors.json
	bulkFailOnErrorsJson []byte

	//go:embed testdata/bulk_test_errors.json
	bulkErrorsFullJson []byte

	bulkTooManyOperationsJson []byte
)

func init() {
	bulkFullJson = removeComments(bulkFullJson)
	bulkFailOnErrorsJson = removeComments(bulkFailOnErrorsJson)
	bulkErrorsFullJson = removeComments(bulkErrorsFullJson)
	bulkTooManyOperationsJson = test.Must(json.Marshal(buildTooManyOperationsRequest()))
}

func TestBulk(t *testing.T) {
	iamOwnerCtx := Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	secondaryOrg := Instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	createdSecondaryOrgUser := createHumanUser(t, iamOwnerCtx, secondaryOrg.OrganizationId, 0)
	bulkMinimalUpdateSecondaryOrgJson := test.Must(json.Marshal(buildMinimalUpdateRequest(createdSecondaryOrgUser.UserId)))

	membershipNotFoundErr := &scim.ScimError{
		Schemas: []string{
			"urn:ietf:params:scim:api:messages:2.0:Error",
			"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
		},
		Detail: "membership not found",
		Status: "404",
		ZitadelDetail: &scim.ZitadelErrorDetail{
			ID:      "AUTHZ-cdgFk",
			Message: "membership not found",
		},
	}

	type wantErr struct {
		scimErrorType string
		status        int
		zitadelErrID  string
	}

	tests := []struct {
		name      string
		body      []byte
		ctx       context.Context
		orgID     string
		want      *scim.BulkResponse
		wantErr   *wantErr
		wantUsers map[string]*resources.ScimUser
	}{
		{
			name: "not authenticated",
			body: bulkFullJson,
			ctx:  context.Background(),
			wantErr: &wantErr{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "no permissions",
			body: bulkFullJson,
			ctx:  Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			want: &scim.BulkResponse{
				Schemas: []schemas.ScimSchemaType{schemas.IdBulkResponse},
				Operations: []*scim.BulkResponseOperation{
					{
						Method:   http.MethodPost,
						Response: membershipNotFoundErr,
						Status:   "404",
					},
					{
						Method:   http.MethodPost,
						BulkID:   "1",
						Response: membershipNotFoundErr,
						Status:   "404",
					},
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail:   "Could not resolve bulkID 1 to created ID",
							Status:   "400",
							ScimType: "invalidValue",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-BLK4",
								Message: "Could not resolve bulkID 1 to created ID",
							},
						},
						Status: "400",
					},
					{
						Method:   http.MethodPost,
						BulkID:   "2",
						Response: membershipNotFoundErr,
						Status:   "404",
					},
					{
						Method: http.MethodPut,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail:   "Could not resolve bulkID 2 to created ID",
							Status:   "400",
							ScimType: "invalidValue",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-BLK4",
								Message: "Could not resolve bulkID 2 to created ID",
							},
						},
						Status: "400",
					},
					{
						Method:   http.MethodPost,
						BulkID:   "3",
						Response: membershipNotFoundErr,
						Status:   "404",
					},
					{
						Method: http.MethodDelete,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail:   "Could not resolve bulkID 3 to created ID",
							Status:   "400",
							ScimType: "invalidValue",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-BLK4",
								Message: "Could not resolve bulkID 3 to created ID",
							},
						},
						Status: "400",
					},
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail: "User could not be found",
							Status: "404",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "COMMAND-ugjs0upun6",
								Message: "Errors.User.NotFound",
							},
						},
						Status: "404",
					},
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail:   "Could not resolve bulkID 99 to created ID",
							Status:   "400",
							ScimType: "invalidValue",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-BLK4",
								Message: "Could not resolve bulkID 99 to created ID",
							},
						},
						Status: "400",
					},
				},
			},
		},
		{
			name: "full",
			body: bulkFullJson,
			want: &scim.BulkResponse{
				Schemas: []schemas.ScimSchemaType{schemas.IdBulkResponse},
				Operations: []*scim.BulkResponseOperation{
					{
						Method: http.MethodPost,
						Status: "201",
					},
					{
						Method: http.MethodPost,
						BulkID: "1",
						Status: "201",
					},
					{
						Method: http.MethodPatch,
						Status: "204",
					},
					{
						Method: http.MethodPost,
						BulkID: "2",
						Status: "201",
					},
					{
						Method: http.MethodPut,
						Status: "200",
					},
					{
						Method: http.MethodPost,
						BulkID: "3",
						Status: "201",
					},
					{
						Method: http.MethodDelete,
						Status: "204",
					},
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail: "User could not be found",
							Status: "404",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "COMMAND-ugjs0upun6",
								Message: "Errors.User.NotFound",
							},
						},
						Status: "404",
					},
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail:   "Could not resolve bulkID 99 to created ID",
							Status:   "400",
							ScimType: "invalidValue",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-BLK4",
								Message: "Could not resolve bulkID 99 to created ID",
							},
						},
						Status: "400",
					},
				},
			},
			wantUsers: map[string]*resources.ScimUser{
				"scim-bulk-created-user-0": {
					ExternalID: "scim-bulk-created-user-0",
					UserName:   "scim-bulk-created-user-0",
					Name: &resources.ScimUserName{
						Formatted:  "scim-bulk-created-user-0-given-name scim-bulk-created-user-0-family-name",
						FamilyName: "scim-bulk-created-user-0-family-name",
						GivenName:  "scim-bulk-created-user-0-given-name",
					},
					DisplayName:       "scim-bulk-created-user-0-given-name scim-bulk-created-user-0-family-name",
					PreferredLanguage: test.Must(language.Parse("en")),
					Active:            schemas.NewRelaxedBool(true),
					Emails: []*resources.ScimEmail{
						{
							Value:   "scim-bulk-created-user-0@example.com",
							Primary: true,
						},
					},
				},
				"scim-bulk-created-user-1": {
					ExternalID: "scim-bulk-created-user-1",
					UserName:   "scim-bulk-created-user-1",
					Name: &resources.ScimUserName{
						Formatted:  "scim-bulk-created-user-1-given-name scim-bulk-created-user-1-family-name",
						FamilyName: "scim-bulk-created-user-1-family-name",
						GivenName:  "scim-bulk-created-user-1-given-name",
					},
					DisplayName:       "scim-bulk-created-user-1-given-name scim-bulk-created-user-1-family-name",
					NickName:          "scim-bulk-created-user-1-nickname-patched",
					PreferredLanguage: test.Must(language.Parse("en")),
					Active:            schemas.NewRelaxedBool(true),
					Emails: []*resources.ScimEmail{
						{
							Value:   "scim-bulk-created-user-1@example.com",
							Primary: true,
						},
					},
					PhoneNumbers: []*resources.ScimPhoneNumber{
						{
							Value:   "+41711231212",
							Primary: true,
						},
					},
				},
				"scim-bulk-created-user-2": {
					ExternalID: "scim-bulk-created-user-2",
					UserName:   "scim-bulk-created-user-2",
					Name: &resources.ScimUserName{
						Formatted:  "scim-bulk-created-user-2-given-name scim-bulk-created-user-2-family-name",
						FamilyName: "scim-bulk-created-user-2-family-name",
						GivenName:  "scim-bulk-created-user-2-given-name",
					},
					DisplayName:       "scim-bulk-created-user-2-given-name scim-bulk-created-user-2-family-name",
					NickName:          "scim-bulk-created-user-2-nickname-patched",
					PreferredLanguage: test.Must(language.Parse("en")),
					Active:            schemas.NewRelaxedBool(true),
					Emails: []*resources.ScimEmail{
						{
							Value:   "scim-bulk-created-user-2@example.com",
							Primary: true,
						},
					},
					PhoneNumbers: []*resources.ScimPhoneNumber{
						{
							Value:   "+41711231212",
							Primary: true,
						},
					},
				},
			},
		},
		{
			name: "errors",
			body: bulkErrorsFullJson,
			want: &scim.BulkResponse{
				Schemas: []schemas.ScimSchemaType{schemas.IdBulkResponse},
				Operations: []*scim.BulkResponseOperation{
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail: "User could not be found",
							Status: "404",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "COMMAND-ugjs0upun6",
								Message: "Errors.User.NotFound",
							},
						},
						Status: "404",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Email is empty",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-EM19",
								Message: "Errors.User.Email.Empty",
							},
						},
						Status: "400",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Could not parse locale",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-MD11",
								Message: "Could not parse locale",
							},
						},
						Status: "400",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Password is too short",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "COMMA-HuJf6",
								Message: "Errors.User.PasswordComplexityPolicy.MinLength",
							},
						},
						Status: "400",
					},
					{
						Method: "POST",
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Could not parse timezone",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-MD12",
								Message: "Could not parse timezone",
							},
						},
						Status: "400",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Errors.Invalid.Argument",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "V2-zzad3",
								Message: "Errors.Invalid.Argument",
							},
						},
						Status: "400",
					},
					{
						Method: "POST",
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Given name in profile is empty",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "USER-UCej2",
								Message: "Errors.User.Profile.FirstNameEmpty",
							},
						},
						Status: "400",
					},
				},
			},
		},
		{
			name: "fail on errors",
			body: withUsername(bulkFailOnErrorsJson, integration.Username()),
			want: &scim.BulkResponse{
				Schemas: []schemas.ScimSchemaType{schemas.IdBulkResponse},
				Operations: []*scim.BulkResponseOperation{
					{
						Method: http.MethodPatch,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							Detail: "User could not be found",
							Status: "404",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "COMMAND-ugjs0upun6",
								Message: "Errors.User.NotFound",
							},
						},
						Status: "404",
					},
					{
						Method: http.MethodPost,
						Status: "201",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Email is empty",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-EM19",
								Message: "Errors.User.Email.Empty",
							},
						},
						Status: "400",
					},
					{
						Method: http.MethodPost,
						Response: &scim.ScimError{
							Schemas: []string{
								"urn:ietf:params:scim:api:messages:2.0:Error",
								"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail",
							},
							ScimType: "invalidValue",
							Detail:   "Could not parse locale",
							Status:   "400",
							ZitadelDetail: &scim.ZitadelErrorDetail{
								ID:      "SCIM-MD11",
								Message: "Could not parse locale",
							},
						},
						Status: "400",
					},
				},
			},
		},
		{
			name: "too many operations",
			body: bulkTooManyOperationsJson,
			wantErr: &wantErr{
				status:        http.StatusRequestEntityTooLarge,
				scimErrorType: "invalidValue",
				zitadelErrID:  "SCIM-BLK19",
			},
		},
		{
			name:  "another organization",
			body:  bulkMinimalUpdateSecondaryOrgJson,
			orgID: secondaryOrg.OrganizationId,
			want: &scim.BulkResponse{
				Schemas: []schemas.ScimSchemaType{schemas.IdBulkResponse},
				Operations: []*scim.BulkResponseOperation{
					{
						Method:   http.MethodPatch,
						Response: membershipNotFoundErr,
						Status:   "404",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			orgID := tt.orgID
			if orgID == "" {
				orgID = Instance.DefaultOrg.Id
			}

			response, err := Instance.Client.SCIM.Bulk(ctx, orgID, tt.body)
			createdUserIDs := buildCreatedIDs(response)

			if tt.wantErr != nil {
				statusCode := tt.wantErr.status
				if statusCode == 0 {
					statusCode = http.StatusBadRequest
				}

				scimErr := scim.RequireScimError(t, statusCode, err)
				assert.Equal(t, tt.wantErr.scimErrorType, scimErr.Error.ScimType)

				if tt.wantErr.zitadelErrID != "" {
					assert.Equal(t, tt.wantErr.zitadelErrID, scimErr.Error.ZitadelDetail.ID)
				}
				return
			}

			require.NoError(t, err)
			require.EqualValues(t, []schemas.ScimSchemaType{schemas.IdBulkResponse}, response.Schemas)

			locationPrefix := "http://" + Instance.Host() + path.Join(schemas.HandlerPrefix, orgID, "Users") + "/"
			for _, responseOperation := range response.Operations {
				// POST operations which result in an error don't expect a location
				if responseOperation.Method == http.MethodPost && responseOperation.Response != nil {
					require.Empty(t, responseOperation.Location)
				} else {
					require.True(t, strings.HasPrefix(responseOperation.Location, locationPrefix))
				}

				// don't assert the location in the deep equal
				responseOperation.Location = ""
			}

			if !reflect.DeepEqual(tt.want, response) {
				x := test.Must(json.Marshal(tt.want))
				x2 := test.Must(json.Marshal(response))
				t.Errorf("want: %v, got: %v", x, x2)
				t.Errorf("want: %#v, got: %#v", tt.want, response)
			}

			if tt.wantUsers != nil {
				for _, createdUserID := range createdUserIDs {
					retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
					require.EventuallyWithT(t, func(ttt *assert.CollectT) {
						user, err := Instance.Client.SCIM.Users.Get(ctx, orgID, createdUserID)
						if err != nil {
							scim.RequireScimError(ttt, http.StatusNotFound, err)
							return
						}

						wantUser, ok := tt.wantUsers[user.UserName]
						if !ok {
							return
						}

						if !test.PartiallyDeepEqual(wantUser, user) {
							ttt.Errorf("want: %#v, got: %#v", wantUser, user)
						}
					}, retryDuration, tick)
				}
			}
		})
	}
}

func buildCreatedIDs(response *scim.BulkResponse) []string {
	createdIds := make([]string, 0, len(response.Operations))
	for _, operation := range response.Operations {
		if operation.Method == http.MethodPost && operation.Status == "201" {
			parts := strings.Split(operation.Location, "/")
			createdIds = append(createdIds, parts[len(parts)-1])
		}
	}

	return createdIds
}

func buildMinimalUpdateRequest(userID string) *scim.BulkRequest {
	return &scim.BulkRequest{
		Schemas: []schemas.ScimSchemaType{schemas.IdBulkRequest},
		Operations: []*scim.BulkRequestOperation{
			{
				Method: http.MethodPatch,
				Path:   "/Users/" + userID,
				Data:   simpleReplacePatchBody("nickname", `"foo-bar-nickname"`),
			},
		},
	}
}

func buildTooManyOperationsRequest() *scim.BulkRequest {
	req := &scim.BulkRequest{
		Schemas:    []schemas.ScimSchemaType{schemas.IdBulkRequest},
		Operations: make([]*scim.BulkRequestOperation, 101), // default config (100) + 1, see defaults.yaml
	}

	for i := 0; i < len(req.Operations); i++ {
		req.Operations[i] = &scim.BulkRequestOperation{
			Method: http.MethodPost,
			Path:   "/Users",
			Data:   withUsername(minimalUserJson, integration.Username()),
		}
	}

	return req
}

func setProvisioningDomain(t require.TestingT, userID, provisioningDomain string) {
	setAndEnsureMetadata(t, userID, "urn:zitadel:scim:provisioningDomain", provisioningDomain)
}

func setAndEnsureMetadata(t require.TestingT, userID, key, value string) {
	_, err := Instance.Client.Mgmt.SetUserMetadata(CTX, &management.SetUserMetadataRequest{
		Id:    userID,
		Key:   key,
		Value: []byte(value),
	})
	require.NoError(t, err)

	// ensure metadata is projected
	ensureMetadataProjected(t, userID, key, value)
}

func ensureMetadataProjected(t require.TestingT, userID, key, value string) {
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.GetUserMetadata(CTX, &management.GetUserMetadataRequest{
			Id:  userID,
			Key: key,
		})
		if !assert.NoError(tt, err) {
			require.Equal(tt, status.Code(err), codes.NotFound)
			return
		}
		assert.Equal(tt, value, string(md.Metadata.Value))
	}, retryDuration, tick)
}

func removeProvisioningDomain(t require.TestingT, userID string) {
	_, err := Instance.Client.Mgmt.RemoveUserMetadata(CTX, &management.RemoveUserMetadataRequest{
		Id:  userID,
		Key: "urn:zitadel:scim:provisioningDomain",
	})
	require.NoError(t, err)
}
