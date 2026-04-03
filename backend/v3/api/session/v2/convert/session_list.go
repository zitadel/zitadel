package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var grpcTimestampOpToDomain = map[objpb.TimestampQueryMethod]database.NumberOperation{
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS:            database.NumberOperationEqual,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER:           database.NumberOperationGreaterThan,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS: database.NumberOperationGreaterThanOrEqual,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS:              database.NumberOperationLessThan,
	objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS:    database.NumberOperationLessThanOrEqual,
}

// ListSessionsRequestGRPCToDomain converts a gRPC [session_grpc.ListSessionsRequest] to
// the domain [domain.ListSessionsRequest]. It validates the request and returns an error
// for any invalid filter values.
func ListSessionsRequestGRPCToDomain(req *session_grpc.ListSessionsRequest) (*domain.ListSessionsRequest, error) {
	filters, err := searchQueriesGRPCToDomain(req.GetQueries())
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

func searchQueriesGRPCToDomain(queries []*session_grpc.SearchQuery) ([]domain.SessionFilter, error) {
	if len(queries) == 0 {
		return nil, nil
	}
	filters := make([]domain.SessionFilter, 0, len(queries))
	for _, q := range queries {
		f, err := searchQueryGRPCToDomain(q)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func searchQueryGRPCToDomain(q *session_grpc.SearchQuery) (domain.SessionFilter, error) {
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
		var id *string
		if typedQuery.CreatorQuery != nil && typedQuery.CreatorQuery.Id != nil {
			if typedQuery.CreatorQuery.GetId() == "" {
				return nil, zerrors.ThrowInvalidArgument(nil, "DOM-x8n24uh", "List.Query.Invalid")
			}
			id = typedQuery.CreatorQuery.Id
		}
		return domain.SessionCreatorFilter{ID: id}, nil

	case *session_grpc.SearchQuery_UserAgentQuery:
		var fp *string
		if typedQuery.UserAgentQuery != nil && typedQuery.UserAgentQuery.FingerprintId != nil {
			if typedQuery.UserAgentQuery.GetFingerprintId() == "" {
				return nil, zerrors.ThrowInvalidArgument(nil, "DOM-x8n23uh", "List.Query.Invalid")
			}
			fp = typedQuery.UserAgentQuery.FingerprintId
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
func DomainSessionToGRPCResponse(s *domain.Session) *session_grpc.Session {
	pb := &session_grpc.Session{
		Id:           s.ID,
		CreationDate: timestamppb.New(s.CreatedAt),
		ChangeDate:   timestamppb.New(s.UpdatedAt),
		Factors:      domainFactorsToGRPC(s.Factors),
		UserAgent:    domainUserAgentToGRPC(s.UserAgent),
	}

	if len(s.Metadata) > 0 {
		pb.Metadata = make(map[string][]byte, len(s.Metadata))
		for _, m := range s.Metadata {
			pb.Metadata[m.Key] = m.Value
		}
	}

	if !s.Expiration.IsZero() {
		pb.ExpirationDate = timestamppb.New(s.Expiration)
	}

	return pb
}

func domainFactorsToGRPC(factors domain.SessionFactors) *session_grpc.Factors {
	userFactor := factors.GetUserFactor()
	if userFactor == nil {
		return nil
	}

	pb := &session_grpc.Factors{
		User: &session_grpc.UserFactor{
			VerifiedAt: timestamppb.New(userFactor.LastVerifiedAt),
			Id:         userFactor.UserID,
			// TODO(IAM-Marco): Should I fetch also the user = session.UserID and insert its data here?
		},
	}

	if f := factors.GetPasswordFactor(); f != nil {
		pb.Password = &session_grpc.PasswordFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := factors.GetPasskeyFactor(); f != nil {
		pb.WebAuthN = &session_grpc.WebAuthNFactor{
			VerifiedAt:   timestamppb.New(f.LastVerifiedAt),
			UserVerified: f.UserVerified,
		}
	}
	if f := factors.GetIDPIntentFactor(); f != nil {
		pb.Intent = &session_grpc.IntentFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := factors.GetTOTPFactor(); f != nil {
		pb.Totp = &session_grpc.TOTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := factors.GetOTPSMSFactor(); f != nil {
		pb.OtpSms = &session_grpc.OTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := factors.GetOTPEmailFactor(); f != nil {
		pb.OtpEmail = &session_grpc.OTPFactor{
			VerifiedAt: timestamppb.New(f.LastVerifiedAt),
		}
	}
	if f := factors.GetRecoveryCodeFactor(); f != nil {
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
