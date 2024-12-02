//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_AddPersonalAccessToken(t *testing.T) {
	t.Parallel()
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, schema)
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	type res struct {
		want *resource_object.Details
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.AddPersonalAccessTokenRequest) error
		req     *user.AddPersonalAccessTokenRequest
		res     res
		wantErr bool
	}{
		{
			name: "pat add, no context",
			ctx:  context.Background(),
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			wantErr: true,
		},
		{
			name: "pat add, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			wantErr: true,
		},
		{
			name: "pat add, pat empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			wantErr: true,
		},
		{
			name: "pat add, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			wantErr: true,
		},
		{
			name: "pat add, user not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id:                  "notexisting",
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			wantErr: true,
		},
		{
			name: "pat add, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
			},
		},
		{
			name: "pat add, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				PersonalAccessToken: &user.SetPersonalAccessToken{},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
			},
		},
		{
			name: "pat add, expirationdate, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{
					ExpirationDate: timestamppb.New(time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)),
				},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
			},
		},
		{
			name: "pat add, expirationdate invalid",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessToken: &user.SetPersonalAccessToken{
					ExpirationDate: timestamppb.New(time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC)),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				assert.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.AddPersonalAccessToken(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}

func TestServer_DeletePersonalAccessToken(t *testing.T) {
	t.Parallel()
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, schema)
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	type res struct {
		want *resource_object.Details
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.RemovePersonalAccessTokenRequest) error
		req     *user.RemovePersonalAccessTokenRequest
		res     res
		wantErr bool
	}{
		{
			name: "pat delete, no context",
			ctx:  context.Background(),
			dep: func(req *user.RemovePersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				patResp := instance.AddAuthenticatorPersonalAccessToken(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId())
				req.PersonalAccessTokenId = patResp.GetPersonalAccessTokenId()
				return nil
			},
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pat delete, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.RemovePersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				patResp := instance.AddAuthenticatorPersonalAccessToken(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId())
				req.PersonalAccessTokenId = patResp.GetPersonalAccessTokenId()
				return nil
			},
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pat remove, id empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessTokenId: "notempty",
			},
			wantErr: true,
		},
		{
			name: "pat delete, userid empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notempty",
			},
			wantErr: true,
		},
		{
			name: "pat remove, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PersonalAccessTokenId: "notempty",
				Id:                    "notexisting",
			},
			wantErr: true,
		},
		{
			name: "pat remove, no pat",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pat remove, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				patResp := instance.AddAuthenticatorPersonalAccessToken(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId())
				req.PersonalAccessTokenId = patResp.GetPersonalAccessTokenId()
				return nil
			},
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
			},
		},
		{
			name: "pat remove, already removed",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePersonalAccessTokenRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				resp := instance.AddAuthenticatorPersonalAccessToken(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				req.PersonalAccessTokenId = resp.GetPersonalAccessTokenId()
				instance.RemoveAuthenticatorPersonalAccessToken(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, req.PersonalAccessTokenId)
				return nil
			},
			req: &user.RemovePersonalAccessTokenRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				assert.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.RemovePersonalAccessToken(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return

			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)

		})
	}
}
