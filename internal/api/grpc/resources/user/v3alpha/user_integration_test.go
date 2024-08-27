//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
	//_, instanceID, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	isolatedIAMOwnerCTX := IAMOwnerCTX
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)
	schema := []byte(`{
		"$schema": "urn:zitadel:schema:v1",
			"type": "object",
			"properties": {
			"name": {
				"type": "string"
			}
		}
	}`)
	schemaResp := Tester.CreateUserSchema(isolatedIAMOwnerCTX, schema)
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
	permissionSchemaResp := Tester.CreateUserSchema(isolatedIAMOwnerCTX, permissionSchema)
	orgResp := Tester.CreateOrganization(isolatedIAMOwnerCTX, gofakeit.Name(), gofakeit.Email())

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
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				User: &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
			},
			wantErr: true,
		},
		{
			name: "user create, no context",
			ctx:  context.Background(),
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
				User:     &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
			},
			wantErr: true,
		},
		{
			name: "user create, no permission",
			ctx:  UserCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
				User:     &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
			},
			wantErr: true,
		},
		{
			name: "user create, invalid schema permission, owner",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: permissionSchemaResp.GetDetails().GetId(),
				User:     &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
			},
			wantErr: true,
		},
		{
			name: "user create, no user",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
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
			name: "user create, no data",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
				User:     &user.User{},
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
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
				User:     &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
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
			ctx:  isolatedIAMOwnerCTX,
			req: &user.CreateUserRequest{
				Organization: &object.Organization{
					Property: &object.Organization_Id{
						Id: orgResp.GetOrganizationId(),
					},
				},
				SchemaId: schemaResp.GetDetails().GetId(),
				User:     &user.User{Data: unmarshalJSON("{\"name\": \"user\"}")},
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
			got, err := Tester.Client.UserV3Alpha.CreateUser(tt.ctx, tt.req)
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

func unmarshalJSON(data string) *structpb.Struct {
	user := new(structpb.Struct)
	err := user.UnmarshalJSON([]byte(data))
	if err != nil {
		logging.OnError(err).Fatal("unmarshalling user json")
	}
	return user
}
