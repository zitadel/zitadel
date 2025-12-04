//go:build integration

package userschema_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

func TestServer_CreateUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *schema.CreateUserSchemaRequest
		want    *schema.CreateUserSchemaResponse
		wantErr bool
	}{
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.Username(),
				},
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: "",
				},
			},
			wantErr: true,
		},
		{
			name: "empty schema, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid schema, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string",
										"required": true
									},
									"description": {
										"type": "string"
									}
								}
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no authenticators, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									},
									"description": {
										"type": "string"
									}
								},
								"required": ["name"]
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
				},
			},
			want: &schema.CreateUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "invalid authenticator, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									},
									"description": {
										"type": "string"
									}
								},
								"required": ["name"]
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "with authenticator, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									},
									"description": {
										"type": "string"
									}
								},
								"required": ["name"]
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME,
					},
				},
			},
			want: &schema.CreateUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "with invalid permission, error",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									},
									"description": {
										"type": "string",
										"urn:zitadel:schema:permission": "read"
									}
								},
								"required": ["name"]
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "with valid permission, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &schema.CreateUserSchemaRequest{
				UserSchema: &schema.UserSchema{
					Type: integration.UserSchemaName(),
					DataType: &schema.UserSchema_Schema{
						Schema: func() *structpb.Struct {
							s := new(structpb.Struct)
							err := s.UnmarshalJSON([]byte(`
							{
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									},
									"description": {
										"type": "string",
										"urn:zitadel:schema:permission": {
											"owner": "rw",
											"self": "r"
										}
									}
								},
								"required": ["name"]
							}
						`))
							require.NoError(t, err)
							return s
						}(),
					},
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME,
					},
				},
			},
			want: &schema.CreateUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.UserSchemaV3.CreateUserSchema(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertResourceDetails(t, tt.want.GetDetails(), got.GetDetails())
		})
	}
}

func TestServer_UpdateUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *schema.PatchUserSchemaRequest
	}
	tests := []struct {
		name    string
		prepare func(request *schema.PatchUserSchemaRequest) error
		args    args
		want    *schema.PatchUserSchemaResponse
		wantErr bool
	}{
		{
			name: "missing permission, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						Type: gu.Ptr(integration.UserSchemaName()),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing id, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{},
			},
			wantErr: true,
		},
		{
			name: "not existing, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				request.Id = "notexisting"
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{},
			},
			wantErr: true,
		},
		{
			name: "empty type, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						Type: gu.Ptr(""),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "update type, ok",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						Type: gu.Ptr(integration.UserSchemaName()),
					},
				},
			},
			want: &schema.PatchUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "empty schema, ok",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						DataType: &schema.PatchUserSchema_Schema{},
					},
				},
			},
			want: &schema.PatchUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "invalid schema, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						DataType: &schema.PatchUserSchema_Schema{
							Schema: func() *structpb.Struct {
								s := new(structpb.Struct)
								err := s.UnmarshalJSON([]byte(`
							{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string",
										"required": true
									},
									"description": {
										"type": "string"
									}
								}
							}
						`))
								require.NoError(t, err)
								return s
							}(),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "update schema, ok",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						DataType: &schema.PatchUserSchema_Schema{
							Schema: func() *structpb.Struct {
								s := new(structpb.Struct)
								err := s.UnmarshalJSON([]byte(`
								{
									"$schema": "urn:zitadel:schema:v1",
									"type": "object",
									"properties": {
										"name": {
											"type": "string"
										},
										"description": {
											"type": "string"
										}
									},
									"required": ["name"]
								}
							`))
								require.NoError(t, err)
								return s
							}(),
						},
					},
				},
			},
			want: &schema.PatchUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "invalid authenticator, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						PossibleAuthenticators: []schema.AuthenticatorType{
							schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "update authenticator, ok",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						PossibleAuthenticators: []schema.AuthenticatorType{
							schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME,
						},
					},
				},
			},
			want: &schema.PatchUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "inactive, error",
			prepare: func(request *schema.PatchUserSchemaRequest) error {
				schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
				_, err := instance.Client.UserSchemaV3.DeactivateUserSchema(isolatedIAMOwnerCTX, &schema.DeactivateUserSchemaRequest{
					Id: schemaID,
				})
				require.NoError(t, err)
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.PatchUserSchemaRequest{
					UserSchema: &schema.PatchUserSchema{
						Type: gu.Ptr(integration.UserSchemaName()),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := instance.Client.UserSchemaV3.PatchUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want.GetDetails(), got.GetDetails())
		})
	}
}

