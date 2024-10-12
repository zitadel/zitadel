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

var publicKeyExample = []byte("-----BEGIN PUBLIC KEY-----\nMIICITANBgkqhkiG9w0BAQEFAAOCAg4AMIICCQKCAgB5tWxwCGRloCqvpgI2ZXPl\nxQ+WZbQPuHTqAxwbXbsKOJoAAq16iHmzriLKpqVDxRUXqTH3cY0P0A1IZbCBB2gG\nyq3Lk08sR5ute+MEQ+QibX2qpk+mccRr+eP6B1otcyBWxRhZ/YtWphDpZ4GCb4oN\nAzTIebU0ztlu1OOnDDSEEhwScu2LhG40bx4hVU8XNgIqEjxiR61J89vfZpCmn0Rl\nsqYvmX9sqtqPokdsKl3LPItRyDAJMG0uhwwGKsHffDNeLDZN1OCZE/ZS7USarJQH\nbtGeqFQKsCL33xsKbNL+QjnAhqHW09bMdwofJvlwYLfL0rGJQr5aVCaERAfKAOE6\npy0nVkEJsRLxvdx/ZbTtZdCBk/LiznkE1xp9J02obQ+kWHtdUYxM1OSJqPRGQpbS\nZTxurdBQ43gRjO07iWNV9CB0i6QN2GtDBmHVb48i6aPdA++uJqnPYzy46FWA3KMA\nSlxiZ1RDcGH+fN9uklC2cwAurctAxed3Me2RYGdxl813udeV4Ef3qaiV2dix/pKA\nvN1KIfPTpTdULCDBLjtaAYflJ2WYXHeWMJMMC4oJc3bcKpA4mWjZibZ3pSGX/STQ\nXwHUtKsGlrVBSeqjjILVpH+2G0rusrqkGOlPKN+qOIsnwJf9x47v+xEw1slqdDWm\n+x3gc+8m9oowCcq20OeNTQIDAQAB\n-----END PUBLIC KEY-----")

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
					Type: &user.SetPublicKey_GenerateKey{GenerateKey: &user.GenerateKey{}},
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
					Type: &user.SetPublicKey_GenerateKey{GenerateKey: &user.GenerateKey{}},
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
			name: "publickey add, publickey invalid format",
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
					Type: &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: []byte("invalid")}},
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
					Type: &user.SetPublicKey_GenerateKey{GenerateKey: &user.GenerateKey{}},
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
					Type: &user.SetPublicKey_GenerateKey{GenerateKey: &user.GenerateKey{}},
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
					Type: &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: publicKeyExample}},
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
					Type: &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: publicKeyExample}},
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
			name: "publickey add, expirationdate invalid",
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
					ExpirationDate: timestamppb.New(time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC)),
					Type:           &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: publicKeyExample}},
				},
			},
			wantErr: true,
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
					ExpirationDate: timestamppb.New(time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)),
					Type:           &user.SetPublicKey_PublicKey{PublicKey: &user.ProvidedPublicKey{PublicKey: publicKeyExample}},
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
					Type: &user.SetPublicKey_GenerateKey{GenerateKey: &user.GenerateKey{}},
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
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), publicKeyExample)
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
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), publicKeyExample)
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
				publickeyResp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), userResp.GetDetails().GetId(), publicKeyExample)
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
				resp := instance.AddAuthenticatorPublicKey(isolatedIAMOwnerCTX, orgResp.GetOrganizationId(), req.Id, publicKeyExample)
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
