//go:build integration

package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestServer_GetSession(t *testing.T) {
	type args struct {
		ctx context.Context
		req *session.GetSessionRequest
		dep func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64
	}
	tests := []struct {
		name                 string
		args                 args
		want                 *session.GetSessionResponse
		wantFactors          []wantFactor
		wantExpirationWindow time.Duration
		wantErr              bool
	}{
		{
			name: "get session, no id provided",
			args: args{
				CTX,
				&session.GetSessionRequest{
					SessionId: "",
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "get session, not found",
			args: args{
				CTX,
				&session.GetSessionRequest{
					SessionId: "unknown",
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "get session, no permission",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					return resp.GetDetails().GetSequence()
				},
			},
			wantErr: true,
		},
		{
			name: "get session, permission, ok",
			args: args{
				CTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					return resp.GetDetails().GetSequence()
				},
			},
			want: &session.GetSessionResponse{
				Session: &session.Session{},
			},
		},
		{
			name: "get session, token, ok",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails().GetSequence()
				},
			},
			want: &session.GetSessionResponse{
				Session: &session.Session{},
			},
		},
		{
			name: "get session, user agent, ok",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("fingerPrintID"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails().GetSequence()
				},
			},
			want: &session.GetSessionResponse{
				Session: &session.Session{
					UserAgent: &session.UserAgent{
						FingerprintId: gu.Ptr("fingerPrintID"),
						Ip:            gu.Ptr("1.2.3.4"),
						Description:   gu.Ptr("Description"),
						Header: map[string]*session.UserAgent_HeaderValues{
							"foo": {Values: []string{"foo", "bar"}},
						},
					},
				},
			},
		},
		{
			name: "get session, lifetime, ok",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{
						Lifetime: durationpb.New(5 * time.Minute),
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails().GetSequence()
				},
			},
			wantExpirationWindow: 5 * time.Minute,
			want: &session.GetSessionResponse{
				Session: &session.Session{},
			},
		},
		{
			name: "get session, metadata, ok",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{
						Metadata: map[string][]byte{"foo": []byte("bar")},
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails().GetSequence()
				},
			},
			want: &session.GetSessionResponse{
				Session: &session.Session{
					Metadata: map[string][]byte{"foo": []byte("bar")},
				},
			},
		},
		{
			name: "get session, user, ok",
			args: args{
				UserCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{
						Checks: &session.Checks{
							User: &session.CheckUser{
								Search: &session.CheckUser_UserId{
									UserId: User.GetUserId(),
								},
							},
						},
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails().GetSequence()
				},
			},
			wantFactors: []wantFactor{wantUserFactor},
			want: &session.GetSessionResponse{
				Session: &session.Session{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sequence uint64
			if tt.args.dep != nil {
				sequence = tt.args.dep(tt.args.ctx, t, tt.args.req)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetSession(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				if !assert.NoError(ttt, err) {
					return
				}

				tt.want.Session.Id = tt.args.req.SessionId
				tt.want.Session.Sequence = sequence
				verifySession(ttt, got.GetSession(), tt.want.GetSession(), time.Minute, tt.wantExpirationWindow, User.GetUserId(), tt.wantFactors...)
			}, retryDuration, tick)
		})
	}
}
