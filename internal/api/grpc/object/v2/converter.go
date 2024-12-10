package object

import (
	"context"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails) *object.Details {
	details := &object.Details{
		Sequence:      objectDetail.Sequence,
		ResourceOwner: objectDetail.ResourceOwner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	return details
}

func ToListDetails(response query.SearchResponse) *object.ListDetails {
	details := &object.ListDetails{
		TotalResult:       response.Count,
		ProcessedSequence: response.Sequence,
		Timestamp:         timestamppb.New(response.EventCreatedAt),
	}

	return details
}
func ListQueryToQuery(query *object.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Offset, uint64(query.Limit), query.Asc
}

func ResourceOwnerFromReq(ctx context.Context, req *object.RequestContext) string {
	if req.GetInstance() {
		return authz.GetInstance(ctx).InstanceID()
	}
	if req.GetOrgId() != "" {
		return req.GetOrgId()
	}
	return authz.GetCtxData(ctx).OrgID
}

func TextMethodToQuery(method object.TextQueryMethod) query.TextComparison {
	switch method {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return query.TextEquals
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return query.TextEqualsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return query.TextStartsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return query.TextStartsWithIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return query.TextContains
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return query.TextContainsIgnoreCase
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return query.TextEndsWith
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return query.TextEndsWithIgnoreCase
	default:
		return -1
	}
}

func AuthMethodsToPb(mfas *query.AuthMethods) []*user_pb.AuthFactor {
	factors := make([]*user_pb.AuthFactor, len(mfas.AuthMethods))
	for i, mfa := range mfas.AuthMethods {
		factors[i] = AuthMethodToPb(mfa)
	}
	return factors
}

func AuthMethodToPb(mfa *query.AuthMethod) *user_pb.AuthFactor {
	factor := &user_pb.AuthFactor{
		State: MFAStateToPb(mfa.State),
	}
	switch mfa.Type {
	case domain.UserAuthMethodTypeTOTP:
		factor.Type = &user_pb.AuthFactor_Otp{
			Otp: &user_pb.AuthFactorOTP{},
		}
	case domain.UserAuthMethodTypeU2F:
		factor.Type = &user_pb.AuthFactor_U2F{
			U2F: &user_pb.AuthFactorU2F{
				Id:   mfa.TokenID,
				Name: mfa.Name,
			},
		}
	case domain.UserAuthMethodTypeOTPSMS:
		factor.Type = &user_pb.AuthFactor_OtpSms{
			OtpSms: &user_pb.AuthFactorOTPSMS{},
		}
	case domain.UserAuthMethodTypeOTPEmail:
		factor.Type = &user_pb.AuthFactor_OtpEmail{
			OtpEmail: &user_pb.AuthFactorOTPEmail{},
		}
	}
	return factor
}

func MFAStateToPb(state domain.MFAState) user_pb.AuthFactorState {
	switch state {
	case domain.MFAStateNotReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY
	case domain.MFAStateReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_READY
	default:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_UNSPECIFIED
	}
}
