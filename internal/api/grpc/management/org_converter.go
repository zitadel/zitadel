package management

import (
	"encoding/json"

	"github.com/caos/logging"

	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func orgsFromModel(orgs []*org_model.Org) []*grpc.Org {
	orgList := make([]*grpc.Org, len(orgs))
	for i, org := range orgs {
		orgList[i] = orgFromModel(org)
	}
	return orgList
}

func orgFromModel(org *org_model.Org) *grpc.Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &grpc.Org{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgViewFromModel(org *org_model.OrgView) *grpc.OrgView {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &grpc.OrgView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromModel(state org_model.OrgState) grpc.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return grpc.OrgState_ORGSTATE_ACTIVE
	case org_model.OrgStateInactive:
		return grpc.OrgState_ORGSTATE_INACTIVE
	default:
		return grpc.OrgState_ORGSTATE_UNSPECIFIED
	}
}

func addOrgDomainToModel(domain *grpc.AddOrgDomainRequest) *org_model.OrgDomain {
	return &org_model.OrgDomain{Domain: domain.Domain}
}

func orgDomainFromModel(domain *org_model.OrgDomain) *grpc.OrgDomain {
	creationDate, err := ptypes.TimestampProto(domain.CreationDate)
	logging.Log("GRPC-u8Ksj").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(domain.ChangeDate)
	logging.Log("GRPC-9osFS").OnError(err).Debug("unable to get timestamp from time")

	return &grpc.OrgDomain{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		OrgId:        domain.AggregateID,
		Domain:       domain.Domain,
		Verified:     domain.Verified,
		Primary:      domain.Primary,
	}
}

func orgDomainViewFromModel(domain *org_model.OrgDomainView) *grpc.OrgDomainView {
	creationDate, err := ptypes.TimestampProto(domain.CreationDate)
	logging.Log("GRPC-7sjDs").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(domain.ChangeDate)
	logging.Log("GRPC-8iSji").OnError(err).Debug("unable to get timestamp from time")

	return &grpc.OrgDomainView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		OrgId:        domain.OrgID,
		Domain:       domain.Domain,
		Verified:     domain.Verified,
		Primary:      domain.Primary,
	}
}

func orgDomainSearchRequestToModel(request *grpc.OrgDomainSearchRequest) *org_model.OrgDomainSearchRequest {
	return &org_model.OrgDomainSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: orgDomainSearchQueriesToModel(request.Queries),
	}
}

func orgDomainSearchQueriesToModel(queries []*grpc.OrgDomainSearchQuery) []*org_model.OrgDomainSearchQuery {
	modelQueries := make([]*org_model.OrgDomainSearchQuery, len(queries))

	for i, query := range queries {
		modelQueries[i] = orgDomainSearchQueryToModel(query)
	}

	return modelQueries
}

func orgDomainSearchQueryToModel(query *grpc.OrgDomainSearchQuery) *org_model.OrgDomainSearchQuery {
	return &org_model.OrgDomainSearchQuery{
		Key:    orgDomainSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func orgDomainSearchKeyToModel(key grpc.OrgDomainSearchKey) org_model.OrgDomainSearchKey {
	switch key {
	case grpc.OrgDomainSearchKey_ORGDOMAINSEARCHKEY_DOMAIN:
		return org_model.OrgDomainSearchKeyDomain
	default:
		return org_model.OrgDomainSearchKeyUnspecified
	}
}

func orgDomainSearchResponseFromModel(resp *org_model.OrgDomainSearchResponse) *grpc.OrgDomainSearchResponse {
	return &grpc.OrgDomainSearchResponse{
		Limit:       resp.Limit,
		Offset:      resp.Offset,
		TotalResult: resp.TotalResult,
		Result:      orgDomainsFromModel(resp.Result),
	}
}
func orgDomainsFromModel(viewDomains []*org_model.OrgDomainView) []*grpc.OrgDomainView {
	domains := make([]*grpc.OrgDomainView, len(viewDomains))

	for i, domain := range viewDomains {
		domains[i] = orgDomainViewFromModel(domain)
	}

	return domains
}

func orgChangesToResponse(response *org_model.OrgChanges, offset uint64, limit uint64) (_ *grpc.Changes) {
	return &grpc.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: orgChangesToMgtAPI(response),
	}
}

func orgChangesToMgtAPI(changes *org_model.OrgChanges) (_ []*grpc.Change) {
	result := make([]*grpc.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &grpc.Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
		}
	}

	return result
}

func orgIamPolicyFromModel(policy *org_model.OrgIamPolicy) *grpc.OrgIamPolicy {
	return &grpc.OrgIamPolicy{
		OrgId:                 policy.AggregateID,
		Description:           policy.Description,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		Default:               policy.Default,
	}
}
