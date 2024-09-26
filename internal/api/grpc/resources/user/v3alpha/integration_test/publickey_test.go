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

func TestServer_AddPublicKey(t *testing.T) {
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
		dep     func(req *user.AddPublicKeyRequest) error
		req     *user.AddPublicKeyRequest
		res     res
		wantErr bool
	}{
		{
			name: "publickey add, no context",
			ctx:  context.Background(),
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey add, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey add, publickey empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: []byte("")}},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey add, user not existing in org",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: "notexisting",
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey add, user not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				Id: "notexisting",
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey add, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
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
			name: "publickey add, no org, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
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
			name: "publickey add, expirationdate, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: []byte(gofakeit.BitcoinPrivateKey())}},
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
			name: "publickey add, generated, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.AddPublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.AddPublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKey: &user.SetPublicKey{
					Type: &user.SetPublicKey_GeneratedKey{GeneratedKey: &user.GeneratedKey{}},
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
			got, err := instance.Client.UserV3Alpha.AddPublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)
		})
	}
}

func TestServer_DeletePublicKey(t *testing.T) {
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
		dep     func(req *user.RemovePublicKeyRequest) error
		req     *user.RemovePublicKeyRequest
		res     res
		wantErr bool
	}{
		{
			name: "publickey delete, no context",
			ctx:  context.Background(),
			dep: func(req *user.RemovePublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), []byte(gofakeit.BitcoinPrivateKey()))
				req.PublicKeyId = publickeyResp.GetPublicKeyId()
				return nil
			},
			req: &user.RemovePublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey delete, no permission",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			dep: func(req *user.RemovePublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), []byte(gofakeit.BitcoinPrivateKey()))
				req.PublicKeyId = publickeyResp.GetPublicKeyId()
				return nil
			},
			req: &user.RemovePublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey remove, id empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemovePublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKeyId: "notempty",
			},
			wantErr: true,
		},
		{
			name: "publickey delete, userid empty",
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeLogin),
			req: &user.RemovePublicKeyRequest{
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
			name: "publickey remove, not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &user.RemovePublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
				PublicKeyId: "notempty",
				Id:          "notexisting",
			},
			wantErr: true,
		},
		{
			name: "publickey remove, no publickey",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				return nil
			},
			req: &user.RemovePublicKeyRequest{
				Organization: &object.Organization{
					Property: &object.Organization_OrgId{
						OrgId: orgResp.GetOrganizationId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "publickey remove, ok",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), []byte(gofakeit.BitcoinPrivateKey()))
				req.PublicKeyId = publickeyResp.GetPublicKeyId()
				return nil
			},
			req: &user.RemovePublicKeyRequest{
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
			name: "publickey remove, already removed",
			ctx:  isolatedIAMOwnerCTX,
			dep: func(req *user.RemovePublicKeyRequest) error {
				userResp := instance.CreateSchemaUser(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), schemaResp.GetDetails().GetId(), []byte("{\"name\": \"user\"}"))
				req.Id = userResp.GetDetails().GetId()
				resp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, []byte(gofakeit.BitcoinPrivateKey()))
				req.PublicKeyId = resp.GetPublicKeyId()
				instance.RemoveAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, req.PublicKeyId)
				return nil
			},
			req: &user.RemovePublicKeyRequest{
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
			got, err := instance.Client.UserV3Alpha.RemovePublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return

			}
			assert.NoError(t, err)
			integration.AssertResourceDetails(t, tt.res.want, got.Details)

		})
	}
}
