//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_SetContactPhone(t *testing.T) {
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

	type res struct {
		want       *resource_object.Details
		returnCode bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.SetContactPhoneRequest) error
		req     *user.SetContactPhoneRequest
		res     res
		wantErr bool
	}{
		{
			name: "phone patch, no context",
			ctx:  context.Background(),
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{
					Number: integration.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "phone patch, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{
					Number: integration.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "phone patch, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
				Phone: &user.SetPhone{
					Number: integration.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "phone patch, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
				Phone: &user.SetPhone{
					Number: integration.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "phone patch, no change",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				number := integration.Phone()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, number)
				req.Phone.Number = number
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{},
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
			name: "phone patch, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Phone: &user.SetPhone{
					Number: integration.Phone(),
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
			name: "phone patch, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{
					Number:       integration.Phone(),
					Verification: &user.SetPhone_ReturnCode{},
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
				returnCode: true,
			},
		},
		{
			name: "phone patch, verified, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{
					Number:       integration.Phone(),
					Verification: &user.SetPhone_IsVerified{IsVerified: true},
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
			name: "phone patch, sent, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Phone: &user.SetPhone{
					Number:       integration.Phone(),
					Verification: &user.SetPhone_SendCode{},
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
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.SetContactPhone(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
			if tt.res.returnCode {
				assert.NotNil(t, got.VerificationCode)
			} else {
				assert.Nil(t, got.VerificationCode)
			}
		})
	}
}

func TestServer_VerifyContactPhone(t *testing.T) {
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

	type res struct {
		want            *resource_object.Details
		returnCodePhone bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.VerifyContactPhoneRequest) error
		req     *user.VerifyContactPhoneRequest
		res     res
		wantErr bool
	}{
		{
			name: "phone verify, no context",
			ctx:  context.Background(),
			dep: func(req *user.VerifyContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone verify, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.VerifyContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone verify, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactPhoneRequest) error {
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id:               "notexisting",
				VerificationCode: "unimportant",
			},
			wantErr: true,
		},
		{
			name: "phone verify, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone verify, wrong code",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactPhoneRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				VerificationCode: "wrong",
			},
			wantErr: true,
		},
		{
			name: "phone verify, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactPhoneRequest{},
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
			name: "phone verify, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactPhoneRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactPhoneRequest{
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
				returnCodePhone: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.VerifyContactPhone(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}

func TestServer_ResendContactPhoneCode(t *testing.T) {
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

	type res struct {
		want       *resource_object.Details
		returnCode bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.ResendContactPhoneCodeRequest) error
		req     *user.ResendContactPhoneCodeRequest
		res     res
		wantErr bool
	}{
		{
			name: "phone resend, no context",
			ctx:  context.Background(),
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone resend, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone resend, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
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
			name: "phone resend, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone resend, no code",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "phone resend, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{},
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
			name: "phone resend, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Verification: &user.ResendContactPhoneCodeRequest_ReturnCode{},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
				returnCode: true,
			},
		},
		{
			name: "phone resend, sent, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactPhoneCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserPhone(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Phone())
				return nil
			},
			req: &user.ResendContactPhoneCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Verification: &user.ResendContactPhoneCodeRequest_SendCode{},
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
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.ResendContactPhoneCode(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
			if tt.res.returnCode {
				assert.NotNil(t, got.VerificationCode)
			} else {
				assert.Nil(t, got.VerificationCode)
			}
		})
	}
}
