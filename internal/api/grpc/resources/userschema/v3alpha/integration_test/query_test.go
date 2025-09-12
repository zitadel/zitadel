//go:build integration

package userschema_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

func TestServer_ListUserSchemas(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *schema.SearchUserSchemasRequest
		prepare func(request *schema.SearchUserSchemasRequest, resp *schema.SearchUserSchemasResponse) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.SearchUserSchemasResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &schema.SearchUserSchemasRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found, error",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.SearchUserSchemasRequest{
					Filters: []*schema.SearchFilter{
						{
							Filter: &schema.SearchFilter_IdFilter{
								IdFilter: &schema.IDFilter{
									Id: "notexisting",
								},
							},
						},
					},
				},
			},
			want: &schema.SearchUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				Result: []*schema.GetUserSchema{},
			},
		},
		{
			name: "single (id), ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.SearchUserSchemasRequest{},
				prepare: func(request *schema.SearchUserSchemasRequest, resp *schema.SearchUserSchemasResponse) error {
					schemaType := integration.UserSchemaName()
					createResp := instance.CreateUserSchemaEmptyWithType(isolatedIAMOwnerCTX, schemaType)
					request.Filters = []*schema.SearchFilter{
						{
							Filter: &schema.SearchFilter_IdFilter{
								IdFilter: &schema.IDFilter{
									Id:     createResp.GetDetails().GetId(),
									Method: object.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					}
					resp.Result[0].Config.Type = schemaType
					resp.Result[0].Details = createResp.GetDetails()
					// as schema is freshly created, the changed date is the created date
					resp.Result[0].Details.Created = resp.Result[0].Details.GetChanged()
					resp.Details.Timestamp = resp.Result[0].Details.GetChanged()
					return nil
				},
			},
			want: &schema.SearchUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Result: []*schema.GetUserSchema{
					{
						State:    schema.State_STATE_ACTIVE,
						Revision: 1,
						Config: &schema.UserSchema{
							Type: "",
							DataType: &schema.UserSchema_Schema{
								Schema: userSchema,
							},
							PossibleAuthenticators: nil,
						},
					},
				},
			},
		},
		{
			name: "multiple (type), ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.SearchUserSchemasRequest{},
				prepare: func(request *schema.SearchUserSchemasRequest, resp *schema.SearchUserSchemasResponse) error {
					schemaType := integration.UserSchemaName()
					schemaType1 := schemaType + "_1"
					schemaType2 := schemaType + "_2"
					createResp := instance.CreateUserSchemaEmptyWithType(isolatedIAMOwnerCTX, schemaType1)
					createResp2 := instance.CreateUserSchemaEmptyWithType(isolatedIAMOwnerCTX, schemaType2)

					request.SortingColumn = gu.Ptr(schema.FieldName_FIELD_NAME_TYPE)
					request.Query = &object.SearchQuery{Desc: false}
					request.Filters = []*schema.SearchFilter{
						{
							Filter: &schema.SearchFilter_TypeFilter{
								TypeFilter: &schema.TypeFilter{
									Type:   schemaType,
									Method: object.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH,
								},
							},
						},
					}

					resp.Result[0].Config.Type = schemaType1
					resp.Result[0].Details = createResp.GetDetails()
					resp.Result[1].Config.Type = schemaType2
					resp.Result[1].Details = createResp2.GetDetails()
					return nil
				},
			},
			want: &schema.SearchUserSchemasResponse{
				Details: &object.ListDetails{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Result: []*schema.GetUserSchema{
					{
						State:    schema.State_STATE_ACTIVE,
						Revision: 1,
						Config: &schema.UserSchema{
							DataType: &schema.UserSchema_Schema{
								Schema: userSchema,
							},
						},
					},
					{
						State:    schema.State_STATE_ACTIVE,
						Revision: 1,
						Config: &schema.UserSchema{
							DataType: &schema.UserSchema_Schema{
								Schema: userSchema,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.UserSchemaV3.SearchUserSchemas(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					for i := range tt.want.Result {
						wantSchema := tt.want.Result[i]
						gotSchema := got.Result[i]

						integration.AssertResourceDetails(ttt, wantSchema.GetDetails(), gotSchema.GetDetails())
						wantSchema.Details = gotSchema.GetDetails()
						grpc.AllFieldsEqual(ttt, wantSchema.ProtoReflect(), gotSchema.ProtoReflect(), grpc.CustomMappers)
					}
				}
				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected user schema result")
		})
	}
}

func TestServer_GetUserSchema(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *schema.GetUserSchemaRequest
		prepare func(request *schema.GetUserSchemaRequest, resp *schema.GetUserSchemaResponse) error
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.GetUserSchemaResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &schema.GetUserSchemaRequest{},
				prepare: func(request *schema.GetUserSchemaRequest, resp *schema.GetUserSchemaResponse) error {
					schemaType := integration.UserSchemaName()
					createResp := instance.CreateUserSchemaEmptyWithType(isolatedIAMOwnerCTX, schemaType)
					request.Id = createResp.GetDetails().GetId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "not existing, error",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.GetUserSchemaRequest{
					Id: "notexisting",
				},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.GetUserSchemaRequest{},
				prepare: func(request *schema.GetUserSchemaRequest, resp *schema.GetUserSchemaResponse) error {
					schemaType := integration.UserSchemaName()
					createResp := instance.CreateUserSchemaEmptyWithType(isolatedIAMOwnerCTX, schemaType)
					request.Id = createResp.GetDetails().GetId()

					resp.UserSchema.Config.Type = schemaType
					resp.UserSchema.Details = createResp.GetDetails()
					return nil
				},
			},
			want: &schema.GetUserSchemaResponse{
				UserSchema: &schema.GetUserSchema{
					State:    schema.State_STATE_ACTIVE,
					Revision: 1,
					Config: &schema.UserSchema{
						DataType: &schema.UserSchema_Schema{
							Schema: userSchema,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(isolatedIAMOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.UserSchemaV3.GetUserSchema(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err, "Error: "+err.Error())
					return
				}
				assert.NoError(ttt, err)

				wantSchema := tt.want.GetUserSchema()
				gotSchema := got.GetUserSchema()
				integration.AssertResourceDetails(ttt, wantSchema.GetDetails(), gotSchema.GetDetails())
				wantSchema.Details = got.GetUserSchema().GetDetails()
				grpc.AllFieldsEqual(ttt, wantSchema.ProtoReflect(), gotSchema.ProtoReflect(), grpc.CustomMappers)
			}, retryDuration, tick, "timeout waiting for expected user schema result")
		})
	}
}
