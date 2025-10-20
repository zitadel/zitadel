package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestDomainOrganizationListModelToGRPCResponse(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)

	tt := []struct {
		name string
		orgs []*domain.Organization
		want []*org.Organization
	}{
		{
			name: "empty result",
			orgs: nil,
			want: []*org.Organization{},
		},
		{
			name: "multiple organizations",
			orgs: []*domain.Organization{
				{
					ID:        "org-1",
					Name:      "org 1",
					State:     domain.OrgStateActive,
					CreatedAt: yesterday,
					UpdatedAt: now,
					Domains: []*domain.OrganizationDomain{
						{Domain: "wrong selected domain"},
						{Domain: "domain.example.com", IsPrimary: true},
					},
				},
				{
					ID:        "org-2",
					Name:      "org 2",
					State:     domain.OrgStateInactive,
					CreatedAt: yesterday,
					UpdatedAt: now,
					Domains: []*domain.OrganizationDomain{
						{Domain: "wrong selected domain 2"},
						{Domain: "domain2.example.com", IsPrimary: true},
					},
				},
			},
			want: []*org.Organization{
				{
					Id:    "org-1",
					Name:  "org 1",
					State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
					Details: &object.Details{
						ChangeDate:   timestamppb.New(now),
						CreationDate: timestamppb.New(yesterday),
					},
					PrimaryDomain: "domain.example.com",
				},
				{
					Id:    "org-2",
					Name:  "org 2",
					State: org.OrganizationState_ORGANIZATION_STATE_INACTIVE,
					Details: &object.Details{
						ChangeDate:   timestamppb.New(now),
						CreationDate: timestamppb.New(yesterday),
					},
					PrimaryDomain: "domain2.example.com",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DomainOrganizationListModelToGRPCResponse(tc.orgs)
			assert.Equal(t, tc.want, got)
		})
	}
}
