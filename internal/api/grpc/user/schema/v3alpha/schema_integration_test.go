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

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
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

func ensureFeatureEnabled(t *testing.T) {
	f, err := Tester.Client.FeatureV2.GetInstanceFeatures(CTX, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.UserSchema.GetEnabled() {
		return
	}
	_, err = Tester.Client.FeatureV2.SetInstanceFeatures(CTX, &feature.SetInstanceFeaturesRequest{
		UserSchema: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := CTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := Tester.Client.FeatureV2.GetInstanceFeatures(CTX, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(ttt, err)
			if f.UserSchema.GetEnabled() {
				return
			}
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature")
}

func TestServer_CreateUserSchema(t *testing.T) {
	ensureFeatureEnabled(t)

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
	ensureFeatureEnabled(t)

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
		{
			name: "inactive, error",
			prepare: func(request *schema.UpdateUserSchemaRequest) error {
				schemaID := Tester.CreateUserSchema(CTX, t).GetId()
				_, err := Client.DeactivateUserSchema(CTX, &schema.DeactivateUserSchemaRequest{
					Id: schemaID,
				})
				require.NoError(t, err)
				request.Id = schemaID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &schema.UpdateUserSchemaRequest{
					Type: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
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
	ensureFeatureEnabled(t)

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
	ensureFeatureEnabled(t)

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
	ensureFeatureEnabled(t)

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

func TestServer_GetUserSchemaByID(t *testing.T) {
	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *schema.GetUserSchemaByIDRequest
		prepare func(request *schema.GetUserSchemaByIDRequest, resp *schema.GetUserSchemaByIDResponse) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.GetUserSchemaByIDResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &schema.GetUserSchemaByIDRequest{},
				prepare: func(request *schema.GetUserSchemaByIDRequest, resp *schema.GetUserSchemaByIDResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano() + 1)
					createResp := Tester.CreateUserSchemaWithType(CTX, t, schemaType)
					request.Id = createResp.GetId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "not existing, error",
			args: args{
				ctx: CTX,
				req: &schema.GetUserSchemaByIDRequest{
					Id: "notexisting",
				},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: CTX,
				req: &schema.GetUserSchemaByIDRequest{},
				prepare: func(request *schema.GetUserSchemaByIDRequest, resp *schema.GetUserSchemaByIDResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano() + 1)
					createResp := Tester.CreateUserSchemaWithType(CTX, t, schemaType)
					request.Id = createResp.GetId()

					resp.Schema.Id = createResp.GetId()
					resp.Schema.Type = schemaType
					resp.Schema.Details = &object.Details{
						Sequence:      createResp.GetDetails().GetSequence(),
						ChangeDate:    createResp.GetDetails().GetChangeDate(),
						ResourceOwner: createResp.GetDetails().GetResourceOwner(),
					}
					return nil
				},
			},
			want: &schema.GetUserSchemaByIDResponse{
				Schema: &schema.UserSchema{
					State:                  schema.State_STATE_ACTIVE,
					Revision:               1,
					Schema:                 userSchema,
					PossibleAuthenticators: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureFeatureEnabled(t)
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration := 5 * time.Second
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetUserSchemaByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)

				integration.AssertDetails(t, tt.want.GetSchema(), got.GetSchema())
				grpc.AllFieldsEqual(t, tt.want.ProtoReflect(), got.ProtoReflect(), grpc.CustomMappers)

			}, retryDuration, time.Millisecond*100, "timeout waiting for expected user schema result")
		})
	}
}

func TestServer_ListUserSchemas(t *testing.T) {
	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *schema.ListUserSchemasRequest
		prepare func(request *schema.ListUserSchemasRequest, resp *schema.ListUserSchemasResponse) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.ListUserSchemasResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &schema.ListUserSchemasRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found, error",
			args: args{
				ctx: CTX,
				req: &schema.ListUserSchemasRequest{
					Queries: []*schema.SearchQuery{
						{
							Query: &schema.SearchQuery_IdQuery{
								IdQuery: &schema.IDQuery{
									Id: "notexisting",
								},
							},
						},
					},
				},
			},
			want: &schema.ListUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
				Result: []*schema.UserSchema{},
			},
		},
		{
			name: "single (id), ok",
			args: args{
				ctx: CTX,
				req: &schema.ListUserSchemasRequest{},
				prepare: func(request *schema.ListUserSchemasRequest, resp *schema.ListUserSchemasResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano() + 1)
					createResp := Tester.CreateUserSchemaWithType(CTX, t, schemaType)
					request.Queries = []*schema.SearchQuery{
						{
							Query: &schema.SearchQuery_IdQuery{
								IdQuery: &schema.IDQuery{
									Id:     createResp.GetId(),
									Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
								},
							},
						},
					}

					resp.Result[0].Id = createResp.GetId()
					resp.Result[0].Type = schemaType
					resp.Result[0].Details = &object.Details{
						Sequence:      createResp.GetDetails().GetSequence(),
						ChangeDate:    createResp.GetDetails().GetChangeDate(),
						ResourceOwner: createResp.GetDetails().GetResourceOwner(),
					}
					return nil
				},
			},
			want: &schema.ListUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*schema.UserSchema{
					{
						State:                  schema.State_STATE_ACTIVE,
						Revision:               1,
						Schema:                 userSchema,
						PossibleAuthenticators: nil,
					},
				},
			},
		},
		{
			name: "multiple (type), ok",
			args: args{
				ctx: CTX,
				req: &schema.ListUserSchemasRequest{},
				prepare: func(request *schema.ListUserSchemasRequest, resp *schema.ListUserSchemasResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano())
					schemaType1 := schemaType + "_1"
					schemaType2 := schemaType + "_2"
					createResp := Tester.CreateUserSchemaWithType(CTX, t, schemaType1)
					createResp2 := Tester.CreateUserSchemaWithType(CTX, t, schemaType2)

					request.SortingColumn = schema.FieldName_FIELD_NAME_TYPE
					request.Query = &object.ListQuery{Asc: true}
					request.Queries = []*schema.SearchQuery{
						{
							Query: &schema.SearchQuery_TypeQuery{
								TypeQuery: &schema.TypeQuery{
									Type:   schemaType,
									Method: object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
								},
							},
						},
					}

					resp.Result[0].Id = createResp.GetId()
					resp.Result[0].Type = schemaType1
					resp.Result[0].Details = &object.Details{
						Sequence:      createResp.GetDetails().GetSequence(),
						ChangeDate:    createResp.GetDetails().GetChangeDate(),
						ResourceOwner: createResp.GetDetails().GetResourceOwner(),
					}
					resp.Result[1].Id = createResp2.GetId()
					resp.Result[1].Type = schemaType2
					resp.Result[1].Details = &object.Details{
						Sequence:      createResp2.GetDetails().GetSequence(),
						ChangeDate:    createResp2.GetDetails().GetChangeDate(),
						ResourceOwner: createResp2.GetDetails().GetResourceOwner(),
					}
					return nil
				},
			},
			want: &schema.ListUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult: 2,
				},
				Result: []*schema.UserSchema{
					{
						State:                  schema.State_STATE_ACTIVE,
						Revision:               1,
						Schema:                 userSchema,
						PossibleAuthenticators: nil,
					},
					{
						State:                  schema.State_STATE_ACTIVE,
						Revision:               1,
						Schema:                 userSchema,
						PossibleAuthenticators: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureFeatureEnabled(t)
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration := 20 * time.Second
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListUserSchemas(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)

				// always first check length, otherwise its failed anyway
				assert.Len(ttt, got.Result, len(tt.want.Result))
				for i := range tt.want.Result {
					//
					grpc.AllFieldsEqual(t, tt.want.Result[i].ProtoReflect(), got.Result[i].ProtoReflect(), grpc.CustomMappers)
				}
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected user schema result")
		})
	}
}
