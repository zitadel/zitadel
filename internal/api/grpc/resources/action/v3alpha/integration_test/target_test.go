//go:build integration

package action_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func TestServer_CreateTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.Target
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &action.Target{
				Name: gofakeit.Name(),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:       gofakeit.Name(),
				TargetType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty webhook url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: gofakeit.Name(),
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty request response url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name: gofakeit.Name(),
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
				Timeout: nil,
			},
			wantErr: true,
		},
		{
			name: "async, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestAsync{
					RestAsync: &action.SetRESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},
		{
			name: "webhook, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},
		{
			name: "webhook, interrupt on error, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},
		{
			name: "call, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},

		{
			name: "call, interruptOnError, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.Target{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.ActionV3Alpha.CreateTarget(tt.ctx, &action.CreateTargetRequest{Target: tt.req})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				integration.AssertResourceDetails(t, tt.want, got.Details)
				assert.NotEmpty(t, got.GetSigningKey())
			}
		})
	}
}

func TestServer_PatchTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *action.PatchTargetRequest
	}
	type want struct {
		details    *resource_object.Details
		signingKey bool
	}
	tests := []struct {
		name    string
		prepare func(request *action.PatchTargetRequest) error
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(gofakeit.Name()),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *action.PatchTargetRequest) error {
				request.Id = "notexisting"
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(gofakeit.Name()),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change name, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Name: gu.Ptr(gofakeit.Name()),
					},
				},
			},
			want: want{
				details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "regenerate signingkey, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						ExpirationSigningKey: durationpb.New(0 * time.Second),
					},
				},
			},
			want: want{
				details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
				signingKey: true,
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						TargetType: &action.PatchTarget_RestCall{
							RestCall: &action.SetRESTCall{
								InterruptOnError: true,
							},
						},
					},
				},
			},
			want: want{
				details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Endpoint: gu.Ptr("https://example.com/hooks/new"),
					},
				},
			},
			want: want{
				details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						Timeout: durationpb.New(20 * time.Second),
					},
				},
			},
			want: want{
				details: &resource_object.Details{
					Changed: timestamppb.Now(),
					Owner: &object.Owner{
						Type: object.OwnerType_OWNER_TYPE_INSTANCE,
						Id:   instance.ID(),
					},
				},
			},
		},
		{
			name: "change type async, ok",
			prepare: func(request *action.PatchTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeAsync, false).GetDetails().GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.PatchTargetRequest{
					Target: &action.PatchTarget{
						TargetType: &action.PatchTarget_RestAsync{
							RestAsync: &action.SetRESTAsync{},
						},
					},
				},
			},
			want: want{
				details: &resource_object.Details{
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
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)
			// We want to have the same response no matter how often we call the function
			instance.Client.ActionV3Alpha.PatchTarget(tt.args.ctx, tt.args.req)
			got, err := instance.Client.ActionV3Alpha.PatchTarget(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				integration.AssertResourceDetails(t, tt.want.details, got.Details)
				if tt.want.signingKey {
					assert.NotEmpty(t, got.SigningKey)
				}
			}
		})
	}
}

func TestServer_DeleteTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	target := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", domain.TargetTypeWebhook, false)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.DeleteTargetRequest
		want    *resource_object.Details
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &action.DeleteTargetRequest{
				Id: target.GetDetails().GetId(),
			},
			wantErr: true,
		},
		{
			name: "empty id",
			ctx:  iamOwnerCtx,
			req: &action.DeleteTargetRequest{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "delete target",
			ctx:  iamOwnerCtx,
			req: &action.DeleteTargetRequest{
				Id: target.GetDetails().GetId(),
			},
			want: &resource_object.Details{
				Changed: timestamppb.Now(),
				Owner: &object.Owner{
					Type: object.OwnerType_OWNER_TYPE_INSTANCE,
					Id:   instance.ID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.ActionV3Alpha.DeleteTarget(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				integration.AssertResourceDetails(t, tt.want, got.Details)
			}
		})
	}
}