func TestServer_DeactivateUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx     context.Context
		req     *schema.DeactivateUserSchemaRequest
		prepare func(request *schema.DeactivateUserSchemaRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.DeactivateUserSchemaResponse
		wantErr bool
	}{
		{
			name: "not existing, error",
			args: args{
				isolatedIAMOwnerCTX,
				&schema.DeactivateUserSchemaRequest{
					Id: "notexisting",
				},
				func(request *schema.DeactivateUserSchemaRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "active, ok",
			args: args{
				isolatedIAMOwnerCTX,
				&schema.DeactivateUserSchemaRequest{},
				func(request *schema.DeactivateUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					return nil
				},
			},
			want: &schema.DeactivateUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "inactive, error",
			args: args{
				isolatedIAMOwnerCTX,
				&schema.DeactivateUserSchemaRequest{},
				func(request *schema.DeactivateUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					_, err := instance.Client.UserSchemaV3.DeactivateUserSchema(isolatedIAMOwnerCTX, &schema.DeactivateUserSchemaRequest{
						Id: schemaID,
					})
					return err
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := instance.Client.UserSchemaV3.DeactivateUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want.GetDetails(), got.GetDetails())
		})
	}
}

func TestServer_ReactivateUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx     context.Context
		req     *schema.ReactivateUserSchemaRequest
		prepare func(request *schema.ReactivateUserSchemaRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.ReactivateUserSchemaResponse
		wantErr bool
	}{
		{
			name: "not existing, error",
			args: args{
				isolatedIAMOwnerCTX,
				&schema.ReactivateUserSchemaRequest{
					Id: "notexisting",
				},
				func(request *schema.ReactivateUserSchemaRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "active, error",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.ReactivateUserSchemaRequest{},
				prepare: func(request *schema.ReactivateUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "inactive, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.ReactivateUserSchemaRequest{},
				prepare: func(request *schema.ReactivateUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					_, err := instance.Client.UserSchemaV3.DeactivateUserSchema(isolatedIAMOwnerCTX, &schema.DeactivateUserSchemaRequest{
						Id: schemaID,
					})
					return err
				},
			},
			want: &schema.ReactivateUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := instance.Client.UserSchemaV3.ReactivateUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want.GetDetails(), got.GetDetails())
		})
	}
}

func TestServer_DeleteUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx     context.Context
		req     *schema.DeleteUserSchemaRequest
		prepare func(request *schema.DeleteUserSchemaRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.DeleteUserSchemaResponse
		wantErr bool
	}{
		{
			name: "not existing, error",
			args: args{
				isolatedIAMOwnerCTX,
				&schema.DeleteUserSchemaRequest{
					Id: "notexisting",
				},
				func(request *schema.DeleteUserSchemaRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "delete, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.DeleteUserSchemaRequest{},
				prepare: func(request *schema.DeleteUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					return nil
				},
			},
			want: &schema.DeleteUserSchemaResponse{
				Details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "deleted, error",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.DeleteUserSchemaRequest{},
				prepare: func(request *schema.DeleteUserSchemaRequest) error {
					schemaID := instance.CreateUserSchemaEmpty(isolatedIAMOwnerCTX).GetDetails().GetId()
					request.Id = schemaID
					_, err := instance.Client.UserSchemaV3.DeleteUserSchema(isolatedIAMOwnerCTX, &schema.DeleteUserSchemaRequest{
						Id: schemaID,
					})
					return err
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := instance.Client.UserSchemaV3.DeleteUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertResourceDetails(t, tt.want.GetDetails(), got.GetDetails())
		})
	}
}
