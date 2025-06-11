package serviceping

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	analytics "github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta"
)

type ServicePing struct {
	TelemetryServiceClient analytics.TelemetryServiceClient
}

func NewClient(ctx context.Context, config Config) {
	connection, err := grpc.NewClient(config.Endpoint)
	_ = err
	analytics.NewTelemetryServiceClient(connection)

}

func GenerateSystemID() (string, error) {
	randBytes := make([]byte, 64)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randBytes), nil
}

func instanceInformationToPb(instances *query.Instances) []*analytics.InstanceInformation {
	instanceInformation := make([]*analytics.InstanceInformation, len(instances.Instances))
	for i, instance := range instances.Instances {
		domains := instanceDomainToPb(instance)
		instanceInformation[i] = &analytics.InstanceInformation{
			Id:        instance.ID,
			Domains:   domains,
			CreatedAt: timestamppb.New(instance.CreationDate),
		}
	}
	return instanceInformation
}

func instanceDomainToPb(instance *query.Instance) []string {
	domains := make([]string, len(instance.Domains))
	for i, domain := range instance.Domains {
		domains[i] = domain.Domain
	}
	return domains
}

func resourceCountsToPb(counts []query.ResourceCount) []*analytics.ResourceCount {
	resourceCounts := make([]*analytics.ResourceCount, len(counts))
	for i, count := range counts {
		resourceCounts[i] = &analytics.ResourceCount{
			InstanceId:   count.InstanceID,
			ParentType:   countParentTypeToPb(count.ParentType),
			ParentId:     count.ParentID,
			ResourceName: count.Resource,
			TableName:    count.TableName,
			UpdatedAt:    timestamppb.New(count.UpdatedAt),
			Amount:       uint32(count.Amount),
		}
	}
	return resourceCounts
}

func countParentTypeToPb(parentType domain.CountParentType) analytics.CountParentType {
	switch parentType {
	case domain.CountParentTypeInstance:
		return analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE
	case domain.CountParentTypeOrganization:
		return analytics.CountParentType_COUNT_PARENT_TYPE_ORGANIZATION
	default:
		return analytics.CountParentType_COUNT_PARENT_TYPE_UNSPECIFIED
	}
}
