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
						Password: fakePassword(),
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
						Password: fakePassword(),
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
						Password: fakePassword(),
					},
					ChangeRequired: false,
				},
			},
			wantErr: true,
		},
		{
			name: "password set, user not existing",
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
						Password: fakePassword(),
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
						Password: fakePassword(),
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
						Password: fakePassword(),
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
			name: "password set, code, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				passwordResp := instance.RequestAuthenticatorPasswordReset(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
				req.NewPassword.Verification = &user.SetPassword_VerificationCode{VerificationCode: passwordResp.GetVerificationCode()}
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
						Password: fakePassword(),
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
			name: "password set, code, failed",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				instance.RequestAuthenticatorPasswordReset(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id)
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
						Password: fakePassword(),
					},
					ChangeRequired: false,
					Verification:   &user.SetPassword_VerificationCode{VerificationCode: "notreally"},
				},
			},
			wantErr: true,
		},
		{
			name: "password set, code, no set",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
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
						Password: fakePassword(),
					},
					ChangeRequired: false,
					Verification:   &user.SetPassword_VerificationCode{VerificationCode: "notreally"},
				},
			},
			wantErr: true,
		},
		{
			name: "password set, current password, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				password := fakePassword()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, password)
				req.NewPassword.Verification = &user.SetPassword_CurrentPassword{CurrentPassword: password}
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
						Password: fakePassword(),
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
			name: "password set, current password, failed",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.SetPasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
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
						Password: fakePassword(),
					},
					ChangeRequired: false,
					Verification:   &user.SetPassword_CurrentPassword{CurrentPassword: fakePassword()},
				},
			},
			wantErr: true,
		},
		{
			name: "password set, current password, not set",
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
						Password: fakePassword(),
					},
					ChangeRequired: false,
					Verification:   &user.SetPassword_CurrentPassword{CurrentPassword: fakePassword()},
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

func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 5) + "Password1!"
}

func TestServer_RequestPasswordReset(t *testing.T) {
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
		code bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(req *user.RequestPasswordResetRequest) error
		req     *user.RequestPasswordResetRequest
		res     res
		wantErr bool
	}{
		{
			name: "password reset, no context",
			ctx:  context.Background(),
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
			},
			wantErr: true,
		},
		{
			name: "password reset, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
			},
			wantErr: true,
		},
		{
			name: "password reset, no password set",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
			},
			wantErr: true,
		},
		{
			name: "password reset, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
			},
			wantErr: true,
		},
		{
			name: "password reset, user not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				Id:     "not existing",
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
			},
			wantErr: true,
		},
		{
			name: "password reset, email, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
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
			name: "password reset, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Medium: &user.RequestPasswordResetRequest_SendEmail{},
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
			name: "password reset, phone, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_SendSms{},
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
			name: "password reset, code, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RequestPasswordResetRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RequestPasswordResetRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Medium: &user.RequestPasswordResetRequest_ReturnCode{},
			},
			res: res{
				want: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_ORG,
						Id:   orgResp.GetOrganizationId(),
					},
				},
				code: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.req)
				assert.NoError(t, err)
			}
			got, err := instance.Client.UserV3Alpha.RequestPasswordReset(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
			if tt.res.code {
				require.NotEmpty(t, got.VerificationCode)
			}
		})
	}
}

func TestServer_RemovePassword(t *testing.T) {
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
		dep     func(req *user.RemovePasswordRequest) error
		req     *user.RemovePasswordRequest
		res     res
		wantErr bool
	}{
		{
			name: "password remove, no context",
			ctx:  context.Background(),
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RemovePasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "password remove, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RemovePasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "password remove, no password set",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.RemovePasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "password remove, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RemovePasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "password remove, user not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.RemovePasswordRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				Id: "not existing",
			},
			wantErr: true,
		},
		{
			name: "password remove, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RemovePasswordRequest{
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
			name: "password remove, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePasswordRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				instance.SetAuthenticatorPassword(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, fakePassword())
				return nil
			},
			req: &user.RemovePasswordRequest{},
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
			got, err := instance.Client.UserV3Alpha.RemovePassword(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}
