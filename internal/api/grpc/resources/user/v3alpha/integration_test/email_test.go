//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func TestServer_SetContactEmail(t *testing.T) {
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
		dep     func(req *user.SetContactEmailRequest) error
		req     *user.SetContactEmailRequest
		res     res
		wantErr bool
	}{
		{
			name: "email patch, no context",
			ctx:  context.Background(),
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address: integration.Email(),
				},
			},
			wantErr: true,
		},
		{
			name: "email patch, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address: integration.Email(),
				},
			},
			wantErr: true,
		},
		{
			name: "email patch, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
				Email: &user.SetEmail{
					Address: integration.Email(),
				},
			},
			wantErr: true,
		},
		{
			name: "email patch, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
				Email: &user.SetEmail{
					Address: integration.Email(),
				},
			},
			wantErr: true,
		},
		{
			name: "email patch, empty",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				email := integration.Email()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, email)
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{},
			},
			wantErr: true,
		},
		{
			name: "email patch, no change",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				email := integration.Email()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, email)
				req.Email.Address = email
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{},
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
			name: "email patch, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Email: &user.SetEmail{
					Address: integration.Email(),
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
			name: "email patch, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address:      integration.Email(),
					Verification: &user.SetEmail_ReturnCode{},
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
			name: "email patch, return, invalid template",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address:      integration.Email(),
					Verification: &user.SetEmail_SendCode{SendCode: &user.SendEmailVerificationCode{UrlTemplate: gu.Ptr("{{")}},
				},
			},
			wantErr: true,
		},
		{
			name: "email patch, verified, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address:      integration.Email(),
					Verification: &user.SetEmail_IsVerified{IsVerified: true},
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
			name: "email patch, template, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address:      integration.Email(),
					Verification: &user.SetEmail_SendCode{SendCode: &user.SendEmailVerificationCode{UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}")}},
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
			name: "email patch, sent, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.SetContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Email: &user.SetEmail{
					Address:      integration.Email(),
					Verification: &user.SetEmail_SendCode{},
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
			got, err := instance.Client.UserV3Alpha.SetContactEmail(tt.ctx, tt.req)
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

func TestServer_VerifyContactEmail(t *testing.T) {
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
		want *resource_object.Details
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.VerifyContactEmailRequest) error
		req     *user.VerifyContactEmailRequest
		res     res
		wantErr bool
	}{
		{
			name: "email verify, no context",
			ctx:  context.Background(),
			dep: func(req *user.VerifyContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email verify, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.VerifyContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email verify, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactEmailRequest) error {
				return nil
			},
			req: &user.VerifyContactEmailRequest{
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
			name: "email verify, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactEmailRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email verify, wrong code",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactEmailRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.VerifyContactEmailRequest{
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
			name: "email verify, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactEmailRequest{},
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
			name: "email verify, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.VerifyContactEmailRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				verifyResp := instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				req.VerificationCode = verifyResp.GetVerificationCode()
				return nil
			},
			req: &user.VerifyContactEmailRequest{
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
				require.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.VerifyContactEmail(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}

func TestServer_ResendContactEmailCode(t *testing.T) {
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
		dep     func(req *user.ResendContactEmailCodeRequest) error
		req     *user.ResendContactEmailCodeRequest
		res     res
		wantErr bool
	}{
		{
			name: "email resend, no context",
			ctx:  context.Background(),
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email resend, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email resend, not found",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
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
			name: "email resend, not found, org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "not existing",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email resend, no code",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				data := "{\"name\": \"user\"}"
				schemaID := schemaResp.GetDetails().GetId()
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaID, []byte(data))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "email resend, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{},
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
			name: "email resend, return, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Verification: &user.ResendContactEmailCodeRequest_ReturnCode{},
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
			name: "email resend, sent, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.ResendContactEmailCodeRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.UpdateSchemaUserEmail(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, integration.Email())
				return nil
			},
			req: &user.ResendContactEmailCodeRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Verification: &user.ResendContactEmailCodeRequest_SendCode{},
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
			got, err := instance.Client.UserV3Alpha.ResendContactEmailCode(tt.ctx, tt.req)
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
