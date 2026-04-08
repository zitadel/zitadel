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
func DomainSessionListToGRPCResponse(sessions []domain.ListSessionResponse) []*session_grpc.Session {
	out := make([]*session_grpc.Session, len(sessions))
	for i, s := range sessions {
		out[i] = DomainSessionToGRPCResponse(s)
	}
	return out
}

// DomainSessionToGRPCResponse converts a single domain session to its gRPC representation.
func DomainSessionToGRPCResponse(lsr domain.ListSessionResponse) *session_grpc.Session {
	pb := &session_grpc.Session{
		Id:           lsr.Session.ID,
		CreationDate: timestamppb.New(lsr.Session.CreatedAt),
		ChangeDate:   timestamppb.New(lsr.Session.UpdatedAt),
		Factors:      domainFactorsToGRPC(lsr.Session.Factors, lsr.User),
		UserAgent:    domainUserAgentToGRPC(lsr.Session.UserAgent),
	}

	if len(lsr.Session.Metadata) > 0 {
		pb.Metadata = make(map[string][]byte, len(lsr.Session.Metadata))
		for _, m := range lsr.Session.Metadata {
			pb.Metadata[m.Key] = m.Value
		}
	}

	if !lsr.Session.Expiration.IsZero() {
		pb.ExpirationDate = timestamppb.New(lsr.Session.Expiration)
	}

	return pb
}

func domainFactorsToGRPC(factors domain.SessionFactors, user *domain.User) *session_grpc.Factors {
	userFactor := factors.GetUserFactor()
	if userFactor == nil {
		return nil
	}

	pb := &session_grpc.Factors{
		User: &session_grpc.UserFactor{
			VerifiedAt: timestamppb.New(userFactor.LastVerifiedAt),
			Id:         userFactor.UserID,
		},
	}

	if user != nil {
		// Set the first loginname in the list, regardless the preferred boolean
		// In case no preferred is found, at least we have one loginname
		preferredLoginName := user.LoginNames[0].LoginName
		for _, ln := range user.LoginNames {
			if ln.IsPreferred {
				preferredLoginName = ln.LoginName
				break
			}
		}

		pb.User.LoginName = preferredLoginName
		pb.User.OrganizationId = user.OrganizationID
		if user.Human != nil {
			pb.User.DisplayName = user.Human.DisplayName
		}
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
