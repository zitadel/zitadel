package management

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func orgFromDomain(org *domain.Org) *management.Org {
	return &management.Org{
		ChangeDate:   timestamppb.New(org.ChangeDate),
		CreationDate: timestamppb.New(org.CreationDate),
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromDomain(org.State),
	}
}

func orgViewFromModel(org *org_model.OrgView) *management.OrgView {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &management.OrgView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromDomain(state domain.OrgState) management.OrgState {
	switch state {
	case domain.OrgStateActive:
		return management.OrgState_ORGSTATE_ACTIVE
	case domain.OrgStateInactive:
		return management.OrgState_ORGSTATE_INACTIVE
	default:
		return management.OrgState_ORGSTATE_UNSPECIFIED
	}
}

func orgStateFromModel(state org_model.OrgState) management.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return management.OrgState_ORGSTATE_ACTIVE
	case org_model.OrgStateInactive:
		return management.OrgState_ORGSTATE_INACTIVE
	default:
		return management.OrgState_ORGSTATE_UNSPECIFIED
	}
}

func addOrgDomainToDomain(ctx context.Context, orgDomain *management.AddOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: orgDomain.Domain,
	}
}

func orgDomainValidationToDomain(ctx context.Context, orgDomain *management.OrgDomainValidationRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain:         orgDomain.Domain,
		ValidationType: orgDomainValidationTypeToDomain(orgDomain.Type),
	}
}

func validateOrgDomainToDomain(ctx context.Context, orgDomain *management.ValidateOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: orgDomain.Domain,
	}
}

func orgDomainValidationTypeToDomain(validationType management.OrgDomainValidationType) domain.OrgDomainValidationType {
	switch validationType {
	case management.OrgDomainValidationType_ORGDOMAINVALIDATIONTYPE_HTTP:
		return domain.OrgDomainValidationTypeHTTP
	case management.OrgDomainValidationType_ORGDOMAINVALIDATIONTYPE_DNS:
		return domain.OrgDomainValidationTypeDNS
	default:
		return domain.OrgDomainValidationTypeUnspecified
	}
}

func orgDomainValidationTypeFromModel(key org_model.OrgDomainValidationType) management.OrgDomainValidationType {
	switch key {
	case org_model.OrgDomainValidationTypeHTTP:
		return management.OrgDomainValidationType_ORGDOMAINVALIDATIONTYPE_HTTP
	case org_model.OrgDomainValidationTypeDNS:
		return management.OrgDomainValidationType_ORGDOMAINVALIDATIONTYPE_DNS
	default:
		return management.OrgDomainValidationType_ORGDOMAINVALIDATIONTYPE_UNSPECIFIED
	}
}

func primaryOrgDomainToDomain(ctx context.Context, ordDomain *management.PrimaryOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: ordDomain.Domain,
	}
}

func removeOrgDomainToDomain(ctx context.Context, ordDomain *management.RemoveOrgDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		Domain: ordDomain.Domain,
	}
}

func orgDomainFromDomain(orgDomain *domain.OrgDomain) *management.OrgDomain {
	return &management.OrgDomain{
		ChangeDate:   timestamppb.New(orgDomain.ChangeDate),
		CreationDate: timestamppb.New(orgDomain.CreationDate),
		OrgId:        orgDomain.AggregateID,
		Domain:       orgDomain.Domain,
		Verified:     orgDomain.Verified,
		Primary:      orgDomain.Primary,
	}
}

func orgDomainViewFromModel(domain *org_model.OrgDomainView) *management.OrgDomainView {
	creationDate, err := ptypes.TimestampProto(domain.CreationDate)
	logging.Log("GRPC-7sjDs").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(domain.ChangeDate)
	logging.Log("GRPC-8iSji").OnError(err).Debug("unable to get timestamp from time")

	return &management.OrgDomainView{
		ChangeDate:     changeDate,
		CreationDate:   creationDate,
		OrgId:          domain.OrgID,
		Domain:         domain.Domain,
		Verified:       domain.Verified,
		Primary:        domain.Primary,
		ValidationType: orgDomainValidationTypeFromModel(domain.ValidationType),
	}
}

func orgDomainSearchRequestToModel(request *management.OrgDomainSearchRequest) *org_model.OrgDomainSearchRequest {
	return &org_model.OrgDomainSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: orgDomainSearchQueriesToModel(request.Queries),
	}
}

func orgDomainSearchQueriesToModel(queries []*management.OrgDomainSearchQuery) []*org_model.OrgDomainSearchQuery {
	modelQueries := make([]*org_model.OrgDomainSearchQuery, len(queries))

	for i, query := range queries {
		modelQueries[i] = orgDomainSearchQueryToModel(query)
	}

	return modelQueries
}

func orgDomainSearchQueryToModel(query *management.OrgDomainSearchQuery) *org_model.OrgDomainSearchQuery {
	return &org_model.OrgDomainSearchQuery{
		Key:    orgDomainSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func orgDomainSearchKeyToModel(key management.OrgDomainSearchKey) org_model.OrgDomainSearchKey {
	switch key {
	case management.OrgDomainSearchKey_ORGDOMAINSEARCHKEY_DOMAIN:
		return org_model.OrgDomainSearchKeyDomain
	default:
		return org_model.OrgDomainSearchKeyUnspecified
	}
}

func orgDomainSearchResponseFromModel(resp *org_model.OrgDomainSearchResponse) *management.OrgDomainSearchResponse {
	timestamp, err := ptypes.TimestampProto(resp.Timestamp)
	logging.Log("GRPC-Mxi9w").OnError(err).Debug("unable to get timestamp from time")
	return &management.OrgDomainSearchResponse{
		Limit:             resp.Limit,
		Offset:            resp.Offset,
		TotalResult:       resp.TotalResult,
		Result:            orgDomainsFromModel(resp.Result),
		ProcessedSequence: resp.Sequence,
		ViewTimestamp:     timestamp,
	}
}
func orgDomainsFromModel(viewDomains []*org_model.OrgDomainView) []*management.OrgDomainView {
	domains := make([]*management.OrgDomainView, len(viewDomains))

	for i, domain := range viewDomains {
		domains[i] = orgDomainViewFromModel(domain)
	}

	return domains
}

func orgChangesToResponse(response *org_model.OrgChanges, offset uint64, limit uint64) (_ *management.Changes) {
	return &management.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: orgChangesToMgtAPI(response),
	}
}

func orgChangesToMgtAPI(changes *org_model.OrgChanges) (_ []*management.Change) {
	result := make([]*management.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &management.Change{
			ChangeDate: change.ChangeDate,
			EventType:  message.NewLocalizedEventType(change.EventType),
			Sequence:   change.Sequence,
			Data:       data,
			Editor:     change.ModifierName,
			EditorId:   change.ModifierId,
		}
	}

	return result
}

func orgIamPolicyViewFromModel(policy *iam_model.OrgIAMPolicyView) *management.OrgIamPolicyView {
	return &management.OrgIamPolicyView{
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		Default:               policy.Default,
	}
}
