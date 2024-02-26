//go:build integration

package schema_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	schema "github.com/zitadel/zitadel/pkg/grpc/user/schema/v3alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client schema.UserSchemaServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		Client = Tester.Client.UserSchemaV3
		return m.Run()
	}())
}

func TestServer_CreateUserSchema(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		req     *schema.CreateUserSchemaRequest
		want    *schema.CreateUserSchemaResponse
		wantErr bool
	}{
		{
			name: "missing permission, error",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: "",
			},
			wantErr: true,
		},
		{
			name: "empty schema, error",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
			},
			wantErr: true,
		},
		{
			name: "invalid schema, error",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			wantErr: true,
		},
		{
			name: "no authenticators, ok",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			want: &schema.CreateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "invalid authenticator, error",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			wantErr: true,
		},
		{
			name: "with authenticator, ok",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			want: &schema.CreateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "with invalid permission, error",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			wantErr: true,
		},
		{
			name: "with valid permission, ok",
			ctx:  CTX,
			req: &schema.CreateUserSchemaRequest{
				Type: fmt.Sprint(time.Now().UnixNano() + 1),
				DataType: &schema.CreateUserSchemaRequest_Schema{
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
			want: &schema.CreateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateUserSchema(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
			assert.NotEmpty(t, got.GetId())
		})
	}
}

func TestServer_UpdateUserSchema(t *testing.T) {
	type args struct {
		ctx context.Context
		req *schema.UpdateUserSchemaRequest
	}
	tests := []struct {
		name    string
		prepare func(request *schema.UpdateUserSchemaRequest) error
		args    args
		want    *schema.UpdateUserSchemaResponse
		wantErr bool
	}{
		{
			name: "missing permission, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &schema.UpdateUserSchemaRequest{
					Type: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
		},
		{
			name: "missing id, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{},
			},
			wantErr: true,
		},
		{
			name: "not existing, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				request.Id = "notexisting"
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{},
			},
			wantErr: true,
		},
		{
			name: "empty type, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					Type: gu.Ptr(""),
				},
			},
			wantErr: true,
		},
		{
			name: "update type, ok",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					Type: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			want: &schema.UpdateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "empty schema, ok",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					DataType: &schema.UpdateUserSchemaRequest_Schema{},
				},
			},
			want: &schema.UpdateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "invalid schema, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					DataType: &schema.UpdateUserSchemaRequest_Schema{
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
			name: "update schema, ok",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					DataType: &schema.UpdateUserSchemaRequest_Schema{
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
			want: &schema.UpdateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "invalid authenticator, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_UNSPECIFIED,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "update authenticator, ok",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					PossibleAuthenticators: []schema.AuthenticatorType{
						schema.AuthenticatorType_AUTHENTICATOR_TYPE_USERNAME,
					},
				},
			},
			want: &schema.UpdateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.UpdateUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeactivateUserSchema(t *testing.T) {
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
				CTX,
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
				CTX,
				&schema.DeactivateUserSchemaRequest{},
				func(request *schema.DeactivateUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					return nil
				},
			},
			want: &schema.DeactivateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "inactive, error",
			args: args{
				CTX,
				&schema.DeactivateUserSchemaRequest{},
				func(request *schema.DeactivateUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					_, err := Client.DeactivateUserSchema(CTX, &schema.DeactivateUserSchemaRequest{
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

			got, err := Client.DeactivateUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_ReactivateUserSchema(t *testing.T) {
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
				CTX,
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
				ctx: CTX,
				req: &schema.ReactivateUserSchemaRequest{},
				prepare: func(request *schema.ReactivateUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "inactive, ok",
			args: args{
				ctx: CTX,
				req: &schema.ReactivateUserSchemaRequest{},
				prepare: func(request *schema.ReactivateUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					_, err := Client.DeactivateUserSchema(CTX, &schema.DeactivateUserSchemaRequest{
						Id: schemaID,
					})
					return err
				},
			},
			want: &schema.ReactivateUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.ReactivateUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteUserSchema(t *testing.T) {
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
				CTX,
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
				ctx: CTX,
				req: &schema.DeleteUserSchemaRequest{},
				prepare: func(request *schema.DeleteUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					return nil
				},
			},
			want: &schema.DeleteUserSchemaResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "deleted, error",
			args: args{
				ctx: CTX,
				req: &schema.DeleteUserSchemaRequest{},
				prepare: func(request *schema.DeleteUserSchemaRequest) error {
					schemaID := Tester.CreateUserSchema(CTX, t).GetId()
					request.Id = schemaID
					_, err := Client.DeleteUserSchema(CTX, &schema.DeleteUserSchemaRequest{
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

			got, err := Client.DeleteUserSchema(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
