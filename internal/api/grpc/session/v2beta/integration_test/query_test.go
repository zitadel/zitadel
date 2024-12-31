//go:build integration

package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
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

type sessionAttr struct {
	ID           string
	UserID       string
	UserAgent    string
	CreationDate *timestamp.Timestamp
	ChangeDate   *timestamppb.Timestamp
	Details      *object.Details
}

type sessionAttrs []*sessionAttr

func (u sessionAttrs) ids() []string {
	ids := make([]string, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ids
}

func createSessions(ctx context.Context, t *testing.T, count int, userID string, userAgent string, lifetime *durationpb.Duration, metadata map[string][]byte) sessionAttrs {
	infos := make([]*sessionAttr, count)
	for i := 0; i < count; i++ {
		infos[i] = createSession(ctx, t, userID, userAgent, lifetime, metadata)
	}
	return infos
}

func createSession(ctx context.Context, t *testing.T, userID string, userAgent string, lifetime *durationpb.Duration, metadata map[string][]byte) *sessionAttr {
	req := &session.CreateSessionRequest{}
	if userID != "" {
		req.Checks = &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: userID,
				},
			},
		}
	}
	if userAgent != "" {
		req.UserAgent = &session.UserAgent{
			FingerprintId: gu.Ptr(userAgent),
			Ip:            gu.Ptr("1.2.3.4"),
			Description:   gu.Ptr("Description"),
			Header: map[string]*session.UserAgent_HeaderValues{
				"foo": {Values: []string{"foo", "bar"}},
			},
		}
	}
	if lifetime != nil {
		req.Lifetime = lifetime
	}
	if metadata != nil {
		req.Metadata = metadata
	}
	resp, err := Client.CreateSession(ctx, req)
	require.NoError(t, err)
	return &sessionAttr{
		resp.GetSessionId(),
		userID,
		userAgent,
		resp.GetDetails().GetChangeDate(),
		resp.GetDetails().GetChangeDate(),
		resp.GetDetails(),
	}
}

func TestServer_ListSessions(t *testing.T) {
	type args struct {
		ctx context.Context
		req *session.ListSessionsRequest
		dep func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr
	}
	tests := []struct {
		name                 string
		args                 args
		want                 *session.ListSessionsResponse
		wantFactors          []wantFactor
		wantExpirationWindow time.Duration
		wantErr              bool
	}{
		{
			name: "list sessions, not found",
			args: args{
				CTX,
				&session.ListSessionsRequest{
					Queries: []*session.SearchQuery{
						{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{"unknown"}}}},
					},
				},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					return []*sessionAttr{}
				},
			},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, wrong creator",
			args: args{
				UserCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, "", "", nil, nil)
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
					return []*sessionAttr{}
				},
			},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, full, ok",
			args: args{
				CTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
					return []*sessionAttr{info}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
				},
			},
		},
		{
			name: "list sessions, multiple, ok",
			args: args{
				CTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					infos := createSessions(ctx, t, 3, User.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: infos.ids()}}})
					return infos
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
				},
			},
		},
		{
			name: "list sessions, userid, ok",
			args: args{
				CTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					createdUser := createFullUser(ctx)
					info := createSession(ctx, t, createdUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_UserIdQuery{UserIdQuery: &session.UserIDQuery{Id: createdUser.GetUserId()}}})
					return []*sessionAttr{info}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header: map[string]*session.UserAgent_HeaderValues{
								"foo": {Values: []string{"foo", "bar"}},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos := tt.args.dep(CTX, t, tt.args.req)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListSessions(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				if !assert.NoError(ttt, err) {
					return
				}

				if !assert.Equal(ttt, got.Details.TotalResult, tt.want.Details.TotalResult) || !assert.Len(ttt, got.Sessions, len(tt.want.Sessions)) {
					return
				}

				for i := range infos {
					tt.want.Sessions[i].Id = infos[i].ID
					tt.want.Sessions[i].Sequence = infos[i].Details.GetSequence()
					tt.want.Sessions[i].CreationDate = infos[i].Details.GetChangeDate()
					tt.want.Sessions[i].ChangeDate = infos[i].Details.GetChangeDate()

					verifySession(ttt, got.Sessions[i], tt.want.Sessions[i], time.Minute, tt.wantExpirationWindow, infos[i].UserID, tt.wantFactors...)
				}
				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick)
		})
	}
}
