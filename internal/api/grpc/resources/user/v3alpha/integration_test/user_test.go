//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_CreateUser(t *testing.T) {
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
	permissionSchema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"urn:zitadel:schema:permission": {
					"owner": "r",
					"self": "r"
				},
				"type": "string"
			}
		}
	}`)
	permissionSchemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, permissionSchema)
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	type res struct {
		want            *resource_object.Details
		returnCodeEmail bool
		returnCodePhone bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		req     *user.CreateUserRequest
		res     res
		wantErr bool
	}{
		{
			name: "user create, no schemaID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{Data: unmarshalJSON("{\"name\": \"user\"}")},
			},
			wantErr: true,
		},
		{
			name: "user create, no context",
			ctx:  context.Background(),
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: schemaResp.GetDetails().GetId(),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user create, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: schemaResp.GetDetails().GetId(),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user create, invalid schema permission, owner",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: permissionSchemaResp.GetDetails().GetId(),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user create, no user data",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: schemaResp.GetDetails().GetId(),
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
			name: "user create, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: schemaResp.GetDetails().GetId(),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
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
			name: "user create, full contact, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.CreateUser{
					SchemaId: schemaResp.GetDetails().GetId(),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
					Contact: &user.SetContact{
						Email: &user.SetEmail{
							Address:      integration.Email(),
							Verification: &user.SetEmail_ReturnCode{ReturnCode: &user.ReturnEmailVerificationCode{}},
						},
						Phone: &user.SetPhone{
							Number:       integration.Phone(),
							Verification: &user.SetPhone_ReturnCode{ReturnCode: &user.ReturnPhoneVerificationCode{}},
						},
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
				returnCodePhone: true,
				returnCodeEmail: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.UserV3Alpha.CreateUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
			if tt.res.returnCodeEmail {
				require.NotNil(t, got.EmailCode)
			}
			if tt.res.returnCodePhone {
				require.NotNil(t, got.PhoneCode)
			}

		})
	}
}

func TestServer_PatchUser(t *testing.T) {
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
	permissionSchema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"urn:zitadel:schema:permission": {
					"owner": "r",
					"self": "r"
				},
				"type": "string"
			}
		}
	}`)
	permissionSchemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, permissionSchema)
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	type res struct {
		want            *resource_object.Details
		returnCodeEmail bool
		returnCodePhone bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.PatchUserRequest) error
		req     *user.PatchUserRequest
		res     res
		wantErr bool
	}{
		{
			name: "user patch, no context",
			ctx:  context.Background(),
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user patch, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user patch, invalid schema permission, owner",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					SchemaId: gu.Ptr(permissionSchemaResp.GetDetails().GetId()),
					Data:     unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user patch, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user patch, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"user\"}"),
				},
			},
			wantErr: true,
		},
		{
			name: "user patch, no change",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				req.User.Data = unmarshalJSON(data)
				req.User.SchemaId = gu.Ptr(schemaID)
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{},
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
			name: "user patch, schema, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				changedSchemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, schema)
				req.User.SchemaId = gu.Ptr(changedSchemaResp.Details.Id)
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{},
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
			name: "user patch, schema and data, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				changedSchemaResp := instance.CreateUserSchema(isolatedIAMOwnerCTX, schema)
				req.User.SchemaId = gu.Ptr(changedSchemaResp.Details.Id)
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"changed\"}"),
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
			name: "user patch, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"changed\"}"),
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
			name: "user patch, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"changed\"}"),
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
			name: "user patch, contact email, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Contact: &user.SetContact{
						Email: &user.SetEmail{
							Address:      integration.Email(),
							Verification: &user.SetEmail_ReturnCode{ReturnCode: &user.ReturnEmailVerificationCode{}},
						},
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
				returnCodeEmail: true,
			},
		},
		{
			name: "user patch, contact phone, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Contact: &user.SetContact{
						Phone: &user.SetPhone{
							Number:       integration.Phone(),
							Verification: &user.SetPhone_ReturnCode{ReturnCode: &user.ReturnPhoneVerificationCode{}},
						},
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
				returnCodePhone: true,
			},
		},
		{
			name: "user patch, full contact, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.PatchUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"changed\"}"),
					Contact: &user.SetContact{
						Email: &user.SetEmail{
							Address:      integration.Email(),
							Verification: &user.SetEmail_ReturnCode{ReturnCode: &user.ReturnEmailVerificationCode{}},
						},
						Phone: &user.SetPhone{
							Number:       integration.Phone(),
							Verification: &user.SetPhone_ReturnCode{ReturnCode: &user.ReturnPhoneVerificationCode{}},
						},
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
				returnCodePhone: true,
				returnCodeEmail: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.PatchUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
			if tt.res.returnCodeEmail {
				assert.NotNil(t, got.EmailCode)
			} else {
				assert.Nil(t, got.EmailCode)
			}
			if tt.res.returnCodePhone {
				assert.NotNil(t, got.PhoneCode)
			} else {
				assert.Nil(t, got.PhoneCode)
			}
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
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
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.DeleteUserRequest) error
		req     *user.DeleteUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user delete, no userID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "user delete, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "user delete, not existing, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user delete, no context",
			ctx:  context.Background(),
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user delete, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user delete, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user delete, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeleteUserRequest{},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user delete, locked, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user delete, deactivated, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeleteUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.DeleteUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.DeleteUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func unmarshalJSON(data string) *structpb.Struct {
	user := new(structpb.Struct)
	err := user.UnmarshalJSON([]byte(data))
	if err != nil {
		logging.OnError(err).Fatal("unmarshalling user json")
	}
	return user
}

func TestServer_LockUser(t *testing.T) {
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
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.LockUserRequest) error
		req     *user.LockUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user lock, no userID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "user lock, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "user lock, not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user lock, no context",
			ctx:  context.Background(),
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user lock, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user lock, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user lock, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.LockUserRequest{},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user lock, already locked",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user lock, deactivated",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.LockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.LockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.LockUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_UnlockUser(t *testing.T) {
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
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.UnlockUserRequest) error
		req     *user.UnlockUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user unlock, no userID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "user unlock, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "user unlock, not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user unlock, no context",
			ctx:  context.Background(),
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user unlock, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user unlock, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.UnlockUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user unlock,no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.UnlockUserRequest{},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user unlock, already unlocked",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.UnlockUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				instance.UnlockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.UnlockUserRequest{
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
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.UnlockUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_DeactivateUser(t *testing.T) {
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
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.DeactivateUserRequest) error
		req     *user.DeactivateUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user deactivate, no userID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "user deactivate, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "user deactivate, not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user deactivate, no context",
			ctx:  context.Background(),
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user deactivate, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user deactivate, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user deactivate, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.DeactivateUserRequest{},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user deactivate, already deactivated",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user deactivate, locked",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.DeactivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.LockSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.DeactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.DeactivateUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_ActivateUser(t *testing.T) {
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
	orgResp := instance.CreateOrganization(isolatedIAMOwnerCTX, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.ActivateUserRequest) error
		req     *user.ActivateUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user activate, no userID",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "user activate, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "user activate, not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user activate, no context",
			ctx:  context.Background(),
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user activate, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user activate, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.ActivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user activate, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.ActivateUserRequest{},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_ORG,
					Id:   orgResp.GetOrganizationId(),
				},
			},
		},
		{
			name: "user activate, already activated",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ActivateUserRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.DeactivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				instance.ActivateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				return nil
			},
			req: &user.ActivateUserRequest{
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
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.ActivateUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}
