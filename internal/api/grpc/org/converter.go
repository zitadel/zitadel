package org

import (
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	object_pb "github.com/caos/zitadel/pkg/grpc/object"
	org_pb "github.com/caos/zitadel/pkg/grpc/org"
	"github.com/golang/protobuf/ptypes"
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
			Method: TextMethodToModel(q.Domain.Method),
			Value:  q.Domain.Domain,
		}, nil
	case *org_pb.OrgQuery_Name:
		//TODO: implement name in backend
		return nil, errors.ThrowUnimplemented(nil, "ADMIN-KGXnX", "name query not implemented")
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ADMIN-vR9nC", "List.Query.Invalid")
	}
}

func TextMethodToModel(method object_pb.TextQueryMethod) model.SearchMethod {
	switch method {
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return model.SearchMethodEquals
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return model.SearchMethodContains
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return model.SearchMethodEndsWith
	case object_pb.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return model.SearchMethodEndsWithIgnoreCase
	default:
		return -1
	}
}

func OrgsToPb(orgs []*org_model.OrgView) []*org_pb.Org {
	o := make([]*org_pb.Org, len(orgs))
	for i, org := range orgs {
		o[i] = OrgToPb(org)
	}
	return o
}

func OrgToPb(org *org_model.OrgView) *org_pb.Org {
	return &org_pb.Org{
		Id:    org.ID,
		State: OrgStateToPb(org.State),
		Name:  org.Name,
		Details: ToDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func OrgStateToPb(state org_model.OrgState) org_pb.OrgState {
	switch state {
	case org_model.OrgStateActive:
		return org_pb.OrgState_ORGSTATE_ACTIVATE
	case org_model.OrgStateInactive:
		return org_pb.OrgState_ORGSTATE_INACTIVATE
	default:
		return org_pb.OrgState_ORGSTATE_UNSPECIFIED
	}
}

func ToDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *object_pb.ObjectDetails {
	creationDatePb, err := ptypes.TimestampProto(creationDate)
	logging.Log("ADMIN-yzma4").OnError(err).Debug("unable to parse creation date")
	changeDatePb, err := ptypes.TimestampProto(changeDate)
	logging.Log("ADMIN-NTgjY").OnError(err).Debug("unable to parse change date")

	return &object_pb.ObjectDetails{
		Sequence:      sequence,
		CreationDate:  creationDatePb,
		ChangeDate:    changeDatePb,
		ResourceOwner: resourceOwner,
	}
}
