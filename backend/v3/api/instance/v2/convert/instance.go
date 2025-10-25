package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/cmd/build"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

/*
 * Domain Model to GRPC v2
 */

func DomainInstanceListModelToGRPCResponse(instances []*domain.Instance) []*instance.Instance {
	toReturn := make([]*instance.Instance, len(instances))

	for i, inst := range instances {
		toReturn[i] = DomainInstanceModelToGRPCResponse(inst)
	}

	return toReturn
}

func DomainInstanceModelToGRPCResponse(inst *domain.Instance) *instance.Instance {
	return &instance.Instance{
		Id:           inst.ID,
		ChangeDate:   timestamppb.New(inst.UpdatedAt),
		CreationDate: timestamppb.New(inst.CreatedAt),
		State:        instance.State_STATE_RUNNING, // TODO(IAM-Marco): Not sure what to put here
		Name:         inst.Name,
		Version:      build.Version(),
		Domains:      domainInstanceDomainListModelToGRPCResponse(inst.Domains),
	}
}

func domainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance.Domain {
	toReturn := make([]*instance.Domain, len(dms))
	for i, domain := range dms {
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		toReturn[i] = &instance.Domain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}
