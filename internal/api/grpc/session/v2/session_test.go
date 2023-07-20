package session

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
)

func Test_sessionsToPb(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)

	sessions := []*query.Session{
		{ // no factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			Metadata:      map[string][]byte{"hello": []byte("world")},
		},
		{ // user factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			UserFactor: query.SessionUserFactor{
				UserID:        "345",
				UserCheckedAt: past,
				LoginName:     "donald",
				DisplayName:   "donald duck",
				ResourceOwner: "org1",
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // password factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			UserFactor: query.SessionUserFactor{
				UserID:        "345",
				UserCheckedAt: past,
				LoginName:     "donald",
				DisplayName:   "donald duck",
				ResourceOwner: "org1",
			},
			PasswordFactor: query.SessionPasswordFactor{
				PasswordCheckedAt: past,
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // passkey factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			UserFactor: query.SessionUserFactor{
				UserID:        "345",
				UserCheckedAt: past,
				LoginName:     "donald",
				DisplayName:   "donald duck",
				ResourceOwner: "org1",
			},
			PasskeyFactor: query.SessionPasskeyFactor{
				PasskeyCheckedAt: past,
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
	}

	want := []*session.Session{
		{ // no factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors:      nil,
			Metadata:     map[string][]byte{"hello": []byte("world")},
		},
		{ // user factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors: &session.Factors{
				User: &session.UserFactor{
					VerifiedAt:     timestamppb.New(past),
					Id:             "345",
					LoginName:      "donald",
					DisplayName:    "donald duck",
					OrganisationId: "org1",
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // password factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors: &session.Factors{
				User: &session.UserFactor{
					VerifiedAt:     timestamppb.New(past),
					Id:             "345",
					LoginName:      "donald",
					DisplayName:    "donald duck",
					OrganisationId: "org1",
				},
				Password: &session.PasswordFactor{
					VerifiedAt: timestamppb.New(past),
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // passkey factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors: &session.Factors{
				User: &session.UserFactor{
					VerifiedAt:     timestamppb.New(past),
					Id:             "345",
					LoginName:      "donald",
					DisplayName:    "donald duck",
					OrganisationId: "org1",
				},
				Passkey: &session.PasskeyFactor{
					VerifiedAt: timestamppb.New(past),
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
	}

	out := sessionsToPb(sessions)
	require.Len(t, out, len(want))

	for i, got := range out {
		if !proto.Equal(got, want[i]) {
			t.Errorf("session %d got:\n%v\nwant:\n%v", i, got, want[i])
		}
	}
}

func mustNewTextQuery(t testing.TB, column query.Column, value string, compare query.TextComparison) query.SearchQuery {
	q, err := query.NewTextQuery(column, value, compare)
	require.NoError(t, err)
	return q
}

func mustNewListQuery(t testing.TB, column query.Column, list []any, compare query.ListComparison) query.SearchQuery {
	q, err := query.NewListQuery(query.SessionColumnID, list, compare)
	require.NoError(t, err)
	return q
}

func Test_listSessionsRequestToQuery(t *testing.T) {
	type args struct {
		ctx context.Context
		req *session.ListSessionsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *query.SessionsSearchQueries
		wantErr error
	}{
		{
			name: "default request",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				req: &session.ListSessionsRequest{},
			},
			want: &query.SessionsSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset: 0,
					Limit:  0,
					Asc:    false,
				},
				Queries: []query.SearchQuery{
					mustNewTextQuery(t, query.SessionColumnCreator, "789", query.TextEquals),
				},
			},
		},
		{
			name: "with list query and sessions",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				req: &session.ListSessionsRequest{
					Query: &object.ListQuery{
						Offset: 10,
						Limit:  20,
						Asc:    true,
					},
					Queries: []*session.SearchQuery{
						{Query: &session.SearchQuery_IdsQuery{
							IdsQuery: &session.IDsQuery{
								Ids: []string{"1", "2", "3"},
							},
						}},
						{Query: &session.SearchQuery_IdsQuery{
							IdsQuery: &session.IDsQuery{
								Ids: []string{"4", "5", "6"},
							},
						}},
					},
				},
			},
			want: &query.SessionsSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset: 10,
					Limit:  20,
					Asc:    true,
				},
				Queries: []query.SearchQuery{
					mustNewListQuery(t, query.SessionColumnID, []interface{}{"1", "2", "3"}, query.ListIn),
					mustNewListQuery(t, query.SessionColumnID, []interface{}{"4", "5", "6"}, query.ListIn),
					mustNewTextQuery(t, query.SessionColumnCreator, "789", query.TextEquals),
				},
			},
		},
		{
			name: "invalid argument error",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				req: &session.ListSessionsRequest{
					Query: &object.ListQuery{
						Offset: 10,
						Limit:  20,
						Asc:    true,
					},
					Queries: []*session.SearchQuery{
						{Query: &session.SearchQuery_IdsQuery{
							IdsQuery: &session.IDsQuery{
								Ids: []string{"1", "2", "3"},
							},
						}},
						{Query: nil},
					},
				},
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listSessionsRequestToQuery(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sessionQueriesToQuery(t *testing.T) {
	type args struct {
		ctx     context.Context
		queries []*session.SearchQuery
	}
	tests := []struct {
		name    string
		args    args
		want    []query.SearchQuery
		wantErr error
	}{
		{
			name: "creator only",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
			},
			want: []query.SearchQuery{
				mustNewTextQuery(t, query.SessionColumnCreator, "789", query.TextEquals),
			},
		},
		{
			name: "invalid argument",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				queries: []*session.SearchQuery{
					{Query: nil},
				},
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
		},
		{
			name: "creator and sessions",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				queries: []*session.SearchQuery{
					{Query: &session.SearchQuery_IdsQuery{
						IdsQuery: &session.IDsQuery{
							Ids: []string{"1", "2", "3"},
						},
					}},
					{Query: &session.SearchQuery_IdsQuery{
						IdsQuery: &session.IDsQuery{
							Ids: []string{"4", "5", "6"},
						},
					}},
				},
			},
			want: []query.SearchQuery{
				mustNewListQuery(t, query.SessionColumnID, []interface{}{"1", "2", "3"}, query.ListIn),
				mustNewListQuery(t, query.SessionColumnID, []interface{}{"4", "5", "6"}, query.ListIn),
				mustNewTextQuery(t, query.SessionColumnCreator, "789", query.TextEquals),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sessionQueriesToQuery(tt.args.ctx, tt.args.queries)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sessionQueryToQuery(t *testing.T) {
	type args struct {
		query *session.SearchQuery
	}
	tests := []struct {
		name    string
		args    args
		want    query.SearchQuery
		wantErr error
	}{
		{
			name: "invalid argument",
			args: args{&session.SearchQuery{
				Query: nil,
			}},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
		},
		{
			name: "query",
			args: args{&session.SearchQuery{
				Query: &session.SearchQuery_IdsQuery{
					IdsQuery: &session.IDsQuery{
						Ids: []string{"1", "2", "3"},
					},
				},
			}},
			want: mustNewListQuery(t, query.SessionColumnID, []interface{}{"1", "2", "3"}, query.ListIn),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sessionQueryToQuery(tt.args.query)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustUserLoginNamesSearchQuery(t testing.TB, value string) query.SearchQuery {
	loginNameQuery, err := query.NewUserLoginNamesSearchQuery("bar")
	require.NoError(t, err)
	return loginNameQuery
}

func Test_userCheck(t *testing.T) {
	type args struct {
		user *session.CheckUser
	}
	tests := []struct {
		name    string
		args    args
		want    userSearch
		wantErr error
	}{
		{
			name: "nil user",
			args: args{nil},
			want: nil,
		},
		{
			name: "by user id",
			args: args{&session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: "foo",
				},
			}},
			want: userSearchByID{"foo"},
		},
		{
			name: "by user id",
			args: args{&session.CheckUser{
				Search: &session.CheckUser_LoginName{
					LoginName: "bar",
				},
			}},
			want: userSearchByLoginName{mustUserLoginNamesSearchQuery(t, "bar")},
		},
		{
			name: "unimplemented error",
			args: args{&session.CheckUser{
				Search: nil,
			}},
			wantErr: caos_errs.ThrowUnimplementedf(nil, "SESSION-d3b4g0", "user search %T not implemented", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userCheck(tt.args.user)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
