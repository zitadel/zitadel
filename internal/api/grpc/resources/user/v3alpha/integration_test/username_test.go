//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_AddUsername(t *testing.T) {
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
		dep     func(req *user.AddUsernameRequest) error
		req     *user.AddUsernameRequest
		res     res
		wantErr bool
	}{
		{
			name: "username add, no context",
			ctx:  context.Background(),
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, username empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					Username:               "",
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, user not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
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
			name: "username add, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: false,
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
			name: "username add, already existing",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				username := gofakeit.Username()
				instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, username, false)
				req.Username.Username = username
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					IsOrganizationSpecific: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, isOrgSpecific, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					Username:               gofakeit.Username(),
					IsOrganizationSpecific: true,
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
			name: "username add, isOrgSpecific, already existing",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				username := gofakeit.Username()
				instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, username, true)
				req.Username.Username = username
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					IsOrganizationSpecific: true,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, isOrgSpecific existing, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				username := gofakeit.Username()
				instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, username, true)
				req.Username.Username = username
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					IsOrganizationSpecific: false,
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
			name: "username add, existing, isOrgSpecific ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				username := gofakeit.Username()
				instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, username, false)
				req.Username.Username = username
				return nil
			},
			req: &user.AddUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Username: &user.SetUsername{
					IsOrganizationSpecific: true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				assert.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.AddUsername(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}

func TestServer_DeleteUsername(t *testing.T) {
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
		dep     func(req *user.RemoveUsernameRequest) error
		req     *user.RemoveUsernameRequest
		res     res
		wantErr bool
	}{
		{
			name: "username delete, no context",
			ctx:  context.Background(),
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				usernameResp := instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), gofakeit.Username(), false)
				req.UsernameId = usernameResp.GetUsernameId()
				return nil
			},
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "username delete, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				usernameResp := instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), gofakeit.Username(), false)
				req.UsernameId = usernameResp.GetUsernameId()
				return nil
			},
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "username remove, id empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				UsernameId: "notempty",
			},
			wantErr: true,
		},
		{
			name: "username delete, userid empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemoveUsernameRequest{
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
			name: "username remove, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				UsernameId: "notempty",
				Id:         "notexisting",
			},
			wantErr: true,
		},
		{
			name: "username remove, no username",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "username remove, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				usernameResp := instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), gofakeit.Username(), false)
				req.UsernameId = usernameResp.GetUsernameId()
				return nil
			},
			req: &user.RemoveUsernameRequest{
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
			name: "username remove, already removed",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				resp := instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, gofakeit.Username(), false)
				req.UsernameId = resp.GetUsernameId()
				instance.RemoveAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, req.UsernameId)
				return nil
			},
			req: &user.RemoveUsernameRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "username remove, isOrgSpecific, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemoveUsernameRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				usernameResp := instance.AddAuthenticatorUsername(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), gofakeit.Username(), true)
				req.UsernameId = usernameResp.GetUsernameId()
				return nil
			},
			req: &user.RemoveUsernameRequest{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				assert.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.RemoveUsername(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return

			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)

		})
	}
}
