package admin

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func orgCreateRequestToDomain(org *admin.CreateOrgRequest) *domain.Org {
	o := &domain.Org{
		Domains: []*domain.OrgDomain{},
		Name:    org.Name,
	}
	if org.Domain != "" {
		o.Domains = append(o.Domains, &domain.OrgDomain{Domain: org.Domain})
	}

	return o
}

func setUpOrgResponseFromModel(setUp *admin_model.SetupOrg) *admin.OrgSetUpResponse {
	return &admin.OrgSetUpResponse{
		Org:  orgFromModel(setUp.Org),
		User: userFromModel(setUp.User),
	}
}

func orgSearchResponseFromModel(request *org_model.OrgSearchResult) *admin.OrgSearchResponse {
	timestamp, err := ptypes.TimestampProto(request.Timestamp)
	logging.Log("GRPC-shu7s").OnError(err).Debug("unable to get timestamp from time")
	return &admin.OrgSearchResponse{
		Result:            orgViewsFromModel(request.Result),
		Limit:             request.Limit,
		Offset:            request.Offset,
		TotalResult:       request.TotalResult,
		ProcessedSequence: request.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func orgViewsFromModel(orgs []*org_model.OrgView) []*admin.Org {
	result := make([]*admin.Org, len(orgs))
	for i, org := range orgs {
		result[i] = orgViewFromModel(org)
	}

	return result
}

func orgFromModel(org *org_model.Org) *admin.Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &admin.Org{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgViewFromModel(org *org_model.OrgView) *admin.Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &admin.Org{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromModel(state org_model.OrgState) admin.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return admin.OrgState_ORGSTATE_ACTIVE
	case org_model.OrgStateInactive:
		return admin.OrgState_ORGSTATE_INACTIVE
	default:
		return admin.OrgState_ORGSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) admin.Gender {
	switch gender {
	case usr_model.GenderFemale:
		return admin.Gender_GENDER_FEMALE
	case usr_model.GenderMale:
		return admin.Gender_GENDER_MALE
	case usr_model.GenderDiverse:
		return admin.Gender_GENDER_DIVERSE
	default:
		return admin.Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender admin.Gender) usr_model.Gender {
	switch gender {
	case admin.Gender_GENDER_FEMALE:
		return usr_model.GenderFemale
	case admin.Gender_GENDER_MALE:
		return usr_model.GenderMale
	case admin.Gender_GENDER_DIVERSE:
		return usr_model.GenderDiverse
	default:
		return usr_model.GenderUnspecified
	}
}

func userStateFromModel(state usr_model.UserState) admin.UserState {
	switch state {
	case usr_model.UserStateActive:
		return admin.UserState_USERSTATE_ACTIVE
	case usr_model.UserStateInactive:
		return admin.UserState_USERSTATE_INACTIVE
	case usr_model.UserStateLocked:
		return admin.UserState_USERSTATE_LOCKED
	default:
		return admin.UserState_USERSTATE_UNSPECIFIED
	}
}

func orgSearchRequestToModel(req *admin.OrgSearchRequest) *org_model.OrgSearchRequest {
	return &org_model.OrgSearchRequest{
		Limit:         req.Limit,
		Asc:           req.Asc,
		Offset:        req.Offset,
		Queries:       orgQueriesToModel(req.Queries),
		SortingColumn: orgQueryKeyToModel(req.SortingColumn),
	}
}

func orgQueriesToModel(queries []*admin.OrgSearchQuery) []*org_model.OrgSearchQuery {
	modelQueries := make([]*org_model.OrgSearchQuery, len(queries))

	for i, query := range queries {
		modelQueries[i] = orgQueryToModel(query)
	}

	return modelQueries
}

func orgQueryToModel(query *admin.OrgSearchQuery) *org_model.OrgSearchQuery {
	return &org_model.OrgSearchQuery{
		Key:    orgQueryKeyToModel(query.Key),
		Value:  query.Value,
		Method: orgQueryMethodToModel(query.Method),
	}
}

func orgQueryKeyToModel(key admin.OrgSearchKey) org_model.OrgSearchKey {
	switch key {
	case admin.OrgSearchKey_ORGSEARCHKEY_DOMAIN:
		return org_model.OrgSearchKeyOrgDomain
	case admin.OrgSearchKey_ORGSEARCHKEY_NAME:
		return org_model.OrgSearchKeyOrgName
	case admin.OrgSearchKey_ORGSEARCHKEY_STATE:
		return org_model.OrgSearchKeyState
	default:
		return org_model.OrgSearchKeyUnspecified
	}
}

func orgQueryMethodToModel(method admin.OrgSearchMethod) model.SearchMethod {
	switch method {
	case admin.OrgSearchMethod_ORGSEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case admin.OrgSearchMethod_ORGSEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case admin.OrgSearchMethod_ORGSEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	default:
		return 0
	}
}

func orgIAMPolicyFromDomain(policy *domain.OrgIAMPolicy) *admin.OrgIamPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-ush36").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-Ps9fW").OnError(err).Debug("unable to get timestamp from time")

	return &admin.OrgIamPolicy{
		OrgId:                 policy.AggregateID,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
	}
}

func orgIAMPolicyViewFromModel(policy *iam_model.OrgIAMPolicyView) *admin.OrgIamPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-ush36").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-Ps9fW").OnError(err).Debug("unable to get timestamp from time")

	return &admin.OrgIamPolicyView{
		OrgId:                 policy.AggregateID,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
	}
}

func orgIAMPolicyRequestToDomain(policy *admin.OrgIamPolicyRequest) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.OrgId,
		},
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}
