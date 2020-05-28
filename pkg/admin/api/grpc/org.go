package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/model"

	org_model "github.com/caos/zitadel/internal/org/model"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *OrgID) (_ *Org, err error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchOrgs(ctx context.Context, request *OrgSearchRequest) (_ *OrgSearchResponse, err error) {
	result, err := s.org.SearchOrgs(ctx, orgSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return &OrgSearchResponse{
		Result:      orgViewsFromModel(result.Result),
		Limit:       request.Limit,
		Offset:      request.Offset,
		TotalResult: result.TotalResult,
	}, nil
}

func (s *Server) IsOrgUnique(ctx context.Context, request *UniqueOrgRequest) (org *UniqueOrgResponse, err error) {
	isUnique, err := s.org.IsOrgUnique(ctx, request.Name, request.Domain)

	return &UniqueOrgResponse{IsUnique: isUnique}, err
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *OrgSetUpRequest) (_ *OrgSetUpResponse, err error) {
	setUp, err := s.org.SetUpOrg(ctx, setUpRequestToModel(orgSetUp))
	if err != nil {
		return nil, err
	}
	return setUpOrgResponseFromModel(setUp), err
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
		return org_model.ORGSEARCHKEY_ORG_DOMAIN
	case OrgSearchKey_ORGSEARCHKEY_ORG_NAME:
		return org_model.ORGSEARCHKEY_ORG_NAME
	case OrgSearchKey_ORGSEARCHKEY_STATE:
		return org_model.ORGSEARCHKEY_STATE
	default:
		return org_model.ORGSEARCHKEY_UNSPECIFIED
	}
}

func orgQueryMethodToModel(method OrgSearchMethod) model.SearchMethod {
	switch method {
	case OrgSearchMethod_ORGSEARCHMETHOD_CONTAINS:
		return model.SEARCHMETHOD_CONTAINS
	case OrgSearchMethod_ORGSEARCHMETHOD_EQUALS:
		return model.SEARCHMETHOD_EQUALS
	case OrgSearchMethod_ORGSEARCHMETHOD_STARTS_WITH:
		return model.SEARCHMETHOD_STARTS_WITH
	default:
		return 0
	}
}
