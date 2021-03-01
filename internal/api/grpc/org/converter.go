package org

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
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
	case *org_pb.OrgQuery_Domain:
		return &org_model.OrgSearchQuery{
			Key:    org_model.OrgSearchKeyOrgDomain,
			Method: object.TextMethodToModel(q.Domain.Method),
			Value:  q.Domain.Domain,
		}, nil
	case *org_pb.OrgQuery_Name:
		//TODO: implement name in backend
		return nil, errors.ThrowUnimplemented(nil, "ADMIN-KGXnX", "name query not implemented")
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
		Details: object.ToDetailsPb(
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
		// Details: object.ToDetailsPb(//TODO: not provided
		// 	org.Sequence,//TODO: not provided
		// 	org.CreationDate,//TODO: not provided
		// 	org.ChangeDate,//TODO: not provided
		// 	org.ResourceOwner,//TODO: not provided
		// ),//TODO: not provided
	}
}

func OrgStateToPb(state org_model.OrgState) org_pb.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return org_pb.OrgState_ORG_STATE_ACTIVATE
	case org_model.OrgStateInactive:
		return org_pb.OrgState_ORG_STATE_INACTIVATE
	default:
		return org_pb.OrgState_ORG_STATE_UNSPECIFIED
	}
}
