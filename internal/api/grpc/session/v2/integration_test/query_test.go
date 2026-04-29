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
	featurepb "github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	instancepb "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_GetSession(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *session.GetSessionRequest
		dep func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					return resp.GetDetails()
				},
			},
			wantErr: true,
		},
		{
			name: "get session, permission, ok",
			args: args{
				IAMOwnerCTX,
				&session.GetSessionRequest{},
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					return resp.GetDetails()
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{})
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails()
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
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
					return resp.GetDetails()
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
						Lifetime: durationpb.New(5 * time.Minute),
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails()
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
					resp, err := Client.CreateSession(ctx, &session.CreateSessionRequest{
						Metadata: map[string][]byte{"foo": []byte("bar")},
					},
					)
					require.NoError(t, err)
					request.SessionId = resp.SessionId
					request.SessionToken = gu.Ptr(resp.SessionToken)
					return resp.GetDetails()
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
				func(ctx context.Context, t *testing.T, request *session.GetSessionRequest) *object.Details {
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
					return resp.GetDetails()
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
			var details *object.Details
			if tt.args.dep != nil {
				details = tt.args.dep(LoginCTX, t, tt.args.req)
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
				tt.want.Session.Sequence = details.GetSequence()
				tt.want.Session.CreationDate = details.GetChangeDate()
				tt.want.Session.ChangeDate = details.GetChangeDate()
				verifySession(ttt, got.GetSession(), tt.want.GetSession(), tt.wantExpirationWindow, User.GetUserId(), tt.wantFactors...)
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

// listSessionsEnv holds the environment context for ListSessions tests,
// allowing the same test cases to run against both the eventstore and relational backends.
type listSessionsEnv struct {
	ownerCtx    context.Context // IAMOwner-level context; can list all sessions
	loginCtx    context.Context // login-user context; used for session creation in deps
	noPermCtx   context.Context // no-permission context; used to verify access restrictions
	client      session.SessionServiceClient
	loginUserID string
	testUser    *user.AddHumanUserResponse
	newUser     func(ctx context.Context) *user.AddHumanUserResponse
	// isRelational controls behavior that currently differs between the eventstore
	// and relational backends (e.g. sequence field, permission-check semantics).
	isRelational bool
}

// runListSessionsTestCases defines and runs the full ListSessions test suite against env.
// This is called by both [TestServer_ListSessions] (eventstore) and [TestServer_ListSessions_Relational].
func runListSessionsTestCases(t *testing.T, env listSessionsEnv) {
	t.Helper()

	// newSess creates a session on env.client via the eventstore model (grpc).
	newSess := func(t *testing.T, ctx context.Context, userID, userAgent string, lifetime *durationpb.Duration, metadata map[string][]byte) *sessionAttr {
		req := &session.CreateSessionRequest{}
		if userID != "" {
			req.Checks = &session.Checks{
				User: &session.CheckUser{Search: &session.CheckUser_UserId{UserId: userID}},
			}
		}
		if userAgent != "" {
			req.UserAgent = &session.UserAgent{
				FingerprintId: gu.Ptr(userAgent),
				Ip:            gu.Ptr("1.2.3.4"),
				Description:   gu.Ptr("Description"),
				Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
			}
		}
		if lifetime != nil {
			req.Lifetime = lifetime
		}
		if metadata != nil {
			req.Metadata = metadata
		}
		resp, err := env.client.CreateSession(ctx, req)
		require.NoError(t, err)
		return &sessionAttr{
			ID:           resp.GetSessionId(),
			UserID:       userID,
			UserAgent:    userAgent,
			CreationDate: resp.GetDetails().GetChangeDate(),
			ChangeDate:   resp.GetDetails().GetChangeDate(),
			Details:      resp.GetDetails(),
		}
	}

	tests := []struct {
		name                 string
		ctx                  context.Context
		req                  *session.ListSessionsRequest
		dep                  func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr
		want                 *session.ListSessionsResponse
		wantFactors          []wantFactor
		wantExpirationWindow time.Duration
		wantErr              bool
	}{
		{
			name: "list sessions, not found",
			ctx:  env.ownerCtx,
			req: &session.ListSessionsRequest{
				Queries: []*session.SearchQuery{
					{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{"unknown"}}}},
				},
			},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				return []*sessionAttr{}
			},
			want: &session.ListSessionsResponse{
				Details:  &object.ListDetails{TotalResult: 0, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{},
			},
		},
		{
			// In the eventstore path, the session exists (TotalResult=1) but the permission
			// callback hides its content (Sessions=[]). In the relational path the permission
			// check is not yet implemented (returns SQL true), so the session is fully visible.
			// TODO(IAM-Marco): align expectations once relational permission checks are implemented.
			name: "list sessions, no permission",
			ctx:  env.noPermCtx,
			req:  &session.ListSessionsRequest{Queries: []*session.SearchQuery{}},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, "", "", nil, nil)
				req.Queries = append(req.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
				if env.isRelational {
					return []*sessionAttr{info}
				}
				return []*sessionAttr{}
			},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: func() []*session.Session {
					if env.isRelational {
						return []*session.Session{{}}
					}
					return []*session.Session{}
				}(),
			},
		},
		{
			name: "list sessions, permission, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, "", "", nil, nil)
				req.Queries = append(req.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
				return []*sessionAttr{info}
			},
			want: &session.ListSessionsResponse{
				Details:  &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{{}},
			},
		},
		{
			name: "list sessions, full, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			name: "list sessions, multiple, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				infos := make([]*sessionAttr, 3)
				ids := make([]string, 3)
				for i := range infos {
					infos[i] = newSess(t, env.loginCtx, env.testUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
					ids[i] = infos[i].ID
				}
				req.Queries = append(req.Queries, &session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: ids}}})
				return infos
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 3, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			name: "list sessions, userid, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				createdUser := env.newUser(env.loginCtx)
				info := newSess(t, env.loginCtx, createdUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries, &session.SearchQuery{Query: &session.SearchQuery_UserIdQuery{UserIdQuery: &session.UserIDQuery{Id: createdUser.GetUserId()}}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			name: "list sessions, own creator, ok",
			ctx:  env.loginCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{}}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			name: "list sessions, creator, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{Id: gu.Ptr(env.loginUserID)}}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("agent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			// CreatorQuery with nil ID means "filter by caller's userID" (ownerCtx user).
			// Session was created by loginCtx user, so ownerUser != creator → no match.
			name: "list sessions, wrong creator",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "agent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{}}})
				return []*sessionAttr{}
			},
			want: &session.ListSessionsResponse{
				Details:  &object.ListDetails{TotalResult: 0, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, empty creator",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_CreatorQuery{CreatorQuery: &session.CreatorQuery{Id: gu.Ptr("")}}})
				return []*sessionAttr{}
			},
			wantErr: true,
		},
		{
			name: "list sessions, useragent, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("useragent")}}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("useragent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
		{
			name: "list sessions, wrong useragent",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("wronguseragent")}}})
				return []*sessionAttr{}
			},
			want: &session.ListSessionsResponse{
				Details:  &object.ListDetails{TotalResult: 0, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{},
			},
		},
		{
			name: "list sessions, empty useragent",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_UserAgentQuery{UserAgentQuery: &session.UserAgentQuery{FingerprintId: gu.Ptr("")}}})
				return []*sessionAttr{}
			},
			wantErr: true,
		},
		{
			name: "list sessions, expiration date query, ok",
			ctx:  env.ownerCtx,
			req:  &session.ListSessionsRequest{},
			dep: func(t *testing.T, req *session.ListSessionsRequest) []*sessionAttr {
				info := newSess(t, env.loginCtx, env.testUser.GetUserId(), "useragent", durationpb.New(time.Minute*5), map[string][]byte{"key": []byte("value")})
				req.Queries = append(req.Queries,
					&session.SearchQuery{Query: &session.SearchQuery_IdsQuery{IdsQuery: &session.IDsQuery{Ids: []string{info.ID}}}},
					&session.SearchQuery{Query: &session.SearchQuery_ExpirationDateQuery{
						ExpirationDateQuery: &session.ExpirationDateQuery{
							ExpirationDate: timestamppb.Now(),
							Method:         objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS,
						},
					}})
				return []*sessionAttr{info}
			},
			wantExpirationWindow: time.Minute * 5,
			wantFactors:          []wantFactor{wantUserFactor},
			want: &session.ListSessionsResponse{
				Details: &object.ListDetails{TotalResult: 1, Timestamp: timestamppb.Now()},
				Sessions: []*session.Session{
					{
						Metadata: map[string][]byte{"key": []byte("value")},
						UserAgent: &session.UserAgent{
							FingerprintId: gu.Ptr("useragent"),
							Ip:            gu.Ptr("1.2.3.4"),
							Description:   gu.Ptr("Description"),
							Header:        map[string]*session.UserAgent_HeaderValues{"foo": {Values: []string{"foo", "bar"}}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			infos := tt.dep(t, tt.req)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := env.client.ListSessions(tt.ctx, tt.req)
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
					if env.isRelational {
						tt.want.Sessions[i].Sequence = 0
					} else {
						tt.want.Sessions[i].Sequence = infos[i].Details.GetSequence()
					}
					tt.want.Sessions[i].CreationDate = infos[i].Details.GetChangeDate()
					tt.want.Sessions[i].ChangeDate = infos[i].Details.GetChangeDate()

					found := false
					for _, s := range got.Sessions {
						if s.Id == infos[i].ID {
							verifySession(ttt, s, tt.want.Sessions[i], tt.wantExpirationWindow, infos[i].UserID, tt.wantFactors...)
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

func TestServer_ListSessions(t *testing.T) {
	t.Parallel()
	runListSessionsTestCases(t, listSessionsEnv{
		ownerCtx:    IAMOwnerCTX,
		loginCtx:    LoginCTX,
		noPermCtx:   UserCTX,
		client:      Client,
		loginUserID: Instance.Users.Get(integration.UserTypeLogin).ID,
		testUser:    User,
		newUser:     func(ctx context.Context) *user.AddHumanUserResponse { return Instance.CreateHumanUser(ctx) },
	})
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

func TestServer_ListSessions_Relational(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	sysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(sysAuthZ)
	integration.EnsureInstanceFeature(t, sysAuthZ, inst,
		&featurepb.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)},
		func(tCollect *assert.CollectT, got *featurepb.GetInstanceFeaturesResponse) {
			assert.True(tCollect, got.GetEnableRelationalTables().GetEnabled())
		},
	)
	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(sysAuthZ, &instancepb.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	instOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	loginCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeLogin)
	noPermCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeNoPermission)

	runListSessionsTestCases(t, listSessionsEnv{
		ownerCtx:     instOwnerCtx,
		loginCtx:     loginCtx,
		noPermCtx:    noPermCtx,
		client:       inst.Client.SessionV2,
		loginUserID:  inst.Users.Get(integration.UserTypeLogin).ID,
		testUser:     inst.CreateHumanUser(instOwnerCtx),
		newUser:      func(ctx context.Context) *user.AddHumanUserResponse { return inst.CreateHumanUser(ctx) },
		isRelational: true,
	})
}
