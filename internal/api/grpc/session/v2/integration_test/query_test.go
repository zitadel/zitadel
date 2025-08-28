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
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestServer_GetSession(t *testing.T) {
	t.Parallel()
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
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
				IAMOwnerCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) uint64 {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
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
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
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
			t.Parallel()
			var sequence uint64
			if tt.args.dep != nil {
				sequence = tt.args.dep(LoginCTX, t, tt.args.req)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetSession(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

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
	t.Parallel()
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
			name: "list sessions, no permission",
			args: args{
				UserCTX,
				&session.ListSessionsRequest{
					Queries: []*session.SearchQuery{},
				},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, "", "", nil, nil)
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
					return []*sessionAttr{}
				},
			},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, permission, ok",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, "", "", nil, nil)
					request.Queries = append(request.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
					return []*sessionAttr{info}
				},
			},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{{}},
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
		{
			name: "list sessions, own creator, ok",
			args: args{
				LoginCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{}}})
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
			name: "list sessions, creator, ok",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{Id: gu.Ptr(Instance.Users.Get(integration.UserTypeLogin).ID)}}})
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
			name: "list sessions, wrong creator",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{}}})
					return []*sessionAttr{}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, empty creator",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{Id: gu.Ptr("")}}})
					return []*sessionAttr{}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			wantErr:              true,
		},
		{
			name: "list sessions, useragent, ok",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("useragent")}}})
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
							FingerprintId: gu.Ptr("useragent"),
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
			name: "list sessions, wrong useragent",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("wronguseragent")}}})
					return []*sessionAttr{}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, empty useragent",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("")}}})
					return []*sessionAttr{}
				},
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			wantErr:              true,
		},
		{
			name: "list sessions, expiration date query, ok",
			args: args{
				IAMOwnerCTX,
				&session.ListSessionsRequest{},
				func(ctx context.Context, t *testing.T, request *session.ListSessionsRequest) []*sessionAttr {
					info := createSession(ctx, t, User.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					request.Queries = append(request.Queries,
						&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
						&session.SearchQuery{Query: &session.SearchQuery_ExpirationDateQuery{
							ExpirationDateQuery: &session.ExpirationDateQuery{ExpirationDate: timestamppb.Now(),
								Method: objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS,
							}}})
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
							FingerprintId: gu.Ptr("useragent"),
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
			t.Parallel()
			infos := tt.args.dep(LoginCTX, t, tt.args.req)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListSessions(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				// expected count of sessions is not equal to created dependencies
				require.Len(ttt, tt.want.Sessions, len(infos))


				// expected count of sessions is not equal to received sessions
				require.Equal(ttt, tt.want.Details.TotalResult, got.Details.TotalResult)
				require.Len(ttt, got.Sessions, len(tt.want.Sessions))

				for i := range infos {
					tt.want.Sessions[i].Id = infos[i].ID
					tt.want.Sessions[i].Sequence = infos[i].Details.GetSequence()
					tt.want.Sessions[i].CreationDate = infos[i].Details.GetChangeDate()
					tt.want.Sessions[i].ChangeDate = infos[i].Details.GetChangeDate()

					// only check for contents of the session, not sorting for now
					found := false
					for _, session := range got.Sessions {
						if session.Id == infos[i].ID {
							verifySession(ttt, session, tt.want.Sessions[i], time.Minute, tt.wantExpirationWindow, infos[i].UserID, tt.wantFactors...)
							found = true
						}
					}
					assert.True(t, found)
				}

				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick)
		})
	}
}

func TestServer_ListSessions_with_expiration_date_filter(t *testing.T) {
	t.Parallel()
	// session with no expiration
	session1, err := Client.CreateSession(IAMOwnerCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)

	// session with expiration
	session2, err := Client.CreateSession(IAMOwnerCTX, &session.CreateSessionRequest{
		Lifetime: durationpb.New(1 * time.Second),
	})
	require.NoError(t, err)

	// wait until the second session expires
	time.Sleep(2 * time.Second)

	// with comparison method GREATER_OR_EQUALS, only the active session should be returned
	listSessionsResponse1, err := Client.ListSessions(IAMOwnerCTX,
		&session.ListSessionsRequest{
			Queries: []*session.SearchQuery{
				{
					Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{session1.SessionId}}},
				},
				{
					Query: &session.SearchQuery_ExpirationDateQuery{
						ExpirationDateQuery: &session.ExpirationDateQuery{
							ExpirationDate: timestamppb.Now(),
							Method:         objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS,
						},
					},
				},
			},
		})
	require.NoError(t, err)
	require.Len(t, listSessionsResponse1.Sessions, 1)
	assert.Equal(t, session1.SessionId, listSessionsResponse1.Sessions[0].Id)

	// with comparison method LESS_OR_EQUALS, only the expired session should be returned
	listSessionsResponse2, err := Client.ListSessions(IAMOwnerCTX,
		&session.ListSessionsRequest{
			Queries: []*session.SearchQuery{
				{
					Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{session2.SessionId}}},
				},
				{
					Query: &session.SearchQuery_ExpirationDateQuery{
						ExpirationDateQuery: &session.ExpirationDateQuery{
							ExpirationDate: timestamppb.Now(),
							Method:         objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS,
						},
					},
				},
			},
		})
	require.NoError(t, err)
	require.Len(t, listSessionsResponse2.Sessions, 1)
	assert.Equal(t, session2.SessionId, listSessionsResponse2.Sessions[0].Id)
}
