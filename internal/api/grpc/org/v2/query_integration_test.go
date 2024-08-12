//go:build integration

package org_test

import (
	"context"
	"fmt"
	"strconv"
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
	Details *object.Details
}

func TestServer_ListOrganizations(t *testing.T) {
	type args struct {
		ctx context.Context
		req *org.ListOrganizationsRequest
		dep func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error)
	}
	tests := []struct {
		name    string
		args    args
		want    *org.ListOrganizationsResponse
		wantErr bool
	}{
		{
			name: "list org by id, ok, multiple",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationIdQuery(Tester.Organisation.ID),
					},
				},
				func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error) {
					count := 3
					orgs := make([]orgAttr, count)
					prefix := fmt.Sprintf("ListOrgs%d", time.Now().UnixNano())
					for i := 0; i < count; i++ {
						name := prefix + strconv.Itoa(i)
						orgResp := Tester.CreateOrganization(ctx, name, fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()))
						orgs[i] = orgAttr{
							ID:      orgResp.GetOrganizationId(),
							Name:    name,
							Details: orgResp.GetDetails(),
						}
					}
					request.Queries = []*org.SearchQuery{
						OrganizationNamePrefixQuery(prefix),
					}
					return orgs, nil
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*org.Organization{
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
					},
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
					},
					{
						State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
					},
				},
			},
		},
		{
			name: "list org by id, ok",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationIdQuery(Tester.Organisation.ID),
					},
				},
				nil,
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
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationNameQuery(Tester.Organisation.Name),
					},
				},
				nil,
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
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery(Tester.Organisation.Domain),
					},
				},
				nil,
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
			name: "list org by domain, ok, sorted",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery(Tester.Organisation.Domain),
					},
					SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				},
				nil,
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 1,
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
			name: "list org, no result",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery("notexisting"),
					},
				},
				nil,
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*org.Organization{},
			},
		},
		{
			name: "list org, no login",
			args: args{
				context.Background(),
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery("nopermission"),
					},
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "list org, no permission",
			args: args{
				UserCTX,
				&org.ListOrganizationsRequest{},
				nil,
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 1,
				Result:        []*org.Organization{},
			},
		},
		{
			name: "list org, no permission org owner",
			args: args{
				OwnerCTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationDomainQuery("nopermission"),
					},
				},
				nil,
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 1,
				Result:        []*org.Organization{},
			},
		},
		{
			name: "list org, org owner",
			args: args{
				OwnerCTX,
				&org.ListOrganizationsRequest{},
				nil,
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 1,
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
			if tt.args.dep != nil {
				orgs, err := tt.args.dep(tt.args.ctx, tt.args.req)
				require.NoError(t, err)
				if len(orgs) > 0 {
					for i, org := range orgs {
						tt.want.Result[i].Name = org.Name
						tt.want.Result[i].Id = org.ID
						tt.want.Result[i].Details = org.Details
					}
				}
			}

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
					// domain from result, as it is generated though the create
					tt.want.Result[i].PrimaryDomain = got.Result[i].PrimaryDomain
					// sequence from result, as it can be with different sequence from create
					tt.want.Result[i].Details.Sequence = got.Result[i].Details.Sequence
				}

				fmt.Println(tt.want.Result)
				fmt.Println(got.Result)
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

func OrganizationNamePrefixQuery(name string) *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_NameQuery{
		NameQuery: &org.OrganizationNameQuery{
			Name:   name,
			Method: object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
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
