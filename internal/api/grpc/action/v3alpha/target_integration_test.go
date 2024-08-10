//go:build integration

package action_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestServer_CreateTarget(t *testing.T) {
	ensureFeatureEnabled(t)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.CreateTargetRequest
		want    *action.CreateTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:       fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty webhook url",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty request response url",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.SetRESTCall{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
				},
				Timeout: nil,
			},
			wantErr: true,
		},
		{
			name: "async, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestAsync{
					RestAsync: &action.SetRESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "webhook, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "webhook, interrupt on error, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "call, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},

		{
			name: "call, interruptOnError, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name:     fmt.Sprint(time.Now().UnixNano() + 1),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.SetRESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateTarget(tt.ctx, tt.req)
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

func TestServer_UpdateTarget(t *testing.T) {
	ensureFeatureEnabled(t)
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
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *action.UpdateTargetRequest) error {
				request.TargetId = "notexisting"
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
		},
		{
			name: "change name, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			want: &action.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestCall{
						RestCall: &action.SetRESTCall{
							InterruptOnError: true,
						},
					},
				},
			},
			want: &action.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					Endpoint: gu.Ptr("https://example.com/hooks/new"),
				},
			},
			want: &action.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					Timeout: durationpb.New(20 * time.Second),
				},
			},
			want: &action.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change type async, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeAsync, false).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestAsync{
						RestAsync: &action.SetRESTAsync{},
					},
				},
			},
			want: &action.UpdateTargetResponse{
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

			got, err := Client.UpdateTarget(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteTarget(t *testing.T) {
	ensureFeatureEnabled(t)
	target := Tester.CreateTarget(CTX, t, "", "https://example.com", domain.TargetTypeWebhook, false)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.DeleteTargetRequest
		want    *action.DeleteTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteTargetRequest{
				TargetId: target.GetId(),
			},
			wantErr: true,
		},
		{
			name: "empty id",
			ctx:  CTX,
			req: &action.DeleteTargetRequest{
				TargetId: "",
			},
			wantErr: true,
		},
		{
			name: "delete target",
			ctx:  CTX,
			req: &action.DeleteTargetRequest{
				TargetId: target.GetId(),
			},
			want: &action.DeleteTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.DeleteTarget(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
