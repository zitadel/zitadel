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
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
)

func TestServer_CreateTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.CreateTargetRequest
		want    *action.CreateTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &action.CreateTargetRequest{
				Name: gofakeit.Name(),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:       gofakeit.Name(),
				TargetType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty webhook url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: gofakeit.Name(),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty request response url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: gofakeit.Name(),
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{},
				},
				Timeout: nil,
			},
			wantErr: true,
		},
		{
			name: "async, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestAsync{
					RestAsync: &action.RESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Id:           "notempty",
				CreationDate: timestamppb.Now(),
				SigningKey:   "notempty",
			},
		},
		{
			name: "webhook, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Id:           "notempty",
				CreationDate: timestamppb.Now(),
				SigningKey:   "notempty",
			},
		},
		{
			name: "webhook, interrupt on error, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Id:           "notempty",
				CreationDate: timestamppb.Now(),
				SigningKey:   "notempty",
			},
		},
		{
			name: "call, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Id:           "notempty",
				CreationDate: timestamppb.Now(),
				SigningKey:   "notempty",
			},
		},

		{
			name: "call, interruptOnError, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     gofakeit.Name(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Id:           "notempty",
				CreationDate: timestamppb.Now(),
				SigningKey:   "notempty",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.ActionV2beta.CreateTarget(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateTargetResponse(t, tt.want, got)
		})
	}
}

func assertCreateTargetResponse(t *testing.T, expectedResp *action.CreateTargetResponse, actualResp *action.CreateTargetResponse) {
	if expectedResp.GetCreationDate() == nil {
		wantCreationDate := expectedResp.GetCreationDate().AsTime()
		assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), wantCreationDate.Add(-time.Minute), wantCreationDate.Add(time.Minute))
	}
	if expectedResp.GetId() != "" {
		assert.NotEmpty(t, actualResp.GetId())
	}
	if expectedResp.GetSigningKey() != "" {
		assert.NotEmpty(t, actualResp.GetSigningKey())
	}
}

func TestServer_UpdateTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	isolatedIAMOwnerCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *action.UpdateTargetRequest
	}
	tests := []struct {
		name    string
		prepare func(request *action.UpdateTargetRequest) error
		args    args
		want    *action.UpdateTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(gofakeit.Name()),
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *action.UpdateTargetRequest) error {
				request.Id = "notexisting"
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(gofakeit.Name()),
				},
			},
			wantErr: true,
		},
		{
			name: "change name, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(gofakeit.Name()),
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: nil,
			},
		},
		{
			name: "regenerate signingkey, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					ExpirationSigningKey: durationpb.New(0 * time.Second),
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: gu.Ptr("notempty"),
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestCall{
						RestCall: &action.RESTCall{
							InterruptOnError: true,
						},
					},
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: nil,
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Endpoint: gu.Ptr("https://example.com/hooks/new"),
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: nil,
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Timeout: durationpb.New(20 * time.Second),
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: nil,
			},
		},
		{
			name: "change type async, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", domain.TargetTypeAsync, false).GetId()
				request.Id = targetID
				return nil
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestAsync{
						RestAsync: &action.RESTAsync{},
					},
				},
			},
			want: &action.UpdateTargetResponse{
				ChangeDate: timestamppb.Now(),
				SigningKey: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)
			// We want to have the same response no matter how often we call the function
			instance.Client.ActionV2beta.UpdateTarget(tt.args.ctx, tt.args.req)
			got, err := instance.Client.ActionV2beta.UpdateTarget(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertUpdateTargetResponse(t, tt.want, got)
		})
	}
}

func assertUpdateTargetResponse(t *testing.T, expectedResp *action.UpdateTargetResponse, actualResp *action.UpdateTargetResponse) {
	if expectedResp.GetChangeDate() == nil {
		wantCreationDate := expectedResp.GetChangeDate().AsTime()
		assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), wantCreationDate.Add(-time.Minute), wantCreationDate.Add(time.Minute))
	}
	if expectedResp.GetSigningKey() != "" {
		assert.NotEmpty(t, actualResp.GetSigningKey())
	}
}

func TestServer_DeleteTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ensureFeatureEnabled(t, instance)
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *action.DeleteTargetRequest)
		req     *action.DeleteTargetRequest
		want    *action.DeleteTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &action.DeleteTargetRequest{
				Id: "notexisting",
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
			name: "delete target, not existing",
			ctx:  iamOwnerCtx,
			req: &action.DeleteTargetRequest{
				Id: "notexisting",
			},
			want: &action.DeleteTargetResponse{
				DeletionDate: timestamppb.Now(),
			},
		},
		{
			name: "delete target",
			ctx:  iamOwnerCtx,
			prepare: func(request *action.DeleteTargetRequest) {
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			req: &action.DeleteTargetRequest{},
			want: &action.DeleteTargetResponse{
				DeletionDate: timestamppb.Now(),
			},
		},
		{
			name: "delete target, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *action.DeleteTargetRequest) {
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				instance.DeleteTarget(iamOwnerCtx, t, targetID)
			},
			req: &action.DeleteTargetRequest{},
			want: &action.DeleteTargetResponse{
				DeletionDate: timestamppb.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			got, err := instance.Client.ActionV2beta.DeleteTarget(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteTargetResponse(t, tt.want, got)
		})
	}
}

func assertDeleteTargetResponse(t *testing.T, expectedResp *action.DeleteTargetResponse, actualResp *action.DeleteTargetResponse) {
	if expectedResp.GetDeletionDate() == nil {
		wantCreationDate := expectedResp.GetDeletionDate().AsTime()
		assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), wantCreationDate.Add(-time.Minute), wantCreationDate.Add(time.Minute))
	}
}
