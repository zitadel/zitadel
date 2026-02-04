//go:build integration

package org_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	object_v2 "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

// TODO(IAM-Marco): When permission checks will be implemented, this test needs to be updated to
// add the feature flag switch:
//
//	relTableState := integration.RelationalTablesEnableMatrix()
//
// See TestServer_ListOrganizations in org/v2beta/integration_test
// See https://github.com/zitadel/zitadel/issues/10219
func TestServer_ListOrganizations(t *testing.T) {
	ListOrgIinstance := integration.NewInstance(CTX)
	listOrgIAmOwnerCtx := ListOrgIinstance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	listOrgClient := ListOrgIinstance.Client.OrgV2

	noOfOrgs := 3
	orgs, orgsName, orgsDomain := createOrgs(listOrgIAmOwnerCtx, t, listOrgClient, noOfOrgs)

	// deactivate org[1]
	_, err := listOrgClient.DeactivateOrganization(listOrgIAmOwnerCtx, &org.DeactivateOrganizationRequest{
		OrganizationId: orgs[1].OrganizationId,
	})
	require.NoError(t, err)

	tests := []struct {
		name  string
		ctx   context.Context
		query []*org.SearchQuery
		want  *org.ListOrganizationsResponse
		err   error
	}{
		{
			name: "list organizations, without required permissions",
			ctx:  ListOrgIinstance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			want: &org.ListOrganizationsResponse{
				Details: &object_v2.ListDetails{
					TotalResult: 4,
				},
			},
		},
		{
			name: "list organizations happy path, no filter",
			ctx:  listOrgIAmOwnerCtx,
			want: &org.ListOrganizationsResponse{
				Details: &object_v2.ListDetails{
					TotalResult: 4,
				},
				Result: []*org.Organization{
					{
						Id:   ListOrgIinstance.DefaultOrg.Id,
						Name: ListOrgIinstance.DefaultOrg.Name,
						Details: &object_v2.Details{
							Sequence:      ListOrgIinstance.DefaultOrg.GetDetails().GetSequence(),
							ChangeDate:    ListOrgIinstance.DefaultOrg.GetDetails().GetChangeDate(),
							ResourceOwner: ListOrgIinstance.DefaultOrg.GetDetails().GetResourceOwner(),
							CreationDate:  ListOrgIinstance.DefaultOrg.GetDetails().GetCreationDate(),
						},
					},
					{
						Id:      orgs[0].OrganizationId,
						Name:    orgsName[0],
						Details: orgs[0].GetDetails(),
					},
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
					{
						Id:      orgs[2].OrganizationId,
						Name:    orgsName[2],
						Details: orgs[2].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations by id happy path",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_IdQuery{
						IdQuery: &org.OrganizationIDQuery{
							Id: orgs[1].OrganizationId,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object_v2.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations by state active",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_StateQuery{
						StateQuery: &org.OrganizationStateQuery{
							State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object_v2.ListDetails{
					TotalResult: 3,
				},
				Result: []*org.Organization{
					{
						Id:   ListOrgIinstance.DefaultOrg.Id,
						Name: ListOrgIinstance.DefaultOrg.Name,
						Details: &object_v2.Details{
							Sequence:      ListOrgIinstance.DefaultOrg.GetDetails().GetSequence(),
							ChangeDate:    ListOrgIinstance.DefaultOrg.GetDetails().GetChangeDate(),
							ResourceOwner: ListOrgIinstance.DefaultOrg.GetDetails().GetResourceOwner(),
							CreationDate:  ListOrgIinstance.DefaultOrg.GetDetails().GetCreationDate(),
						},
					},
					{
						Id:      orgs[0].OrganizationId,
						Name:    orgsName[0],
						Details: orgs[0].GetDetails(),
					},
					{
						Id:      orgs[2].OrganizationId,
						Name:    orgsName[2],
						Details: orgs[2].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations by state inactive",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_StateQuery{
						StateQuery: &org.OrganizationStateQuery{
							State: org.OrganizationState_ORGANIZATION_STATE_INACTIVE,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object_v2.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations by id bad id",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_IdQuery{
						IdQuery: &org.OrganizationIDQuery{
							Id: "bad id",
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
				Result: nil,
			},
		},
		{
			name: "list organizations specify org name equals",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_NameQuery{
						NameQuery: &org.OrganizationNameQuery{
							Name:   orgsName[1],
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations specify org name contains",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_NameQuery{
						NameQuery: &org.OrganizationNameQuery{
							Name: func() string {
								return orgsName[1][1 : len(orgsName[1])-2]
							}(),
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations specify org name contains IGNORE CASE",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_NameQuery{
						NameQuery: &org.OrganizationNameQuery{
							Name: func() string {
								return strings.ToUpper(orgsName[1][1 : len(orgsName[1])-2])
							}(),
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations specify domain name equals",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_DomainQuery{
						DomainQuery: &org.OrganizationDomainQuery{
							Domain: orgsDomain[1],
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations specify domain name contains",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_DomainQuery{
						DomainQuery: &org.OrganizationDomainQuery{
							Domain: orgsDomain[1][1 : len(orgsDomain[1])-2],
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
		{
			name: "list organizations specify org name contains IGNORE CASE",
			ctx:  listOrgIAmOwnerCtx,
			query: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_DomainQuery{
						DomainQuery: &org.OrganizationDomainQuery{
							Domain: strings.ToUpper(orgsDomain[1][1 : len(orgsDomain[1])-2]),
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
						},
					},
				},
			},
			want: &org.ListOrganizationsResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				Result: []*org.Organization{
					{
						Id:      orgs[1].OrganizationId,
						Name:    orgsName[1],
						Details: orgs[1].GetDetails(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := listOrgClient.ListOrganizations(tt.ctx, &org.ListOrganizationsRequest{
					Queries: tt.query,
					Query: &object.ListQuery{
						Asc: true,
					},
					SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_CREATION_DATE,
				})
				if tt.err != nil {
					require.ErrorContains(ttt, err, tt.err.Error())
					return
				}
				require.NoError(ttt, err)

				integration.AssertListDetails(ttt, tt.want, got)

				require.Len(ttt, got.Result, len(tt.want.Result))
				for i, got := range got.Result {
					integration.AssertDetails(t, tt.want.Result[i], got)

					assert.Equal(ttt, tt.want.Result[i].Id, got.Id)
					assert.Equal(ttt, tt.want.Result[i].Name, got.Name)
				}
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_ListOrganizationDomains(t *testing.T) {
	domain := integration.DomainName()

	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId

	var primaryDomain string
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Second)
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		organizations, err := Client.ListOrganizations(CTX, &org.ListOrganizationsRequest{
			Queries: []*org.SearchQuery{
				{Query: &org.SearchQuery_IdQuery{
					IdQuery: &org.OrganizationIDQuery{Id: orgId},
				}},
			},
		})
		require.NoError(t, err)
		require.Len(t, organizations.GetResult(), 1)
		primaryDomain = organizations.GetResult()[0].GetPrimaryDomain()
	}, retryDuration, tick, "could not find primary domain")

	_, err := Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
		OrganizationId: orgId,
		Domain:         domain,
	})
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		request *org.ListOrganizationDomainsRequest
	}
	type want struct {
		response *org.ListOrganizationDomainsResponse
		err      bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "non existing organization",
			args: args{
				ctx:     CTX,
				request: &org.ListOrganizationDomainsRequest{OrganizationId: "not-existing"},
			},
			want: want{
				response: &org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 0,
					},
					Domains: nil,
				},
			},
		},
		{
			name: "no permission (different organization), error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			},
			want: want{
				err: true,
			},
		},
		{
			name: "list org domain, all domains",
			args: args{
				ctx: CTX,
				request: &org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
					SortingColumn:  org.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				},
			},
			want: want{
				response: &org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 2,
					},
					Domains: []*org.Domain{
						{
							OrganizationId: orgId,
							Domain:         primaryDomain,
							IsVerified:     true,
							IsPrimary:      true,
							ValidationType: 0,
						},
						{
							OrganizationId: orgId,
							Domain:         domain,
							IsVerified:     true,
							IsPrimary:      false,
							ValidationType: 0,
						},
					},
				},
			},
		},
		{
			name: "list specific domain",
			args: args{
				ctx: CTX,
				request: &org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
					Filters: []*org.DomainSearchFilter{
						{Filter: &org.DomainSearchFilter_DomainFilter{DomainFilter: &org.OrganizationDomainQuery{Domain: domain}}},
					},
				},
			},
			want: want{
				response: &org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 1,
					},
					Domains: []*org.Domain{
						{
							OrganizationId: orgId,
							Domain:         domain,
							IsVerified:     true,
							IsPrimary:      false,
							ValidationType: 0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, 5*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				queryRes, err := Client.ListOrganizationDomains(tt.args.ctx, tt.args.request)
				if tt.want.err {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				assert.Len(ttt, queryRes.Domains, int(tt.want.response.GetPagination().GetTotalResult()))
				assert.EqualExportedValues(ttt, tt.want.response.GetPagination(), queryRes.GetPagination())
				assert.ElementsMatch(ttt, tt.want.response.GetDomains(), queryRes.GetDomains())
			}, retryDuration, tick, "timeout waiting for adding domain")
		})
	}
}

func TestServer_ListOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId
	setRespoonse, err := Client.SetOrganizationMetadata(CTX, &org.SetOrganizationMetadataRequest{
		OrganizationId: orgId,
		Metadata: []*org.Metadata{
			{
				Key:   "key1",
				Value: []byte("value1"),
			},
			{
				Key:   "key2",
				Value: []byte("value2"),
			},
			{
				Key:   "key2.1",
				Value: []byte("value3"),
			},
			{
				Key:   "key2.2",
				Value: []byte("value4"),
			},
		},
	})
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		request *org.ListOrganizationMetadataRequest
	}
	type want struct {
		response *org.ListOrganizationMetadataResponse
		err      error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "list org metadata happy path",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 4,
					},
					Metadata: []*metadata.Metadata{
						{
							Key:          "key1",
							Value:        []byte("value1"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2",
							Value:        []byte("value2"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.1",
							Value:        []byte("value3"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.2",
							Value:        []byte("value4"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
					},
				},
			},
		},
		{
			name: "list org metadata filter key",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
					Pagination: &filter.PaginationRequest{
						Offset: 1,
						Limit:  2,
					},
					Filters: []*metadata.MetadataSearchFilter{
						{
							Filter: &metadata.MetadataSearchFilter_KeyFilter{
								KeyFilter: &metadata.MetadataKeyFilter{
									Key:    "key2",
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH,
								},
							},
						},
					},
				},
			},
			want: want{
				response: &org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult:  3,
						AppliedLimit: 2,
					},
					Metadata: []*metadata.Metadata{
						{
							Key:          "key2.1",
							Value:        []byte("value3"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.2",
							Value:        []byte("value4"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
					},
				},
			},
		},
		{
			name: "list org metadata for non existent org",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &org.ListOrganizationMetadataRequest{
					OrganizationId: "non existent orgid",
				},
			},
			want: want{
				response: &org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{},
				},
			},
		},
		{
			name: "list org metadata without permission (other organization)",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				request: &org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 1*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListOrganizationMetadata(tt.args.ctx, tt.args.request)
				require.NoError(ttt, err)

				assert.EqualExportedValues(ttt, tt.want.response, got)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}
