package convert

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	objv2 "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestSearchQueryGRPCToDomain(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	nowPB := timestamppb.New(now)
	creatorID := "creator-123"
	fingerprintID := "fp-123"
	emptyStr := ""

	tests := []struct {
		name    string
		input   *session_grpc.SearchQuery
		want    domain.SessionFilter
		wantErr error
	}{
		{
			name: "IdsQuery maps to SessionIDsFilter",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_IdsQuery{
					IdsQuery: &session_grpc.IDsQuery{Ids: []string{"id1", "id2"}},
				},
			},
			want: domain.SessionIDsFilter{IDs: []string{"id1", "id2"}},
		},
		{
			name: "UserIdQuery maps to SessionUserIDFilter",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_UserIdQuery{
					UserIdQuery: &session_grpc.UserIDQuery{Id: "user-123"},
				},
			},
			want: domain.SessionUserIDFilter{UserID: "user-123"},
		},
		{
			name: "CreationDateQuery with EQUALS method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session_grpc.CreationDateQuery{
						CreationDate: nowPB,
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS,
					},
				},
			},
			want: domain.SessionCreationDateFilter{Op: database.NumberOperationEqual, Date: now},
		},
		{
			name: "CreationDateQuery with GREATER method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session_grpc.CreationDateQuery{
						CreationDate: nowPB,
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER,
					},
				},
			},
			want: domain.SessionCreationDateFilter{Op: database.NumberOperationGreaterThan, Date: now},
		},
		{
			name: "CreationDateQuery with GREATER_OR_EQUALS method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session_grpc.CreationDateQuery{
						CreationDate: nowPB,
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS,
					},
				},
			},
			want: domain.SessionCreationDateFilter{Op: database.NumberOperationGreaterThanOrEqual, Date: now},
		},
		{
			name: "CreationDateQuery with LESS method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session_grpc.CreationDateQuery{
						CreationDate: nowPB,
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS,
					},
				},
			},
			want: domain.SessionCreationDateFilter{Op: database.NumberOperationLessThan, Date: now},
		},
		{
			name: "CreationDateQuery with LESS_OR_EQUALS method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreationDateQuery{
					CreationDateQuery: &session_grpc.CreationDateQuery{
						CreationDate: nowPB,
						Method:       objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS,
					},
				},
			},
			want: domain.SessionCreationDateFilter{Op: database.NumberOperationLessThanOrEqual, Date: now},
		},
		{
			name: "CreatorQuery with nil inner query uses nil ID",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreatorQuery{CreatorQuery: nil},
			},
			want: domain.SessionCreatorFilter{ID: nil},
		},
		{
			name: "CreatorQuery with nil ID uses nil ID",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreatorQuery{
					CreatorQuery: &session_grpc.CreatorQuery{},
				},
			},
			want: domain.SessionCreatorFilter{ID: nil},
		},
		{
			name: "CreatorQuery with empty ID returns error",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreatorQuery{
					CreatorQuery: &session_grpc.CreatorQuery{Id: &emptyStr},
				},
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-x8n24uh", "List.Query.Invalid"),
		},
		{
			name: "CreatorQuery with valid ID uses it",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_CreatorQuery{
					CreatorQuery: &session_grpc.CreatorQuery{Id: &creatorID},
				},
			},
			want: domain.SessionCreatorFilter{ID: &creatorID},
		},
		{
			name: "UserAgentQuery with nil inner query uses nil FingerprintID",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_UserAgentQuery{UserAgentQuery: nil},
			},
			want: domain.SessionUserAgentFilter{FingerprintID: nil},
		},
		{
			name: "UserAgentQuery with nil FingerprintID uses nil FingerprintID",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_UserAgentQuery{
					UserAgentQuery: &session_grpc.UserAgentQuery{},
				},
			},
			want: domain.SessionUserAgentFilter{FingerprintID: nil},
		},
		{
			name: "UserAgentQuery with empty FingerprintID returns error",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_UserAgentQuery{
					UserAgentQuery: &session_grpc.UserAgentQuery{FingerprintId: &emptyStr},
				},
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-x8n23uh", "List.Query.Invalid"),
		},
		{
			name: "UserAgentQuery with valid FingerprintID uses it",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_UserAgentQuery{
					UserAgentQuery: &session_grpc.UserAgentQuery{FingerprintId: &fingerprintID},
				},
			},
			want: domain.SessionUserAgentFilter{FingerprintID: &fingerprintID},
		},
		{
			name: "ExpirationDateQuery with EQUALS method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_ExpirationDateQuery{
					ExpirationDateQuery: &session_grpc.ExpirationDateQuery{
						ExpirationDate: nowPB,
						Method:         objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS,
					},
				},
			},
			want: domain.SessionExpirationDateFilter{Op: database.NumberOperationEqual, Date: now},
		},
		{
			name: "ExpirationDateQuery with GREATER method",
			input: &session_grpc.SearchQuery{
				Query: &session_grpc.SearchQuery_ExpirationDateQuery{
					ExpirationDateQuery: &session_grpc.ExpirationDateQuery{
						ExpirationDate: nowPB,
						Method:         objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER,
					},
				},
			},
			want: domain.SessionExpirationDateFilter{Op: database.NumberOperationGreaterThan, Date: now},
		},
		{
			name:    "unknown query type returns error",
			input:   &session_grpc.SearchQuery{},
			wantErr: zerrors.ThrowInvalidArgumentf(nil, "CONV-Cz5s3t", "session search query %T not implemented", nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := searchQueryGRPCToDomain(tc.input)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSearchQueriesGRPCToDomain(t *testing.T) {
	t.Parallel()

	emptyStr := ""

	tests := []struct {
		name    string
		input   []*session_grpc.SearchQuery
		want    []domain.SessionFilter
		wantErr bool
	}{
		{
			name:  "nil slice returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty slice returns nil",
			input: []*session_grpc.SearchQuery{},
			want:  nil,
		},
		{
			name: "valid queries are converted to filters",
			input: []*session_grpc.SearchQuery{
				{Query: &session_grpc.SearchQuery_IdsQuery{IdsQuery: &session_grpc.IDsQuery{Ids: []string{"id1"}}}},
				{Query: &session_grpc.SearchQuery_UserIdQuery{UserIdQuery: &session_grpc.UserIDQuery{Id: "user-1"}}},
			},
			want: []domain.SessionFilter{
				domain.SessionIDsFilter{IDs: []string{"id1"}},
				domain.SessionUserIDFilter{UserID: "user-1"},
			},
		},
		{
			name: "invalid query inside list propagates error",
			input: []*session_grpc.SearchQuery{
				{
					Query: &session_grpc.SearchQuery_CreatorQuery{
						CreatorQuery: &session_grpc.CreatorQuery{Id: &emptyStr},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := searchQueriesGRPCToDomain(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestListSessionsRequestGRPCToDomain(t *testing.T) {
	t.Parallel()

	emptyStr := ""

	tests := []struct {
		name    string
		input   *session_grpc.ListSessionsRequest
		want    *domain.ListSessionsRequest
		wantErr bool
	}{
		{
			name:  "nil request returns empty domain request",
			input: nil,
			want: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnUnspecified,
			},
		},
		{
			name: "invalid query propagates error",
			input: &session_grpc.ListSessionsRequest{
				Queries: []*session_grpc.SearchQuery{
					{
						Query: &session_grpc.SearchQuery_CreatorQuery{
							CreatorQuery: &session_grpc.CreatorQuery{Id: &emptyStr},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "SESSION_FIELD_NAME_CREATION_DATE maps to SessionSortColumnCreationDate",
			input: &session_grpc.ListSessionsRequest{
				SortingColumn: session_grpc.SessionFieldName_SESSION_FIELD_NAME_CREATION_DATE,
			},
			want: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnCreationDate,
			},
		},
		{
			name: "SESSION_FIELD_NAME_UNSPECIFIED maps to SessionSortColumnUnspecified",
			input: &session_grpc.ListSessionsRequest{
				SortingColumn: session_grpc.SessionFieldName_SESSION_FIELD_NAME_UNSPECIFIED,
			},
			want: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnUnspecified,
			},
		},
		{
			name: "query limit, offset and ascending are mapped",
			input: &session_grpc.ListSessionsRequest{
				Query: &objv2.ListQuery{Limit: 20, Offset: 5, Asc: true},
			},
			want: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnUnspecified,
				Limit:      20,
				Offset:     5,
				Ascending:  true,
			},
		},
		{
			name: "valid queries are included in result",
			input: &session_grpc.ListSessionsRequest{
				Queries: []*session_grpc.SearchQuery{
					{Query: &session_grpc.SearchQuery_UserIdQuery{UserIdQuery: &session_grpc.UserIDQuery{Id: "user-1"}}},
				},
			},
			want: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnUnspecified,
				Filters:    []domain.SessionFilter{domain.SessionUserIDFilter{UserID: "user-1"}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ListSessionsRequestGRPCToDomain(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDomainSessionListToGRPCResponse(t *testing.T) {
	t.Parallel()

	t.Run("empty slice returns empty response slice", func(t *testing.T) {
		t.Parallel()
		got := DomainSessionListToGRPCResponse([]*domain.Session{})
		assert.Empty(t, got)
	})

	t.Run("each session is converted", func(t *testing.T) {
		t.Parallel()
		sessions := []*domain.Session{
			{ID: "sess-1"},
			{ID: "sess-2"},
		}
		got := DomainSessionListToGRPCResponse(sessions)
		assert.Len(t, got, 2)
		assert.Equal(t, "sess-1", got[0].Id)
		assert.Equal(t, "sess-2", got[1].Id)
	})
}

func TestDomainSessionToGRPCResponse(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	expiration := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	t.Run("basic fields are mapped", func(t *testing.T) {
		t.Parallel()
		s := &domain.Session{
			ID:        "sess-abc",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		got := DomainSessionToGRPCResponse(s)
		assert.Equal(t, "sess-abc", got.Id)
		assert.Equal(t, timestamppb.New(createdAt), got.CreationDate)
		assert.Equal(t, timestamppb.New(updatedAt), got.ChangeDate)
		assert.EqualValues(t, 0, got.Sequence)
		assert.Nil(t, got.ExpirationDate)
		assert.Nil(t, got.Metadata)
	})

	t.Run("expiration is set when non-zero", func(t *testing.T) {
		t.Parallel()
		s := &domain.Session{Expiration: expiration}
		got := DomainSessionToGRPCResponse(s)
		assert.Equal(t, timestamppb.New(expiration), got.ExpirationDate)
	})

	t.Run("zero expiration is not set", func(t *testing.T) {
		t.Parallel()
		s := &domain.Session{}
		got := DomainSessionToGRPCResponse(s)
		assert.Nil(t, got.ExpirationDate)
	})

	t.Run("metadata is mapped to key-value pairs", func(t *testing.T) {
		t.Parallel()
		s := &domain.Session{
			Metadata: []*domain.SessionMetadata{
				{Metadata: domain.Metadata{Key: "k1", Value: []byte("v1")}},
				{Metadata: domain.Metadata{Key: "k2", Value: []byte("v2")}},
			},
		}
		got := DomainSessionToGRPCResponse(s)
		assert.Equal(t, map[string][]byte{"k1": []byte("v1"), "k2": []byte("v2")}, got.Metadata)
	})

	t.Run("empty metadata is not set", func(t *testing.T) {
		t.Parallel()
		s := &domain.Session{Metadata: []*domain.SessionMetadata{}}
		got := DomainSessionToGRPCResponse(s)
		assert.Nil(t, got.Metadata)
	})
}

func TestDomainFactorsToGRPC(t *testing.T) {
	t.Parallel()

	verifiedAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	verifiedAtPB := timestamppb.New(verifiedAt)

	t.Run("no user factor returns nil", func(t *testing.T) {
		t.Parallel()
		got := domainFactorsToGRPC(domain.SessionFactors{})
		assert.Nil(t, got)
	})

	t.Run("no factors returns nil", func(t *testing.T) {
		t.Parallel()
		got := domainFactorsToGRPC(nil)
		assert.Nil(t, got)
	})

	t.Run("only user factor maps User field", func(t *testing.T) {
		t.Parallel()
		factors := domain.SessionFactors{
			&domain.SessionFactorUser{UserID: "user-1", LastVerifiedAt: verifiedAt},
		}
		got := domainFactorsToGRPC(factors)
		assert.NotNil(t, got)
		assert.Equal(t, "user-1", got.User.Id)
		assert.Equal(t, verifiedAtPB, got.User.VerifiedAt)
		assert.Nil(t, got.Password)
		assert.Nil(t, got.WebAuthN)
		assert.Nil(t, got.Intent)
		assert.Nil(t, got.Totp)
		assert.Nil(t, got.OtpSms)
		assert.Nil(t, got.OtpEmail)
		assert.Nil(t, got.RecoveryCode)
	})

	t.Run("all factors are mapped", func(t *testing.T) {
		t.Parallel()
		factors := domain.SessionFactors{
			&domain.SessionFactorUser{UserID: "user-1", LastVerifiedAt: verifiedAt},
			&domain.SessionFactorPassword{LastVerifiedAt: verifiedAt},
			&domain.SessionFactorPasskey{LastVerifiedAt: verifiedAt, UserVerified: true},
			&domain.SessionFactorIdentityProviderIntent{LastVerifiedAt: verifiedAt},
			&domain.SessionFactorTOTP{LastVerifiedAt: verifiedAt},
			&domain.SessionFactorOTPSMS{LastVerifiedAt: verifiedAt},
			&domain.SessionFactorOTPEmail{LastVerifiedAt: verifiedAt},
			&domain.SessionFactorRecoveryCode{LastVerifiedAt: verifiedAt},
		}
		got := domainFactorsToGRPC(factors)
		assert.NotNil(t, got)
		assert.Equal(t, verifiedAtPB, got.User.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.Password.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.WebAuthN.VerifiedAt)
		assert.True(t, got.WebAuthN.UserVerified)
		assert.Equal(t, verifiedAtPB, got.Intent.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.Totp.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.OtpSms.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.OtpEmail.VerifiedAt)
		assert.Equal(t, verifiedAtPB, got.RecoveryCode.VerifiedAt)
	})
}

func TestDomainUserAgentToGRPC(t *testing.T) {
	t.Parallel()

	fpID := "fp-abc"
	desc := "Chrome on macOS"
	ip := net.ParseIP("192.168.1.1")
	ipStr := ip.String()

	t.Run("nil user agent returns nil", func(t *testing.T) {
		t.Parallel()
		got := domainUserAgentToGRPC(nil)
		assert.Nil(t, got)
	})

	t.Run("empty user agent returns empty proto", func(t *testing.T) {
		t.Parallel()
		got := domainUserAgentToGRPC(&domain.SessionUserAgent{})
		assert.NotNil(t, got)
		assert.Nil(t, got.FingerprintId)
		assert.Nil(t, got.Description)
		assert.Nil(t, got.Ip)
		assert.Nil(t, got.Header)
	})

	t.Run("fingerprint ID is mapped", func(t *testing.T) {
		t.Parallel()
		got := domainUserAgentToGRPC(&domain.SessionUserAgent{FingerprintID: &fpID})
		assert.Equal(t, &fpID, got.FingerprintId)
	})

	t.Run("description is mapped", func(t *testing.T) {
		t.Parallel()
		got := domainUserAgentToGRPC(&domain.SessionUserAgent{Description: &desc})
		assert.Equal(t, &desc, got.Description)
	})

	t.Run("IP is mapped as string", func(t *testing.T) {
		t.Parallel()
		got := domainUserAgentToGRPC(&domain.SessionUserAgent{IP: ip})
		assert.Equal(t, &ipStr, got.Ip)
	})

	t.Run("headers are mapped", func(t *testing.T) {
		t.Parallel()
		ua := &domain.SessionUserAgent{
			Header: http.Header{"Accept": []string{"application/json", "text/html"}},
		}
		got := domainUserAgentToGRPC(ua)
		assert.Equal(t, []string{"application/json", "text/html"}, got.Header["Accept"].Values)
	})
}
