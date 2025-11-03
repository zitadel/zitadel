package session

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var (
	timestampComparisons = map[objpb.TimestampQueryMethod]query.TimestampComparison{
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS:            query.TimestampEquals,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER:           query.TimestampGreater,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS: query.TimestampGreaterOrEquals,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS:              query.TimestampLess,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS:    query.TimestampLessOrEquals,
	}
)

func (s *Server) GetSession(ctx context.Context, req *connect.Request[session.GetSessionRequest]) (*connect.Response[session.GetSessionResponse], error) {
	res, err := s.query.SessionByID(ctx, true, req.Msg.GetSessionId(), req.Msg.GetSessionToken(), s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&session.GetSessionResponse{
		Session: sessionToPb(res),
	}), nil
}

func (s *Server) ListSessions(ctx context.Context, req *connect.Request[session.ListSessionsRequest]) (*connect.Response[session.ListSessionsResponse], error) {
	queries, err := listSessionsRequestToQuery(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	sessions, err := s.query.SearchSessions(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&session.ListSessionsResponse{
		Details:  object.ToListDetails(sessions.SearchResponse),
		Sessions: sessionsToPb(sessions.Sessions),
	}), nil
}

func listSessionsRequestToQuery(ctx context.Context, req *session.ListSessionsRequest) (*query.SessionsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := sessionQueriesToQuery(ctx, req.GetQueries())
	if err != nil {
		return nil, err
	}
	return &query.SessionsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToSessionColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func sessionQueriesToQuery(ctx context.Context, queries []*session.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, v := range queries {
		q[i], err = sessionQueryToQuery(ctx, v)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func sessionQueryToQuery(ctx context.Context, sq *session.SearchQuery) (query.SearchQuery, error) {
	switch q := sq.Query.(type) {
	case *session.SearchQuery_IdsQuery:
		return idsQueryToQuery(q.IdsQuery)
	case *session.SearchQuery_UserIdQuery:
		return query.NewUserIDSearchQuery(q.UserIdQuery.GetId())
	case *session.SearchQuery_CreationDateQuery:
		return creationDateQueryToQuery(q.CreationDateQuery)
	case *session.SearchQuery_CreatorQuery:
		if q.CreatorQuery != nil && q.CreatorQuery.Id != nil {
			if q.CreatorQuery.GetId() != "" {
				return query.NewSessionCreatorSearchQuery(q.CreatorQuery.GetId())
			}
		} else {
			if userID := authz.GetCtxData(ctx).UserID; userID != "" {
				return query.NewSessionCreatorSearchQuery(userID)
			}
		}
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-x8n24uh", "List.Query.Invalid")
	case *session.SearchQuery_UserAgentQuery:
		if q.UserAgentQuery != nil && q.UserAgentQuery.FingerprintId != nil {
			if *q.UserAgentQuery.FingerprintId != "" {
				return query.NewSessionUserAgentFingerprintIDSearchQuery(q.UserAgentQuery.GetFingerprintId())
			}
		} else {
			if agentID := authz.GetCtxData(ctx).AgentID; agentID != "" {
				return query.NewSessionUserAgentFingerprintIDSearchQuery(agentID)
			}
		}
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-x8n23uh", "List.Query.Invalid")
	case *session.SearchQuery_ExpirationDateQuery:
		return expirationDateQueryToQuery(q.ExpirationDateQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid")
	}
}

func idsQueryToQuery(q *session.IDsQuery) (query.SearchQuery, error) {
	return query.NewSessionIDsSearchQuery(q.Ids)
}

func creationDateQueryToQuery(q *session.CreationDateQuery) (query.SearchQuery, error) {
	comparison := timestampComparisons[q.GetMethod()]
	return query.NewCreationDateQuery(q.GetCreationDate().AsTime(), comparison)
}

func expirationDateQueryToQuery(q *session.ExpirationDateQuery) (query.SearchQuery, error) {
	comparison := timestampComparisons[q.GetMethod()]

	// to obtain sessions with a set expiration date
	expirationDateQuery, err := query.NewExpirationDateQuery(q.GetExpirationDate().AsTime(), comparison)
	if err != nil {
		return nil, err
	}

	switch comparison {
	case query.TimestampEquals, query.TimestampLess, query.TimestampLessOrEquals:
		return expirationDateQuery, nil
	case query.TimestampGreater, query.TimestampGreaterOrEquals:
		// to obtain sessions without an expiration date
		expirationDateIsNullQuery, err := query.NewIsNullQuery(query.SessionColumnExpiration)
		if err != nil {
			return nil, err
		}
		return query.NewOrQuery(expirationDateQuery, expirationDateIsNullQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-Dwigt", "List.Query.InvalidComparisonMethod")
	}
}

func fieldNameToSessionColumn(field session.SessionFieldName) query.Column {
	switch field {
	case session.SessionFieldName_SESSION_FIELD_NAME_CREATION_DATE:
		return query.SessionColumnCreationDate
	case session.SessionFieldName_SESSION_FIELD_NAME_UNSPECIFIED:
		return query.Column{}
	default:
		return query.Column{}
	}
}

func sessionsToPb(sessions []*query.Session) []*session.Session {
	s := make([]*session.Session, len(sessions))
	for i, session := range sessions {
		s[i] = sessionToPb(session)
	}
	return s
}

func sessionToPb(s *query.Session) *session.Session {
	return &session.Session{
		Id:             s.ID,
		CreationDate:   timestamppb.New(s.CreationDate),
		ChangeDate:     timestamppb.New(s.ChangeDate),
		Sequence:       s.Sequence,
		Factors:        factorsToPb(s),
		Metadata:       s.Metadata,
		UserAgent:      userAgentToPb(s.UserAgent),
		ExpirationDate: expirationToPb(s.Expiration),
	}
}

func userAgentToPb(ua domain.UserAgent) *session.UserAgent {
	if ua.IsEmpty() {
		return nil
	}

	out := &session.UserAgent{
		FingerprintId: ua.FingerprintID,
		Description:   ua.Description,
	}
	if ua.IP != nil {
		out.Ip = gu.Ptr(ua.IP.String())
	}
	if ua.Header == nil {
		return out
	}
	out.Header = make(map[string]*session.UserAgent_HeaderValues, len(ua.Header))
	for k, v := range ua.Header {
		out.Header[k] = &session.UserAgent_HeaderValues{
			Values: v,
		}
	}
	return out
}

func expirationToPb(expiration time.Time) *timestamppb.Timestamp {
	if expiration.IsZero() {
		return nil
	}
	return timestamppb.New(expiration)
}

func factorsToPb(s *query.Session) *session.Factors {
	user := userFactorToPb(s.UserFactor)
	if user == nil {
		return nil
	}
	return &session.Factors{
		User:     user,
		Password: passwordFactorToPb(s.PasswordFactor),
		WebAuthN: webAuthNFactorToPb(s.WebAuthNFactor),
		Intent:   intentFactorToPb(s.IntentFactor),
		Totp:     totpFactorToPb(s.TOTPFactor),
		OtpSms:   otpFactorToPb(s.OTPSMSFactor),
		OtpEmail: otpFactorToPb(s.OTPEmailFactor),
	}
}

func passwordFactorToPb(factor query.SessionPasswordFactor) *session.PasswordFactor {
	if factor.PasswordCheckedAt.IsZero() {
		return nil
	}
	return &session.PasswordFactor{
		VerifiedAt: timestamppb.New(factor.PasswordCheckedAt),
	}
}

func intentFactorToPb(factor query.SessionIntentFactor) *session.IntentFactor {
	if factor.IntentCheckedAt.IsZero() {
		return nil
	}
	return &session.IntentFactor{
		VerifiedAt: timestamppb.New(factor.IntentCheckedAt),
	}
}

func webAuthNFactorToPb(factor query.SessionWebAuthNFactor) *session.WebAuthNFactor {
	if factor.WebAuthNCheckedAt.IsZero() {
		return nil
	}
	return &session.WebAuthNFactor{
		VerifiedAt:   timestamppb.New(factor.WebAuthNCheckedAt),
		UserVerified: factor.UserVerified,
	}
}

func totpFactorToPb(factor query.SessionTOTPFactor) *session.TOTPFactor {
	if factor.TOTPCheckedAt.IsZero() {
		return nil
	}
	return &session.TOTPFactor{
		VerifiedAt: timestamppb.New(factor.TOTPCheckedAt),
	}
}

func otpFactorToPb(factor query.SessionOTPFactor) *session.OTPFactor {
	if factor.OTPCheckedAt.IsZero() {
		return nil
	}
	return &session.OTPFactor{
		VerifiedAt: timestamppb.New(factor.OTPCheckedAt),
	}
}

func userFactorToPb(factor query.SessionUserFactor) *session.UserFactor {
	if factor.UserID == "" || factor.UserCheckedAt.IsZero() {
		return nil
	}
	return &session.UserFactor{
		VerifiedAt:     timestamppb.New(factor.UserCheckedAt),
		Id:             factor.UserID,
		LoginName:      factor.LoginName,
		DisplayName:    factor.DisplayName,
		OrganizationId: factor.ResourceOwner,
	}
}
