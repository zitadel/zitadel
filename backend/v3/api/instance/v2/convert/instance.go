package convert

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/cmd/build"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func DomainInstanceListModelToGRPCResponse(instances []*domain.Instance) []*instance_v2.Instance {
	toReturn := make([]*instance_v2.Instance, len(instances))

	for i, inst := range instances {
		toReturn[i] = DomainInstanceModelToGRPCResponse(inst)
	}

	return toReturn
}

func DomainInstanceModelToGRPCResponse(inst *domain.Instance) *instance_v2.Instance {
	return &instance_v2.Instance{
		Id:            inst.ID,
		ChangeDate:    timestamppb.New(inst.UpdatedAt),
		CreationDate:  timestamppb.New(inst.CreatedAt),
		State:         instance_v2.State_STATE_RUNNING, // TODO(IAM-Marco): Not sure what to put here
		Name:          inst.Name,
		Version:       build.Version(),
		CustomDomains: DomainInstanceDomainListModelToGRPCResponse(inst.Domains),
	}
}
