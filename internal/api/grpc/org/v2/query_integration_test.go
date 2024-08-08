//go:build integration

package org_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

type orgAttr struct {
	ID      string
	Name    string
	Domain  string
	Changed *timestamppb.Timestamp
	Details *object.Details
}

func TestServer_ListOrganizations(t *testing.T) {
	type args struct {
		ctx   context.Context
		count int
		req   *org.ListOrganizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *org.ListOrganizationsResponse
		wantErr bool
	}{
		{
			name: "list org by id, ok",
			args: args{
				CTX,
				1,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationIdQuery(Tester.Organisation.ID),
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*org.Organization{
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
						Name:  Tester.Organisation.Name,
						Details: &object.Details{
							Sequence:      Tester.Organisation.Sequence,
							ChangeDate:    timestamppb.New(Tester.Organisation.ChangeDate),
							ResourceOwner: Tester.Organisation.ResourceOwner,
						},
						Id:            Tester.Organisation.ID,
						PrimaryDomain: Tester.Organisation.Domain,
					},
				},
			},
		},
		{
			name: "list org by name, ok",
			args: args{
				CTX,
				1,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationNameQuery(Tester.Organisation.Name),
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*org.Organization{
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
						Name:  Tester.Organisation.Name,
						Details: &object.Details{
							Sequence:      Tester.Organisation.Sequence,
							ChangeDate:    timestamppb.New(Tester.Organisation.ChangeDate),
							ResourceOwner: Tester.Organisation.ResourceOwner,
						},
						Id:            Tester.Organisation.ID,
						PrimaryDomain: Tester.Organisation.Domain,
					},
				},
			},
		},
		{
			name: "list org by domain, ok",
			args: args{
				CTX,
				1,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery(Tester.Organisation.Domain),
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*org.Organization{
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
						Name:  Tester.Organisation.Name,
						Details: &object.Details{
							Sequence:      Tester.Organisation.Sequence,
							ChangeDate:    timestamppb.New(Tester.Organisation.ChangeDate),
							ResourceOwner: Tester.Organisation.ResourceOwner,
						},
						Id:            Tester.Organisation.ID,
						PrimaryDomain: Tester.Organisation.Domain,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			retryDuration := time.Minute
			if ctxDeadline, ok := CTX.Deadline(); ok {
				retryDuration = time.Until(ctxDeadline)
			}
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Client.ListOrganizations(tt.args.ctx, tt.args.req)
				assertErr := assert.NoError
				if tt.wantErr {
					assertErr = assert.Error
				}
				assertErr(ttt, listErr)
				if listErr != nil {
					return
				}
				// always only give back dependency infos which are required for the response
				assert.Len(ttt, tt.want.Result, int(tt.want.Details.TotalResult))
				// always first check length, otherwise its failed anyway
				assert.Len(ttt, got.Result, len(tt.want.Result))

				for i := range tt.want.Result {
					tt.want.Result[i].PrimaryDomain = got.Result[i].PrimaryDomain
				}

				for i := range tt.want.Result {
					assert.Contains(ttt, got.Result, tt.want.Result[i])
				}
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, time.Millisecond*100, "timeout waiting for expected user result")
		})
	}
}

func OrganizationIdQuery(resourceowner string) *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_IdQuery{
		IdQuery: &org.OrganizationIDQuery{
			Id: resourceowner,
		},
	}}
}

func OrganizationNameQuery(name string) *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_NameQuery{
		NameQuery: &org.OrganizationNameQuery{
			Name: name,
		},
	}}
}

func OrganizationDomainQuery(domain string) *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_DomainQuery{
		DomainQuery: &org.OrganizationDomainQuery{
			Domain: domain,
		},
	}}
}
