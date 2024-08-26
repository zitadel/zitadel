//go:build integration

package userschema_test

import (
	"context"
	"fmt"
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
	_, _, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)

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
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
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
					TotalResult: 0,
				},
				Result: []*schema.UserSchema{},
			},
		},
		{
			name: "single (id), ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.SearchUserSchemasRequest{},
				prepare: func(request *schema.SearchUserSchemasRequest, resp *schema.SearchUserSchemasResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano() + 1)
					createResp := Tester.CreateUserSchemaWithType(isolatedIAMOwnerCTX, t, schemaType)
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

					resp.Result[0].Type = schemaType
					resp.Result[0].Details = &object.Details{
						Id:      createResp.GetDetails().GetId(),
						Created: createResp.GetDetails().GetCreated(),
						Changed: createResp.GetDetails().GetChanged(),
						Owner:   createResp.GetDetails().GetOwner(),
					}
					return nil
				},
			},
			want: &schema.SearchUserSchemasResponse{
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
				ctx: isolatedIAMOwnerCTX,
				req: &schema.SearchUserSchemasRequest{},
				prepare: func(request *schema.SearchUserSchemasRequest, resp *schema.SearchUserSchemasResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano())
					schemaType1 := schemaType + "_1"
					schemaType2 := schemaType + "_2"
					createResp := Tester.CreateUserSchemaWithType(isolatedIAMOwnerCTX, t, schemaType1)
					createResp2 := Tester.CreateUserSchemaWithType(isolatedIAMOwnerCTX, t, schemaType2)

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

					resp.Result[0].Type = schemaType1
					resp.Result[0].Details = &object.Details{
						Id:      createResp.GetDetails().GetId(),
						Created: createResp.GetDetails().GetCreated(),
						Changed: createResp.GetDetails().GetChanged(),
						Owner:   createResp.GetDetails().GetOwner(),
					}
					resp.Result[1].Type = schemaType2
					resp.Result[1].Details = &object.Details{
						Id:      createResp2.GetDetails().GetId(),
						Created: createResp2.GetDetails().GetCreated(),
						Changed: createResp2.GetDetails().GetChanged(),
						Owner:   createResp2.GetDetails().GetOwner(),
					}
					return nil
				},
			},
			want: &schema.SearchUserSchemasResponse{
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
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration := 20 * time.Second
			if ctxDeadline, ok := isolatedIAMOwnerCTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.SearchUserSchemas(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)

				fmt.Println(tt.want)
				fmt.Println(got)

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

func TestServer_GetUserSchemaByID(t *testing.T) {
	_, _, _, isolatedIAMOwnerCTX := Tester.UseIsolatedInstance(t, IAMOwnerCTX, SystemCTX)
	ensureFeatureEnabled(t, isolatedIAMOwnerCTX)

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
					createResp := Tester.CreateUserSchemaWithType(isolatedIAMOwnerCTX, t, schemaType)
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
				req: &schema.GetUserSchemaByIDRequest{
					Id: "notexisting",
				},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &schema.GetUserSchemaByIDRequest{},
				prepare: func(request *schema.GetUserSchemaByIDRequest, resp *schema.GetUserSchemaByIDResponse) error {
					schemaType := fmt.Sprint(time.Now().UnixNano() + 1)
					createResp := Tester.CreateUserSchemaWithType(isolatedIAMOwnerCTX, t, schemaType)
					request.Id = createResp.GetDetails().GetId()

					resp.Schema.Type = schemaType
					resp.Schema.Details = createResp.GetDetails()
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
			if tt.args.prepare != nil {
				err := tt.args.prepare(tt.args.req, tt.want)
				require.NoError(t, err)
			}

			retryDuration := 5 * time.Second
			if ctxDeadline, ok := isolatedIAMOwnerCTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}

			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetUserSchemaByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)

				fmt.Println(tt.want)
				fmt.Println(got)

				integration.AssertResourceDetails(t, tt.want.GetSchema().GetDetails(), got.GetSchema().GetDetails())
				grpc.AllFieldsEqual(t, tt.want.ProtoReflect(), got.ProtoReflect(), grpc.CustomMappers)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected user schema result")
		})
	}
}
