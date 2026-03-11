package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func TestTrustedDomainInstanceDomainListModelToGRPCResponse(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	tt := []struct {
		name     string
		input    []*domain.InstanceDomain
		expected []*instance_v2.TrustedDomain
	}{
		{
			name: "single domain",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test.com",
				},
			},
			expected: []*instance_v2.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test.com",
				},
			},
		},
		{
			name: "multiple domains",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test1.com",
				},
				{
					InstanceID: "instance2",
					CreatedAt:  testTime,
					Domain:     "test2.com",
				},
			},
			expected: []*instance_v2.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test1.com",
				},
				{
					InstanceId:   "instance2",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test2.com",
				},
			},
		},
		{
			name:     "empty list",
			input:    []*domain.InstanceDomain{},
			expected: []*instance_v2.TrustedDomain{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := TrustedDomainInstanceDomainListModelToGRPCResponse(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
