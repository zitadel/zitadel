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

	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
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
				TargetType: &action.CreateTargetRequest_RestRequestResponse{
					RestRequestResponse: &action.SetRESTRequestResponse{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com",
					},
				},
				Timeout:       nil,
				ExecutionType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty execution type, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com",
					},
				},
				Timeout:       durationpb.New(10 * time.Second),
				ExecutionType: nil,
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "async execution, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.CreateTargetRequest_IsAsync{
					IsAsync: true,
				},
			},
			want: &action.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "interrupt on error execution, ok",
			ctx:  CTX,
			req: &action.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.CreateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
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
				targetID := Tester.CreateTarget(CTX, t).GetId()
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
				targetID := Tester.CreateTarget(CTX, t).GetId()
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
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestRequestResponse{
						RestRequestResponse: &action.SetRESTRequestResponse{
							Url: "https://example.com",
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
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestWebhook{
						RestWebhook: &action.SetRESTWebhook{
							Url: "https://example.com/hooks/new",
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
			name: "change timeout, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
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
			name: "change execution type, ok",
			prepare: func(request *action.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &action.UpdateTargetRequest{
					ExecutionType: &action.UpdateTargetRequest_IsAsync{
						IsAsync: true,
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
	target := Tester.CreateTarget(CTX, t)
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
