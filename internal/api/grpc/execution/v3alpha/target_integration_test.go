//go:build integration

package execution_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client execution.ExecutionServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.ExecutionV3

		CTX, _ = Tester.WithAuthorization(ctx, integration.IAMOwner), errCtx
		return m.Run()
	}())
}

func TestServer_CreateTarget(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.CreateTargetRequest
		want    *execution.CreateTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				Type: 0,
			},
			wantErr: true,
		},
		{
			name: "empty url",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name: fmt.Sprint(time.Now().UnixNano() + 1),
				Type: execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:  "",
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name:          fmt.Sprint(time.Now().UnixNano() + 1),
				Type:          execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:           "https://example.com",
				Timeout:       nil,
				ExecutionType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty execution type, ok",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name:          fmt.Sprint(time.Now().UnixNano() + 1),
				Type:          execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:           "https://example.com",
				Timeout:       durationpb.New(10 * time.Second),
				ExecutionType: nil,
			},
			want: &execution.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "async execution, ok",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name:    fmt.Sprint(time.Now().UnixNano() + 1),
				Type:    execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:     "https://example.com",
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.CreateTargetRequest_IsAsync{
					IsAsync: true,
				},
			},
			want: &execution.CreateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "interrupt on error execution, ok",
			ctx:  CTX,
			req: &execution.CreateTargetRequest{
				Name:    fmt.Sprint(time.Now().UnixNano() + 1),
				Type:    execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:     "https://example.com",
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.CreateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
			},
			want: &execution.CreateTargetResponse{
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
	type args struct {
		ctx context.Context
		req *execution.UpdateTargetRequest
	}
	tests := []struct {
		name    string
		prepare func(request *execution.UpdateTargetRequest) error
		args    args
		want    *execution.UpdateTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: Tester.WithAuthorization(context.Background(), integration.OrgOwner),
				req: &execution.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *execution.UpdateTargetRequest) error {
				request.TargetId = "notexisting"
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			wantErr: true,
		},
		{
			name: "change name, ok",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					Name: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			want: &execution.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					Type: gu.Ptr(execution.TargetType_TARGET_TYPE_REST_REQUEST_RESPONSE),
				},
			},
			want: &execution.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					Url: gu.Ptr("https://example.com/hooks/new"),
				},
			},
			want: &execution.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					Timeout: durationpb.New(20 * time.Second),
				},
			},
			want: &execution.UpdateTargetResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "change execution type, ok",
			prepare: func(request *execution.UpdateTargetRequest) error {
				targetID := Tester.CreateTarget(CTX, t).GetId()
				request.TargetId = targetID
				return nil
			},
			args: args{
				ctx: CTX,
				req: &execution.UpdateTargetRequest{
					ExecutionType: &execution.UpdateTargetRequest_IsAsync{
						IsAsync: true,
					},
				},
			},
			want: &execution.UpdateTargetResponse{
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
	target := Tester.CreateTarget(CTX, t)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.DeleteTargetRequest
		want    *execution.DeleteTargetResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.DeleteTargetRequest{
				TargetId: target.GetId(),
			},
			wantErr: true,
		},
		{
			name: "empty id",
			ctx:  CTX,
			req: &execution.DeleteTargetRequest{
				TargetId: "",
			},
			wantErr: true,
		},
		{
			name: "delete target",
			ctx:  CTX,
			req: &execution.DeleteTargetRequest{
				TargetId: target.GetId(),
			},
			want: &execution.DeleteTargetResponse{
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
