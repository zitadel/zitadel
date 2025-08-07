package serviceping

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	analytics "github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta"
)

const (
	pathBaseInformation = "/instances"
	pathResourceCounts  = "/resource_counts"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func (c Client) ReportBaseInformation(ctx context.Context, in *analytics.ReportBaseInformationRequest, opts ...grpc.CallOption) (*analytics.ReportBaseInformationResponse, error) {
	reportResponse := new(analytics.ReportBaseInformationResponse)
	err := c.callTelemetryService(ctx, pathBaseInformation, in, reportResponse)
	if err != nil {
		return nil, err
	}
	return reportResponse, nil
}

func (c Client) ReportResourceCounts(ctx context.Context, in *analytics.ReportResourceCountsRequest, opts ...grpc.CallOption) (*analytics.ReportResourceCountsResponse, error) {
	reportResponse := new(analytics.ReportResourceCountsResponse)
	err := c.callTelemetryService(ctx, pathResourceCounts, in, reportResponse)
	if err != nil {
		return nil, err
	}
	return reportResponse, nil
}

func (c Client) callTelemetryService(ctx context.Context, path string, in proto.Message, out proto.Message) error {
	requestBody, err := protojson.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+path, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &TelemetryError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(body, out)
}

func NewClient(config *Config) Client {
	return Client{
		httpClient: http.DefaultClient,
		endpoint:   config.Endpoint,
	}
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

type TelemetryError struct {
	StatusCode int
	Body       []byte
}

func (e *TelemetryError) Error() string {
	return fmt.Sprintf("telemetry error %d: %s", e.StatusCode, e.Body)
}
