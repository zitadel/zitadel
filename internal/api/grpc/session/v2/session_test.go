package session

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var (
	creationDate = time.Date(2023, 10, 10, 14, 15, 0, 0, time.UTC)
)

func Test_sessionsToPb(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)

	sessions := []*query.Session{
		{ // no factor, with user agent and expiration
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			Metadata:      map[string][]byte{"hello": []byte("world")},
			UserAgent: domain.UserAgent{
				FingerprintID: gu.Ptr("fingerprintID"),
				Description:   gu.Ptr("description"),
				IP:            net.IPv4(1, 2, 3, 4),
				Header:        http.Header{"foo": []string{"foo", "bar"}},
			},
			Expiration: now,
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
		{ // webAuthN factor
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
			WebAuthNFactor: query.SessionWebAuthNFactor{
				WebAuthNCheckedAt: past,
				UserVerified:      true,
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // totp factor
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
			TOTPFactor: query.SessionTOTPFactor{
				TOTPCheckedAt: past,
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
	}

	want := []*session.Session{
		{ // no factor, with user agent and expiration
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors:      nil,
			Metadata:     map[string][]byte{"hello": []byte("world")},
			UserAgent: &session.UserAgent{
				FingerprintId: gu.Ptr("fingerprintID"),
				Description:   gu.Ptr("description"),
				Ip:            gu.Ptr("1.2.3.4"),
				Header: map[string]*session.UserAgent_HeaderValues{
					"foo": {Values: []string{"foo", "bar"}},
				},
			},
			ExpirationDate: timestamppb.New(now),
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
					OrganizationId: "org1",
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
					OrganizationId: "org1",
				},
				Password: &session.PasswordFactor{
					VerifiedAt: timestamppb.New(past),
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // webAuthN factor
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
					OrganizationId: "org1",
				},
				WebAuthN: &session.WebAuthNFactor{
					VerifiedAt:   timestamppb.New(past),
					UserVerified: true,
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // totp factor
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
					OrganizationId: "org1",
				},
				Totp: &session.TOTPFactor{
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

func Test_userAgentToPb(t *testing.T) {
	type args struct {
		ua domain.UserAgent
	}
	tests := []struct {
		name string
		args args
		want *session.UserAgent
	}{
		{
			name: "empty",
			args: args{domain.UserAgent{}},
		},
		{
			name: "fingerprint id and description",
			args: args{domain.UserAgent{
				FingerprintID: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
			}},
			want: &session.UserAgent{
				FingerprintId: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
			},
		},
		{
			name: "with ip",
			args: args{domain.UserAgent{
				FingerprintID: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
				IP:            net.IPv4(1, 2, 3, 4),
			}},
			want: &session.UserAgent{
				FingerprintId: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
				Ip:            gu.Ptr("1.2.3.4"),
			},
		},
		{
			name: "with header",
			args: args{domain.UserAgent{
				FingerprintID: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
				Header: http.Header{
					"foo":   []string{"foo", "bar"},
					"hello": []string{"world"},
				},
			}},
			want: &session.UserAgent{
				FingerprintId: gu.Ptr("fingerPrintID"),
				Description:   gu.Ptr("description"),
				Header: map[string]*session.UserAgent_HeaderValues{
					"foo":   {Values: []string{"foo", "bar"}},
					"hello": {Values: []string{"world"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userAgentToPb(tt.args.ua)
			assert.Equal(t, tt.want, got)
		})
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

func mustNewTimestampQuery(t testing.TB, column query.Column, ts time.Time, compare query.TimestampComparison) query.SearchQuery {
	q, err := query.NewTimestampQuery(column, ts, compare)
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
			name: "default request with sorting column",
			args: args{
				ctx: authz.NewMockContext("123", "456", "789"),
				req: &session.ListSessionsRequest{
					SortingColumn: session.SessionFieldName_SESSION_FIELD_NAME_CREATION_DATE,
				},
			},
			want: &query.SessionsSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					SortingColumn: query.SessionColumnCreationDate,
					Asc:           false,
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
						{Query: &session.SearchQuery_UserIdQuery{
							UserIdQuery: &session.UserIDQuery{
								Id: "10",
							},
						}},
						{Query: &session.SearchQuery_CreationDateQuery{
							CreationDateQuery: &session.CreationDateQuery{
								CreationDate: timestamppb.New(creationDate),
								Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER,
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
					mustNewTextQuery(t, query.SessionColumnUserID, "10", query.TextEquals),
					mustNewTimestampQuery(t, query.SessionColumnCreationDate, creationDate, query.TimestampGreater),
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
			wantErr: zerrors.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
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
			wantErr: zerrors.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
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
			wantErr: zerrors.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid"),
		},
		{
			name: "ids query",
			args: args{&session.SearchQuery{
				Query: &session.SearchQuery_IdsQuery{
					IdsQuery: &session.IDsQuery{
						Ids: []string{"1", "2", "3"},
					},
				},
			}},
			want: mustNewListQuery(t, query.SessionColumnID, []interface{}{"1", "2", "3"}, query.ListIn),
		},
		{
			name: "user id query",
			args: args{&session.SearchQuery{
				Query: &session.SearchQuery_UserIdQuery{
					UserIdQuery: &session.UserIDQuery{
						Id: "10",
					},
				},
			}},
			want: mustNewTextQuery(t, query.SessionColumnUserID, "10", query.TextEquals),
		},
		{
			name: "creation date query",
			args: args{&session.SearchQuery{
				Query: &session.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session.CreationDateQuery{
						CreationDate: timestamppb.New(creationDate),
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS,
					},
				},
			}},
			want: mustNewTimestampQuery(t, query.SessionColumnCreationDate, creationDate, query.TimestampLess),
		},
		{
			name: "creation date query with default method",
			args: args{&session.SearchQuery{
				Query: &session.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session.CreationDateQuery{
						CreationDate: timestamppb.New(creationDate),
					},
				},
			}},
			want: mustNewTimestampQuery(t, query.SessionColumnCreationDate, creationDate, query.TimestampEquals),
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
			want: userSearchByLoginName{"bar"},
		},
		{
			name: "unimplemented error",
			args: args{&session.CheckUser{
				Search: nil,
			}},
			wantErr: zerrors.ThrowUnimplementedf(nil, "SESSION-d3b4g0", "user search %T not implemented", nil),
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

func Test_userVerificationRequirementToDomain(t *testing.T) {
	type args struct {
		req session.UserVerificationRequirement
	}
	tests := []struct {
		args args
		want domain.UserVerificationRequirement
	}{
		{
			args: args{session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_UNSPECIFIED},
			want: domain.UserVerificationRequirementUnspecified,
		},
		{
			args: args{session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED},
			want: domain.UserVerificationRequirementRequired,
		},
		{
			args: args{session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED},
			want: domain.UserVerificationRequirementPreferred,
		},
		{
			args: args{session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_DISCOURAGED},
			want: domain.UserVerificationRequirementDiscouraged,
		},
		{
			args: args{999},
			want: domain.UserVerificationRequirementUnspecified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.req.String(), func(t *testing.T) {
			got := userVerificationRequirementToDomain(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_userAgentToCommand(t *testing.T) {
	type args struct {
		userAgent *session.UserAgent
	}
	tests := []struct {
		name string
		args args
		want *domain.UserAgent
	}{
		{
			name: "nil",
			args: args{nil},
			want: nil,
		},
		{
			name: "all fields",
			args: args{&session.UserAgent{
				FingerprintId: gu.Ptr("fp1"),
				Ip:            gu.Ptr("1.2.3.4"),
				Description:   gu.Ptr("firefox"),
				Header: map[string]*session.UserAgent_HeaderValues{
					"hello": {
						Values: []string{"foo", "bar"},
					},
				},
			}},
			want: &domain.UserAgent{
				FingerprintID: gu.Ptr("fp1"),
				IP:            net.ParseIP("1.2.3.4"),
				Description:   gu.Ptr("firefox"),
				Header: http.Header{
					"hello": []string{"foo", "bar"},
				},
			},
		},
		{
			name: "invalid ip",
			args: args{&session.UserAgent{
				FingerprintId: gu.Ptr("fp1"),
				Ip:            gu.Ptr("oops"),
				Description:   gu.Ptr("firefox"),
				Header: map[string]*session.UserAgent_HeaderValues{
					"hello": {
						Values: []string{"foo", "bar"},
					},
				},
			}},
			want: &domain.UserAgent{
				FingerprintID: gu.Ptr("fp1"),
				IP:            nil,
				Description:   gu.Ptr("firefox"),
				Header: http.Header{
					"hello": []string{"foo", "bar"},
				},
			},
		},
		{
			name: "nil fields",
			args: args{&session.UserAgent{}},
			want: &domain.UserAgent{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userAgentToCommand(tt.args.userAgent)
			assert.Equal(t, tt.want, got)
		})
	}
}
