package convert

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

// ListSessionsRequestGRPCToDomain converts a gRPC [session_grpc.ListSessionsRequest] to
// the domain [domain.ListSessionsRequest]. It validates the request and returns an error
// for any invalid filter values.
func ListSessionsRequestGRPCToDomain(ctx context.Context, req *session_grpc.ListSessionsRequest) (*domain.ListSessionsRequest, error) {
	filters, err := searchQueriesGRPCToDomain(ctx, req.GetQueries())
	if err != nil {
		return nil, err
	}

	sortCol := domain.SessionSortColumnUnspecified
	if req.GetSortingColumn() == session_grpc.SessionFieldName_SESSION_FIELD_NAME_CREATION_DATE {
		sortCol = domain.SessionSortColumnCreationDate
	}

	return &domain.ListSessionsRequest{
		Filters:    filters,
		SortColumn: sortCol,
		Ascending:  req.GetQuery().GetAsc(),
		Limit:      req.GetQuery().GetLimit(),
		Offset:     uint32(req.GetQuery().GetOffset()),
	}, nil
}

func searchQueriesGRPCToDomain(ctx context.Context, queries []*session_grpc.SearchQuery) ([]domain.SessionFilter, error) {
	if len(queries) == 0 {
		return nil, nil
	}
	filters := make([]domain.SessionFilter, 0, len(queries))
	for _, q := range queries {
		f, err := searchQueryGRPCToDomain(ctx, q)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func searchQueryGRPCToDomain(ctx context.Context, q *session_grpc.SearchQuery) (domain.SessionFilter, error) {
	switch typedQuery := q.GetQuery().(type) {
	case *session_grpc.SearchQuery_IdsQuery:
		return domain.SessionIDsFilter{IDs: typedQuery.IdsQuery.GetIds()}, nil

	case *session_grpc.SearchQuery_UserIdQuery:
		return domain.SessionUserIDFilter{UserID: typedQuery.UserIdQuery.GetId()}, nil

	case *session_grpc.SearchQuery_CreationDateQuery:
		return domain.SessionCreationDateFilter{
			Op:   grpcTimestampOpToDomain[typedQuery.CreationDateQuery.GetMethod()],
			Date: typedQuery.CreationDateQuery.GetCreationDate().AsTime(),
		}, nil

	case *session_grpc.SearchQuery_CreatorQuery:
		id := authz.GetCtxData(ctx).UserID
		if typedQuery.CreatorQuery != nil && typedQuery.CreatorQuery.Id != nil {
			if typedQuery.CreatorQuery.GetId() == "" {
				return nil, zerrors.ThrowInvalidArgument(nil, "DOM-x8n24uh", "List.Query.Invalid")
			}
			id = typedQuery.CreatorQuery.GetId()
		}
		return domain.SessionCreatorFilter{ID: id}, nil

	case *session_grpc.SearchQuery_UserAgentQuery:
		fp := authz.GetCtxData(ctx).AgentID
		if typedQuery.UserAgentQuery != nil && typedQuery.UserAgentQuery.FingerprintId != nil {
			if typedQuery.UserAgentQuery.GetFingerprintId() == "" {
				return nil, zerrors.ThrowInvalidArgument(nil, "DOM-x8n23uh", "List.Query.Invalid")
			}
			fp = typedQuery.UserAgentQuery.GetFingerprintId()
		}
		return domain.SessionUserAgentFilter{FingerprintID: fp}, nil

	case *session_grpc.SearchQuery_ExpirationDateQuery:
		return domain.SessionExpirationDateFilter{
			Op:   grpcTimestampOpToDomain[typedQuery.ExpirationDateQuery.GetMethod()],
			Date: typedQuery.ExpirationDateQuery.GetExpirationDate().AsTime(),
		}, nil

	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "CONV-Cz5s3t", "session search query %T not implemented", typedQuery)
	}
}

// DomainSessionListToGRPCResponse converts a slice of domain sessions to gRPC response sessions.
func DomainSessionListToGRPCResponse(sessions []*domain.Session) []*session_grpc.Session {
	out := make([]*session_grpc.Session, len(sessions))
	for i, s := range sessions {
		out[i] = DomainSessionToGRPCResponse(s)
	}
	return out
}

// DomainSessionToGRPCResponse converts a single domain session to its gRPC representation.
func DomainSessionToGRPCResponse(session *domain.Session) *session_grpc.Session {
	pb := &session_grpc.Session{
		Id:           session.ID,
		CreationDate: timestamppb.New(session.CreatedAt),
		ChangeDate:   timestamppb.New(session.UpdatedAt),
		Factors:      domainFactorsToGRPC(session),
		UserAgent:    domainUserAgentToGRPC(session.UserAgent),
	}

	if len(session.Metadata) > 0 {
		pb.Metadata = make(map[string][]byte, len(session.Metadata))
		for _, m := range session.Metadata {
			pb.Metadata[m.Key] = m.Value
		}
	}

	if !session.Expiration.IsZero() {
		pb.ExpirationDate = timestamppb.New(session.Expiration)
	}

	return pb
}

func domainFactorsToGRPC(session *domain.Session) *session_grpc.Factors {
	userFactor := session.Factors.GetUserFactor()
	if userFactor == nil {
		return nil
	}

	pb := &session_grpc.Factors{
		User: &session_grpc.UserFactor{
			VerifiedAt: timestamppb.New(userFactor.LastVerifiedAt),
			Id:         userFactor.UserID,
		},
	}

	pb.User.LoginName = session.UserPreferredLoginName
	if session.UserOrganizationID != nil {
		pb.User.OrganizationId = *session.UserOrganizationID
	}
	if session.UserDisplayName != nil {
		pb.User.DisplayName = *session.UserDisplayName
	}

	if f := session.Factors.GetPasswordFactor(); f != nil {
		pb.Password = &session_grpc.PasswordFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := session.Factors.GetPasskeyFactor(); f != nil {
		pb.WebAuthN = &session_grpc.WebAuthNFactor{
			VerifiedAt:   timestamppb.New(f.LastVerifiedAt),
			UserVerified: f.UserVerified,
		}
	}
	if f := session.Factors.GetIDPIntentFactor(); f != nil {
		pb.Intent = &session_grpc.IntentFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := session.Factors.GetTOTPFactor(); f != nil {
		pb.Totp = &session_grpc.TOTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := session.Factors.GetOTPSMSFactor(); f != nil {
		pb.OtpSms = &session_grpc.OTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := session.Factors.GetOTPEmailFactor(); f != nil {
		pb.OtpEmail = &session_grpc.OTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := session.Factors.GetRecoveryCodeFactor(); f != nil {
		pb.RecoveryCode = &session_grpc.RecoveryCodeFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}

	return pb
}

func domainUserAgentToGRPC(ua *domain.SessionUserAgent) *session_grpc.UserAgent {
	if ua == nil {
		return nil
	}

	pb := &session_grpc.UserAgent{
		FingerprintId: ua.FingerprintID,
		Description:   ua.Description,
	}

	if ua.IP != nil {
		ip := ua.IP.String()
		pb.Ip = &ip
	}
	if len(ua.Header) > 0 {
		pb.Header = make(map[string]*session_grpc.UserAgent_HeaderValues, len(ua.Header))
		for k, v := range ua.Header {
			pb.Header[k] = &session_grpc.UserAgent_HeaderValues{Values: v}
		}
	}

	return pb
}
