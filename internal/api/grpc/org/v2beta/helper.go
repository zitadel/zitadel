package org

import (
	"time"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"

	v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"

	org_pb "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrganizationViewToPb(org *query.Org) *org_pb.Organization {
	return &org_pb.Organization{
		Id:            org.ID,
		State:         OrgStateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.Domain,
		Details: ToViewDetailsPb(
			org.Sequence,
			org.CreationDate,
			org.ChangeDate,
			org.ResourceOwner,
		),
	}
}

func OrgStateToPb(state domain.OrgState) org_pb.OrganizationState {
	switch state {
	case domain.OrgStateActive:
		return org_pb.OrganizationState_ORGANIZATION_STATE_ACTIVE
	case domain.OrgStateInactive:
		return org_pb.OrganizationState_ORGANIZATION_STATE_INACTIVE
	default:
		return org_pb.OrganizationState_ORGANIZATION_STATE_UNSPECIFIED
	}
}

func ToViewDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *v2beta.Details {
	details := &v2beta.Details{
		Sequence:      sequence,
		ResourceOwner: resourceOwner,
	}
	if !creationDate.IsZero() {
		details.CreationDate = timestamppb.New(creationDate)
	}
	if !changeDate.IsZero() {
		details.ChangeDate = timestamppb.New(changeDate)
	}
	return details
}

func createdOrganizationToPb(createdOrg *command.CreatedOrg) (_ *org.CreateOrganizationResponse, err error) {
	admins := make([]*org.CreateOrganizationResponse_CreatedAdmin, len(createdOrg.CreatedAdmins))
	for i, admin := range createdOrg.CreatedAdmins {
		admins[i] = &org.CreateOrganizationResponse_CreatedAdmin{
			UserId:    admin.ID,
			EmailCode: admin.EmailCode,
			PhoneCode: admin.PhoneCode,
		}
	}
	return &org.CreateOrganizationResponse{
		Details:        object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganizationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}, nil
}
