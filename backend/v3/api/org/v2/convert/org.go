package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	grpc_object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func DomainOrganizationListModelToGRPCResponse(orgs []*domain.Organization) []*v2_org.Organization {
	toReturn := make([]*v2_org.Organization, len(orgs))

	for i, org := range orgs {
		toReturn[i] = domainOrganizationModelToGRPCResponse(org)
	}

	return toReturn
}

func domainOrganizationModelToGRPCResponse(org *domain.Organization) *v2_org.Organization {
	return &v2_org.Organization{
		Id: org.ID,
		Details: &grpc_object.Details{
			ChangeDate:   timestamppb.New(org.UpdatedAt),
			CreationDate: timestamppb.New(org.CreatedAt),
		},
		State:         v2_org.OrganizationState(org.State),
		Name:          org.Name,
		PrimaryDomain: org.PrimaryDomain(),
	}
}
