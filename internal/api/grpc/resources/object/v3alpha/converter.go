package object

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/pkg/grpc/resources/object/v3alpha"
)

func DomainToDetailsPb(objectDetail *domain.ObjectDetails, ownerType object.OwnerType, ownerId string) *resource_object.Details {
	details := &resource_object.Details{
		Id: objectDetail.ID,
		Owner: &object.Owner{
			Type: ownerType,
			Id:   ownerId,
		},
	}
	if !objectDetail.EventDate.IsZero() {
		details.Changed = timestamppb.New(objectDetail.EventDate)
	}
	if !objectDetail.CreationDate.IsZero() {
		details.Created = timestamppb.New(objectDetail.CreationDate)
	}
	return details
}

func ToSearchDetailsPb(request query.SearchRequest, response query.SearchResponse) *resource_object.ListDetails {
	details := &resource_object.ListDetails{
		AppliedLimit: request.Limit,
		TotalResult:  response.Count,
		Timestamp:    timestamppb.New(response.EventCreatedAt),
	}

	return details
}

func TextMethodPbToQuery(method resource_object.TextFilterMethod) query.TextComparison {
	switch method {
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS:
		return query.TextEquals
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS_IGNORE_CASE:
		return query.TextEqualsIgnoreCase
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH:
		return query.TextStartsWith
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH_IGNORE_CASE:
		return query.TextStartsWithIgnoreCase
	case resource_object.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS:
		return query.TextContains
	default:
		return -1
	}
}

func SearchQueryPbToQuery(defaults systemdefaults.SystemDefaults, query *resource_object.SearchQuery) (offset, limit uint64, asc bool, err error) {
	limit = defaults.DefaultQueryLimit
	asc = true
	if query == nil {
		return 0, limit, asc, nil
	}
	offset = query.Offset
	if query.Desc {
		asc = false
	}
	if defaults.MaxQueryLimit > 0 && uint64(query.Limit) > defaults.MaxQueryLimit {
		return 0, 0, false, zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", query.Limit, defaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded")
	}
	if query.Limit > 0 {
		limit = uint64(query.Limit)
	}
	return offset, limit, asc, nil
}

func ResourceOwnerFromOrganization(organization *object.Organization) string {
	if organization == nil {
		return ""
	}

	if organization.GetOrgId() != "" {
		return organization.GetOrgId()
	}
	if organization.GetOrgDomain() != "" {
		// TODO get org from domain
		return ""
	}
	return ""
}
