//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
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
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
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
	permissionSchemaResp := Instance.CreateUserSchema(IAMOwnerCTX, permissionSchema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

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
			ctx:  IAMOwnerCTX,
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
			ctx:  UserCTX,
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
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
		}, {
			name: "user create, full contact, ok",
			ctx:  IAMOwnerCTX,
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
							Address:      gofakeit.Email(),
							Verification: &user.SetEmail_ReturnCode{ReturnCode: &user.ReturnEmailVerificationCode{}},
						},
						Phone: &user.SetPhone{
							Number:       gofakeit.Phone(),
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
			got, err := Instance.Client.UserV3Alpha.CreateUser(tt.ctx, tt.req)
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
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
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
	permissionSchemaResp := Instance.CreateUserSchema(IAMOwnerCTX, permissionSchema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	type res struct {
		want            *resource_object.Details
		returnCodeEmail bool
		returnCodePhone bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.PatchUserRequest) error
		req     *user.PatchUserRequest
		res     res
		wantErr bool
	}{
		{
			name: "user create, no context",
			ctx:  context.Background(),
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user create, no permission",
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user create, invalid schema permission, owner",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user create, no change",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				req.User.Data = unmarshalJSON(data)
				req.User.SchemaId = gu.Ptr(schemaID)
				return nil
			},
			req: &user.PatchUserRequest{
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
			name: "user update, schema, ok",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				changedSchemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
				req.User.SchemaId = gu.Ptr(changedSchemaResp.Details.Id)
				return nil
			},
			req: &user.PatchUserRequest{
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
			name: "user update, schema and data, ok",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				changedSchemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
				req.User.SchemaId = gu.Ptr(changedSchemaResp.Details.Id)
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
			name: "user update, ok",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user update, full contact, ok",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.PatchUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.PatchUserRequest{
				User: &user.PatchUser{
					Data: unmarshalJSON("{\"name\": \"changed\"}"),
					Contact: &user.SetContact{
						Email: &user.SetEmail{
							Address:      gofakeit.Email(),
							Verification: &user.SetEmail_ReturnCode{ReturnCode: &user.ReturnEmailVerificationCode{}},
						},
						Phone: &user.SetPhone{
							Number:       gofakeit.Phone(),
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
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.PatchUser(tt.ctx, tt.req)
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

func TestServer_DeleteUser(t *testing.T) {
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.DeleteUserRequest) error
		req     *user.DeleteUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user delete, no userID",
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
			name: "user delete, no context",
			ctx:  context.Background(),
			dep: func(ctx context.Context, req *user.DeleteUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.DeleteUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.DeleteUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.DeleteUser(tt.ctx, tt.req)
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
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.LockUserRequest) error
		req     *user.LockUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user lock, no userID",
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user lock, already locked",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.LockUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.LockUser(tt.ctx, tt.req)
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
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.UnlockUserRequest) error
		req     *user.UnlockUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user unlock, no userID",
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.UnlockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
			dep: func(ctx context.Context, req *user.UnlockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.UnlockUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.UnlockUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
			name: "user unlock, already unlocked",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.UnlockUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
				Instance.UnlockSchemaUser(ctx, "", req.Id)
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
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.UnlockUser(tt.ctx, tt.req)
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
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.DeactivateUserRequest) error
		req     *user.DeactivateUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user deactivate, no userID",
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
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
			name: "user deactivate, already deactivated",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
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
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.DeactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.LockSchemaUser(ctx, "", req.Id)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.DeactivateUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}

func TestServer_ReactivateUser(t *testing.T) {
	ensureFeatureEnabled(t, IAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Instance.CreateUserSchema(IAMOwnerCTX, schema)
	orgResp := Instance.CreateOrganization(IAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, req *user.ReactivateUserRequest) error
		req     *user.ReactivateUserRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "user reactivate, no userID",
			ctx:  IAMOwnerCTX,
			req: &user.ReactivateUserRequest{
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
			name: "user reactivate, not existing",
			ctx:  IAMOwnerCTX,
			req: &user.ReactivateUserRequest{
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
			name: "user reactivate, not existing in org",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.ReactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
				return nil
			},
			req: &user.ReactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user reactivate, no context",
			ctx:  context.Background(),
			dep: func(ctx context.Context, req *user.ReactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
				return nil
			},
			req: &user.ReactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user reactivate, no permission",
			ctx:  UserCTX,
			dep: func(ctx context.Context, req *user.ReactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(IAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
				return nil
			},
			req: &user.ReactivateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "user reactivate, ok",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.ReactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
				return nil
			},
			req: &user.ReactivateUserRequest{
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
			name: "user reactivate, already reactivated",
			ctx:  IAMOwnerCTX,
			dep: func(ctx context.Context, req *user.ReactivateUserRequest) error {
				userResp := Instance.CreateSchemaUser(ctx, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				Instance.DeactivateSchemaUser(ctx, "", req.Id)
				Instance.ReactivateSchemaUser(ctx, "", req.Id)
				return nil
			},
			req: &user.ReactivateUserRequest{
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
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}
			got, err := Instance.Client.UserV3Alpha.ReactivateUser(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want, got.Details)
		})
	}
}
