package grpc

import (
	"encoding/json"

	"github.com/caos/logging"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func orgsFromModel(orgs []*org_model.Org) []*Org {
	orgList := make([]*Org, len(orgs))
	for i, org := range orgs {
		orgList[i] = orgFromModel(org)
	}
	return orgList
}

func orgFromModel(org *org_model.Org) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgViewFromModel(org *org_model.OrgView) *OrgView {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &OrgView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromModel(state org_model.OrgState) OrgState {
	switch state {
	case org_model.ORGSTATE_ACTIVE:
		return OrgState_ORGSTATE_ACTIVE
	case org_model.ORGSTATE_INACTIVE:
		return OrgState_ORGSTATE_INACTIVE
	default:
		return OrgState_ORGSTATE_UNSPECIFIED
	}
}

func addOrgDomainToModel(domain *AddOrgDomainRequest) *org_model.OrgDomain {
	return &org_model.OrgDomain{Domain: domain.Domain}
}

func orgDomainFromModel(domain *org_model.OrgDomain) *OrgDomain {
	creationDate, err := ptypes.TimestampProto(domain.CreationDate)
	logging.Log("GRPC-u8Ksj").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(domain.ChangeDate)
	logging.Log("GRPC-9osFS").OnError(err).Debug("unable to get timestamp from time")

	return &OrgDomain{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		OrgId:        domain.AggregateID,
		Domain:       domain.Domain,
		Verified:     domain.Verified,
		Primary:      domain.Primary,
	}
}

func orgDomainViewFromModel(domain *org_model.OrgDomainView) *OrgDomainView {
	creationDate, err := ptypes.TimestampProto(domain.CreationDate)
	logging.Log("GRPC-7sjDs").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(domain.ChangeDate)
	logging.Log("GRPC-8iSji").OnError(err).Debug("unable to get timestamp from time")

	return &OrgDomainView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		OrgId:        domain.OrgID,
		Domain:       domain.Domain,
		Verified:     domain.Verified,
		Primary:      domain.Primary,
	}
}

func orgDomainSearchRequestToModel(request *OrgDomainSearchRequest) *org_model.OrgDomainSearchRequest {
	return &org_model.OrgDomainSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: orgDomainSearchQueriesToModel(request.Queries),
	}
}

func orgDomainSearchQueriesToModel(queries []*OrgDomainSearchQuery) []*org_model.OrgDomainSearchQuery {
	modelQueries := make([]*org_model.OrgDomainSearchQuery, len(queries))

	for i, query := range queries {
		modelQueries[i] = orgDomainSearchQueryToModel(query)
	}

	return modelQueries
}

func orgDomainSearchQueryToModel(query *OrgDomainSearchQuery) *org_model.OrgDomainSearchQuery {
	return &org_model.OrgDomainSearchQuery{
		Key:    orgDomainSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func orgDomainSearchKeyToModel(key OrgDomainSearchKey) org_model.OrgDomainSearchKey {
	switch key {
	case OrgDomainSearchKey_ORGDOMAINSEARCHKEY_DOMAIN:
		return org_model.ORGDOMAINSEARCHKEY_DOMAIN
	default:
		return org_model.ORGDOMAINSEARCHKEY_UNSPECIFIED
	}
}

func orgDomainSearchResponseFromModel(resp *org_model.OrgDomainSearchResponse) *OrgDomainSearchResponse {
	return &OrgDomainSearchResponse{
		Limit:       resp.Limit,
		Offset:      resp.Offset,
		TotalResult: resp.TotalResult,
		Result:      orgDomainsFromModel(resp.Result),
	}
}
func orgDomainsFromModel(viewDomains []*org_model.OrgDomainView) []*OrgDomainView {
	domains := make([]*OrgDomainView, len(viewDomains))

	for i, domain := range viewDomains {
		domains[i] = orgDomainViewFromModel(domain)
	}

	return domains
}

func orgChangesToResponse(response *org_model.OrgChanges, offset uint64, limit uint64) (_ *Changes) {
	return &Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: orgChangesToMgtAPI(response),
	}
}

func orgChangesToMgtAPI(changes *org_model.OrgChanges) (_ []*Change) {
	result := make([]*Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
		}
	}

	return result
}

func orgIamPolicyFromModel(policy *org_model.OrgIamPolicy) *OrgIamPolicy {
	return &OrgIamPolicy{
		OrgId:                 policy.AggregateID,
		Description:           policy.Description,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		Default:               policy.Default,
	}
}
