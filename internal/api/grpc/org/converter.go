package org

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	org_pb "github.com/caos/zitadel/pkg/grpc/org"
)

func OrgQueriesToModel(queries []*org_pb.OrgQuery) (_ []*org_model.OrgSearchQuery, err error) {
	q := make([]*org_model.OrgSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = OrgQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func OrgQueryToModel(query *org_pb.OrgQuery) (*org_model.OrgSearchQuery, error) {
	switch q := query.Query.(type) {
	case *org_pb.OrgQuery_DomainQuery:
		return &org_model.OrgSearchQuery{
			Key:    org_model.OrgSearchKeyOrgDomain,
			Method: object.TextMethodToModel(q.DomainQuery.Method),
			Value:  q.DomainQuery.Domain,
		}, nil
	case *org_pb.OrgQuery_NameQuery:
		return &org_model.OrgSearchQuery{
			Key:    org_model.OrgSearchKeyOrgName,
			Method: object.TextMethodToModel(q.NameQuery.Method),
			Value:  q.NameQuery.Name,
		}, nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-vR9nC", "List.Query.Invalid")
	}
}

func OrgViewsToPb(orgs []*org_model.OrgView) []*org_pb.Org {
	o := make([]*org_pb.Org, len(orgs))
	for i, org := range orgs {
		o[i] = OrgViewToPb(org)
	}
	return o
}

func OrgViewToPb(org *org_model.OrgView) *org_pb.Org {
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

func OrgStateToPb(state org_model.OrgState) org_pb.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return org_pb.OrgState_ORG_STATE_ACTIVE
	case org_model.OrgStateInactive:
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
