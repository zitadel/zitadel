package org

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/query"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	org_pb "github.com/caos/zitadel/pkg/grpc/org"
)

func OrgQueriesToModel(queries []*org_pb.OrgQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToModel(apiQuery *org_pb.OrgQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *org_pb.OrgQuery_DomainQuery:
		return query.NewOrgDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.Method), q.DomainQuery.Domain)
	case *org_pb.OrgQuery_NameQuery:
		return query.NewOrgNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func OrgQueriesToUserGrantModel(queries []*org_pb.OrgQuery) (_ []*grant_model.UserGrantSearchQuery, err error) {
	q := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToUserGrantQueryModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToUserGrantQueryModel(query *org_pb.OrgQuery) (*grant_model.UserGrantSearchQuery, error) {
	switch q := query.Query.(type) {
	case *org_pb.OrgQuery_DomainQuery:
		return &grant_model.UserGrantSearchQuery{
			Key:    grant_model.UserGrantSearchKeyOrgDomain,
			Method: object.TextMethodToModel(q.DomainQuery.Method),
			Value:  q.DomainQuery.Domain,
		}, nil
	case *org_pb.OrgQuery_NameQuery:
		return &grant_model.UserGrantSearchQuery{
			Key:    grant_model.UserGrantSearchKeyOrgName,
			Method: object.TextMethodToModel(q.NameQuery.Method),
			Value:  q.NameQuery.Name,
		}, nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-ADvsd", "List.Query.Invalid")
	}
}

func OrgViewsToPb(orgs []*query.Org) []*org_pb.Org {
	o := make([]*org_pb.Org, len(orgs))
	for i, org := range orgs {
		o[i] = OrgViewToPb(org)
	}
	return o
}

func OrgViewToPb(org *query.Org) *org_pb.Org {
	return &org_pb.Org{
		Id:    org.ID,
		State: OrgStateToPb(org.State),
		Name:  org.Name,
		Details: object.ToViewDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func OrgsToPb(orgs []*grant_model.Org) []*org_pb.Org {
	o := make([]*org_pb.Org, len(orgs))
	for i, org := range orgs {
		o[i] = OrgToPb(org)
	}
	return o
}

func OrgToPb(org *grant_model.Org) *org_pb.Org {
	return &org_pb.Org{
		Id:   org.OrgID,
		Name: org.OrgName,
		// State: OrgStateToPb(org.State), //TODO: not provided
		// Details: object.ChangeToDetailsPb(//TODO: not provided
		// 	org.Sequence,//TODO: not provided
		// 	org.CreationDate,//TODO: not provided
		// 	org.EventDate,//TODO: not provided
		// 	org.ResourceOwner,//TODO: not provided
		// ),//TODO: not provided
	}
}

func OrgStateToPb(state domain.OrgState) org_pb.OrgState {
	switch state {
	case domain.OrgStateActive:
		return org_pb.OrgState_ORG_STATE_ACTIVE
	case domain.OrgStateInactive:
		return org_pb.OrgState_ORG_STATE_INACTIVE
	default:
		return org_pb.OrgState_ORG_STATE_UNSPECIFIED
	}
}

func DomainQueriesToModel(queries []*org_pb.DomainSearchQuery) (_ []*org_model.OrgDomainSearchQuery, err error) {
	q := make([]*org_model.OrgDomainSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = DomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func DomainQueryToModel(query *org_pb.DomainSearchQuery) (*org_model.OrgDomainSearchQuery, error) {
	switch q := query.Query.(type) {
	case *org_pb.DomainSearchQuery_DomainNameQuery:
		return DomainNameQueryToModel(q.DomainNameQuery)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func DomainNameQueryToModel(query *org_pb.DomainNameQuery) (*org_model.OrgDomainSearchQuery, error) {
	return &org_model.OrgDomainSearchQuery{
		Key:    org_model.OrgDomainSearchKeyDomain,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Name,
	}, nil
}

func DomainsToPb(domains []*org_model.OrgDomainView) []*org_pb.Domain {
	d := make([]*org_pb.Domain, len(domains))
	for i, domain := range domains {
		d[i] = DomainToPb(domain)
	}
	return d
}

func DomainToPb(domain *org_model.OrgDomainView) *org_pb.Domain {
	return &org_pb.Domain{
		OrgId:          domain.OrgID,
		DomainName:     domain.Domain,
		IsVerified:     domain.Verified,
		IsPrimary:      domain.Primary,
		ValidationType: DomainValidationTypeFromModel(domain.ValidationType),
		Details: object.ToViewDetailsPb(
			0,
			domain.CreationDate,
			domain.ChangeDate,
			"",
		),
	}
}

func DomainValidationTypeToDomain(validationType org_pb.DomainValidationType) domain.OrgDomainValidationType {
	switch validationType {
	case org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP:
		return domain.OrgDomainValidationTypeHTTP
	case org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS:
		return domain.OrgDomainValidationTypeDNS
	default:
		return domain.OrgDomainValidationTypeUnspecified
	}
}

func DomainValidationTypeFromModel(validationType org_model.OrgDomainValidationType) org_pb.DomainValidationType {
	switch validationType {
	case org_model.OrgDomainValidationTypeDNS:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS
	case org_model.OrgDomainValidationTypeHTTP:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP
	default:
		return org_pb.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED
	}
}
