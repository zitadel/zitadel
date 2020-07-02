package grpc

import (
	"github.com/caos/logging"
	admin_model "github.com/caos/zitadel/internal/admin/model"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func setUpRequestToModel(setUp *OrgSetUpRequest) *admin_model.SetupOrg {
	return &admin_model.SetupOrg{
		Org:  orgCreateRequestToModel(setUp.Org),
		User: userCreateRequestToModel(setUp.User),
	}
}

func orgCreateRequestToModel(org *CreateOrgRequest) *org_model.Org {
	return &org_model.Org{
		Domains: []*org_model.OrgDomain{&org_model.OrgDomain{Domain: org.Domain}},
		Name:    org.Name,
	}
}

func userCreateRequestToModel(user *CreateUserRequest) *usr_model.User {
	preferredLanguage, err := language.Parse(user.PreferredLanguage)
	logging.Log("GRPC-30hwz").OnError(err).Debug("unable to parse language")
	result := &usr_model.User{
		Profile: &usr_model.Profile{
			UserName:          user.UserName,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			NickName:          user.NickName,
			PreferredLanguage: preferredLanguage,
			Gender:            genderToModel(user.Gender),
		},
		Password: &usr_model.Password{
			SecretString: user.Password,
		},
		Email: &usr_model.Email{
			EmailAddress:    user.Email,
			IsEmailVerified: user.IsEmailVerified,
		},
		Address: &usr_model.Address{
			Country:       user.Country,
			Locality:      user.Locality,
			PostalCode:    user.PostalCode,
			Region:        user.Region,
			StreetAddress: user.StreetAddress,
		},
	}
	if user.Phone != "" {
		result.Phone = &usr_model.Phone{PhoneNumber: user.Phone, IsPhoneVerified: user.IsPhoneVerified}
	}
	return result
}

func setUpOrgResponseFromModel(setUp *admin_model.SetupOrg) *OrgSetUpResponse {
	return &OrgSetUpResponse{
		Org:  orgFromModel(setUp.Org),
		User: userFromModel(setUp.User),
	}
}

func orgViewsFromModel(orgs []*org_model.OrgView) []*Org {
	result := make([]*Org, len(orgs))
	for i, org := range orgs {
		result[i] = orgViewFromModel(org)
	}

	return result
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

func orgViewFromModel(org *org_model.OrgView) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func userFromModel(user *usr_model.User) *User {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	converted := &User{
		Id:                user.AggregateID,
		State:             userStateFromModel(user.State),
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          user.Sequence,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		Gender:            genderFromModel(user.Gender),
	}
	if user.Email != nil {
		converted.Email = user.EmailAddress
		converted.IsEmailVerified = user.IsEmailVerified
	}
	if user.Phone != nil {
		converted.Phone = user.PhoneNumber
		converted.IsPhoneVerified = user.IsPhoneVerified
	}
	if user.Address != nil {
		converted.Country = user.Country
		converted.Locality = user.Locality
		converted.PostalCode = user.PostalCode
		converted.Region = user.Region
		converted.StreetAddress = user.StreetAddress
	}
	return converted
}

func orgStateFromModel(state org_model.OrgState) OrgState {
	switch state {
	case org_model.OrgStateActive:
		return OrgState_ORGSTATE_ACTIVE
	case org_model.OrgStateInactive:
		return OrgState_ORGSTATE_INACTIVE
	default:
		return OrgState_ORGSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) Gender {
	switch gender {
	case usr_model.GenderFemale:
		return Gender_GENDER_FEMALE
	case usr_model.GenderMale:
		return Gender_GENDER_MALE
	case usr_model.GenderDiverse:
		return Gender_GENDER_DIVERSE
	default:
		return Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender Gender) usr_model.Gender {
	switch gender {
	case Gender_GENDER_FEMALE:
		return usr_model.GenderFemale
	case Gender_GENDER_MALE:
		return usr_model.GenderMale
	case Gender_GENDER_DIVERSE:
		return usr_model.GenderDiverse
	default:
		return usr_model.GenderUnspecified
	}
}

func userStateFromModel(state usr_model.UserState) UserState {
	switch state {
	case usr_model.UserStateActive:
		return UserState_USERSTATE_ACTIVE
	case usr_model.UserStateInactive:
		return UserState_USERSTATE_INACTIVE
	case usr_model.UserStateLocked:
		return UserState_USERSTATE_LOCKED
	default:
		return UserState_USERSTATE_UNSPECIFIED
	}
}

func orgSearchRequestToModel(req *OrgSearchRequest) *org_model.OrgSearchRequest {
	return &org_model.OrgSearchRequest{
		Limit:         req.Limit,
		Asc:           req.Asc,
		Offset:        req.Offset,
		Queries:       orgQueriesToModel(req.Queries),
		SortingColumn: orgQueryKeyToModel(req.SortingColumn),
	}
}

func orgQueriesToModel(queries []*OrgSearchQuery) []*org_model.OrgSearchQuery {
	modelQueries := make([]*org_model.OrgSearchQuery, len(queries))

	for i, query := range queries {
		modelQueries[i] = orgQueryToModel(query)
	}

	return modelQueries
}

func orgQueryToModel(query *OrgSearchQuery) *org_model.OrgSearchQuery {
	return &org_model.OrgSearchQuery{
		Key:    orgQueryKeyToModel(query.Key),
		Value:  query.Value,
		Method: orgQueryMethodToModel(query.Method),
	}
}

func orgQueryKeyToModel(key OrgSearchKey) org_model.OrgSearchKey {
	switch key {
	case OrgSearchKey_ORGSEARCHKEY_DOMAIN:
		return org_model.OrgSearchKeyOrgDomain
	case OrgSearchKey_ORGSEARCHKEY_ORG_NAME:
		return org_model.OrgSearchKeyOrgName
	case OrgSearchKey_ORGSEARCHKEY_STATE:
		return org_model.OrgSearchKeyState
	default:
		return org_model.OrgSearchKeyUnspecified
	}
}

func orgQueryMethodToModel(method OrgSearchMethod) model.SearchMethod {
	switch method {
	case OrgSearchMethod_ORGSEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case OrgSearchMethod_ORGSEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case OrgSearchMethod_ORGSEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	default:
		return 0
	}
}

func orgIamPolicyFromModel(policy *org_model.OrgIamPolicy) *OrgIamPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-ush36").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-Ps9fW").OnError(err).Debug("unable to get timestamp from time")

	return &OrgIamPolicy{
		OrgId:                 policy.AggregateID,
		Description:           policy.Description,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		Default:               policy.Default,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
	}
}

func orgIamPolicyRequestToModel(policy *OrgIamPolicyRequest) *org_model.OrgIamPolicy {
	return &org_model.OrgIamPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: policy.OrgId,
		},
		Description:           policy.Description,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}
