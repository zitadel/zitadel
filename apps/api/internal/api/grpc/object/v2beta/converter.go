package object

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	org_pb "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails) *object.Details {
	details := &object.Details{
		Sequence:      objectDetail.Sequence,
		ResourceOwner: objectDetail.ResourceOwner,
	}
	if !objectDetail.EventDate.IsZero() {
		details.ChangeDate = timestamppb.New(objectDetail.EventDate)
	}
	if !objectDetail.CreationDate.IsZero() {
		details.CreationDate = timestamppb.New(objectDetail.CreationDate)
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

func ListQueryToModel(query *object.ListQuery) (offset, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Offset, uint64(query.Limit), query.Asc
}

func DomainsToPb(domains []*query.Domain) []*org_pb.Domain {
	d := make([]*org_pb.Domain, len(domains))
	for i, domain := range domains {
		d[i] = DomainToPb(domain)
	}
	return d
}

func DomainToPb(d *query.Domain) *org_pb.Domain {
	return &org_pb.Domain{
		OrganizationId: d.OrgID,
		DomainName:     d.Domain,
		IsVerified:     d.IsVerified,
		IsPrimary:      d.IsPrimary,
		ValidationType: DomainValidationTypeFromModel(d.ValidationType),
	}
}

func DomainValidationTypeFromModel(validationType domain.OrgDomainValidationType) org_pb.DomainValidationType {
	switch validationType {
	case domain.OrgDomainValidationTypeDNS:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS
	case domain.OrgDomainValidationTypeHTTP:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP
	case domain.OrgDomainValidationTypeUnspecified:
		// added to please golangci-lint
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED
	default:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED
	}
}

func DomainValidationTypeToDomain(validationType org_pb.DomainValidationType) domain.OrgDomainValidationType {
	switch validationType {
	case org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP:
		return domain.OrgDomainValidationTypeHTTP
	case org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS:
		return domain.OrgDomainValidationTypeDNS
	case org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED:
		// added to please golangci-lint
		return domain.OrgDomainValidationTypeUnspecified
	default:
		return domain.OrgDomainValidationTypeUnspecified
	}
}
