package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func DomainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance_v2.CustomDomain {
	toReturn := make([]*instance_v2.CustomDomain, len(dms))
	for i, domain := range dms {
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		toReturn[i] = &instance_v2.CustomDomain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}

func TrustedDomainInstanceDomainListModelToGRPCResponse(dms []*domain.InstanceDomain) []*instance_v2.TrustedDomain {
	toReturn := make([]*instance_v2.TrustedDomain, len(dms))
	for i, domain := range dms {
		toReturn[i] = &instance_v2.TrustedDomain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
		}
	}

	return toReturn
}
