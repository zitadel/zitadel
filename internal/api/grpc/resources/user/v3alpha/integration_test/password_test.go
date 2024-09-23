//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_SetPassword(t *testing.T) {
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
		dep     func(req *user.SetPasswordRequest) error
		req     *user.SetPasswordRequest
		res     res
		wantErr bool
	}{
		{
			name: "password set, no context",
			ctx:  context.Background(),
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "password set, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "password set, password empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: "",
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "password set, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "username add, user not existing",
			ctx:  isolatedIAMOwnerCTX,

			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				Id: "not existing",
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "password set, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
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
			name: "password set, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetPasswordRequest{
				NewPassword: &user.SetPassword{
					Type: &user.SetPassword_Password{
						Password: gofakeit.Password(true, true, true, true, false, 12),
					},
					ChangeRequired: false,
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
			got, err := instance.Client.UserV3Alpha.SetPassword(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}
