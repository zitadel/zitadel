//go:build integration

package org_test

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

type orgAttr struct {
	ID      string
	Name    string
	Details *object.Details
}

func createOrganization(ctx context.Context, name string) orgAttr {
	orgResp := Instance.CreateOrganization(ctx, name, integration.Email())
	orgResp.Details.CreationDate = orgResp.Details.ChangeDate
	return orgAttr{
		ID:      orgResp.GetOrganizationId(),
		Name:    name,
		Details: orgResp.GetDetails(),
	}
}

func createOrganizationWithCustomOrgID(ctx context.Context, name string, orgID string) orgAttr {
	orgResp := Instance.CreateOrganizationWithCustomOrgID(ctx, name, orgID)
	orgResp.Details.CreationDate = orgResp.Details.ChangeDate
	return orgAttr{
		ID:      orgResp.GetOrganizationId(),
		Name:    name,
		Details: orgResp.GetDetails(),
	}
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
			name: "list org by default, ok",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						DefaultOrganizationQuery(),
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
						Id:            Instance.DefaultOrg.Id,
						Name:          Instance.DefaultOrg.Name,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
						State:         org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.CreationDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
					},
				},
			},
		},
		{
			name: "list org by id, ok, multiple",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{
						OrganizationIdQuery(Instance.DefaultOrg.Id),
					},
					SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				},
				func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error) {
					count := 3
					orgs := make([]orgAttr, count)
					prefix := integration.OrganizationName()
					for i := 0; i < count; i++ {
						name := prefix + strconv.Itoa(i)
						orgs[i] = createOrganization(ctx, name)
					}
					request.Queries = []*org.SearchQuery{
						OrganizationNamePrefixQuery(prefix),
					}

					slices.SortFunc(orgs, func(a, b orgAttr) int {
						return -1 * strings.Compare(a.Name, b.Name)
					})
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
						OrganizationIdQuery(Instance.DefaultOrg.Id),
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
						Name:  Instance.DefaultOrg.Name,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.CreationDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
						Id:            Instance.DefaultOrg.Id,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
					},
				},
			},
		},
		{
			name: "list org by custom id, ok",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{},
				func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error) {
					orgs := make([]orgAttr, 1)
					name := integration.OrganizationName()
					orgID := integration.ID()
					orgs[0] = createOrganizationWithCustomOrgID(ctx, name, orgID)
					request.Queries = []*org.SearchQuery{
						OrganizationIdQuery(orgID),
					}
					return orgs, nil
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
						OrganizationNameQuery(Instance.DefaultOrg.Name),
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
						Name:  Instance.DefaultOrg.Name,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.CreationDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
						Id:            Instance.DefaultOrg.Id,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
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
						OrganizationDomainQuery(Instance.DefaultOrg.PrimaryDomain),
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
						Name:  Instance.DefaultOrg.Name,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.CreationDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
						Id:            Instance.DefaultOrg.Id,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
					},
				},
			},
		},
		{
			name: "list org by domain (non primary), ok",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{},
				func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error) {
					orgs := make([]orgAttr, 1)
					orgs[0] = createOrganization(ctx, integration.OrganizationName())
					domain := integration.DomainName()
					_, err := Instance.Client.Mgmt.AddOrgDomain(integration.SetOrgID(ctx, orgs[0].ID), &management.AddOrgDomainRequest{
						Domain: domain,
					})
					if err != nil {
						return nil, err
					}
					request.Queries = []*org.SearchQuery{
						OrganizationDomainQuery(domain),
					}
					return orgs, nil
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
					},
				},
			},
		},
		{
			name: "list org by inactive state, ok",
			args: args{
				CTX,
				&org.ListOrganizationsRequest{
					Queries: []*org.SearchQuery{},
				},
				func(ctx context.Context, request *org.ListOrganizationsRequest) ([]orgAttr, error) {
					name := integration.OrganizationName()
					orgResp := createOrganization(ctx, name)
					deactivateOrgResp := Instance.DeactivateOrganization(ctx, orgResp.ID)
					request.Queries = []*org.SearchQuery{
						OrganizationIdQuery(orgResp.ID),
						OrganizationStateQuery(org.OrganizationState_ORGANIZATION_STATE_INACTIVE),
					}
					return []orgAttr{{
						ID:   orgResp.ID,
						Name: name,
						Details: &object.Details{
							ResourceOwner: deactivateOrgResp.GetDetails().GetResourceOwner(),
							Sequence:      deactivateOrgResp.GetDetails().GetSequence(),
							CreationDate:  orgResp.Details.GetCreationDate(),
							ChangeDate:    deactivateOrgResp.GetDetails().GetChangeDate(),
						},
					}}, nil
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
						State:   org.OrganizationState_ORGANIZATION_STATE_INACTIVE,
						Details: &object.Details{},
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
						OrganizationDomainQuery(Instance.DefaultOrg.PrimaryDomain),
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
						Name:  Instance.DefaultOrg.Name,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.ChangeDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
						Id:            Instance.DefaultOrg.Id,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
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
						Name:  Instance.DefaultOrg.Name,
						Details: &object.Details{
							Sequence:      Instance.DefaultOrg.Details.Sequence,
							CreationDate:  Instance.DefaultOrg.Details.ChangeDate,
							ChangeDate:    Instance.DefaultOrg.Details.ChangeDate,
							ResourceOwner: Instance.DefaultOrg.Details.ResourceOwner,
						},
						Id:            Instance.DefaultOrg.Id,
						PrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
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

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListOrganizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				// totalResult is unrelated to the tests here so gets carried over, can vary from the count of results due to permissions
				tt.want.Details.TotalResult = got.Details.TotalResult
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					for i := range tt.want.Result {
						// domain from result, as it is generated though the create
						tt.want.Result[i].PrimaryDomain = got.Result[i].PrimaryDomain
						// sequence from result, as it can be with different sequence from create
						tt.want.Result[i].Details.Sequence = got.Result[i].Details.Sequence
					}

					for i := range tt.want.Result {
						assert.Contains(ttt, got.Result, tt.want.Result[i])
					}
				}
				integration.AssertListDetails(t, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func DefaultOrganizationQuery() *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_DefaultQuery{
		DefaultQuery: &org.DefaultOrganizationQuery{},
	}}
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

func OrganizationStateQuery(state org.OrganizationState) *org.SearchQuery {
	return &org.SearchQuery{Query: &org.SearchQuery_StateQuery{
		StateQuery: &org.OrganizationStateQuery{
			State: state,
		},
	}}
}
