//go:build integration

package org_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/zitadel/zitadel/internal/integration"

	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

var (
	CTX      context.Context
	Instance *integration.Instance
	Client   v2beta_org.OrganizationServiceClient
	User     *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OrgV2beta

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		User = Instance.CreateHumanUser(CTX)
		return m.Run()
	}())
}

func TestServer_CreateOrganization(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *v2beta_org.CreateOrganizationRequest
		want    *v2beta_org.CreateOrganizationResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &v2beta_org.CreateOrganizationRequest{
				Name:   "name",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name:   "",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid admin type",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: gofakeit.AppName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "admin with init",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: gofakeit.AppName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_Human{
							Human: &user_v2beta.AddHumanUserRequest{
								Profile: &user_v2beta.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user_v2beta.SetHumanEmail{
									Email: gofakeit.Email(),
									Verification: &user_v2beta.SetHumanEmail_ReturnCode{
										ReturnCode: &user_v2beta.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				OrganizationId: integration.NotEmpty,
				CreatedAdmins: []*v2beta_org.CreateOrganizationResponse_CreatedAdmin{
					{
						UserId:    integration.NotEmpty,
						EmailCode: gu.Ptr(integration.NotEmpty),
						PhoneCode: nil,
					},
				},
			},
		},
		{
			name: "existing user and new human with idp",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: gofakeit.AppName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_Human{
							Human: &user_v2beta.AddHumanUserRequest{
								Profile: &user_v2beta.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user_v2beta.SetHumanEmail{
									Email: gofakeit.Email(),
									Verification: &user_v2beta.SetHumanEmail_IsVerified{
										IsVerified: true,
									},
								},
								IdpLinks: []*user_v2beta.IDPLink{
									{
										IdpId:    idpResp.Id,
										UserId:   "userID",
										UserName: "username",
									},
								},
							},
						},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				CreatedAdmins: []*v2beta_org.CreateOrganizationResponse_CreatedAdmin{
					// a single admin is expected, because the first provided already exists
					{
						UserId: integration.NotEmpty,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check details
			assert.NotZero(t, got.GetDetails().GetSequence())
			gotCD := got.GetDetails().GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEmpty(t, got.GetDetails().GetResourceOwner())

			// organization id must be the same as the resourceOwner
			assert.Equal(t, got.GetDetails().GetResourceOwner(), got.GetOrganizationId())

			// check the admins
			require.Len(t, got.GetCreatedAdmins(), len(tt.want.GetCreatedAdmins()))
			for i, admin := range tt.want.GetCreatedAdmins() {
				gotAdmin := got.GetCreatedAdmins()[i]
				assertCreatedAdmin(t, admin, gotAdmin)
			}
		})
	}
}

func TestServer_UpdateOrganization(t *testing.T) {
	orgs, orgsName, err := createOrgs(1)
	if err != nil {
		assert.Fail(t, "unable to create org")
	}
	orgId := orgs[0].OrganizationId
	orgName := orgsName[0]

	tests := []struct {
		name    string
		ctx     context.Context
		req     *v2beta_org.UpdateOrganizationRequest
		want    *v2beta_org.UpdateOrganizationResponse
		wantErr bool
	}{
		{
			name: "update org with new name",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id:   orgId,
				Name: "new org name",
			},
		},
		{
			name: "update org with same name",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id:   orgId,
				Name: orgName,
			},
		},
		{
			name: "update org with no id",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id: orgId,
				// Name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check details
			assert.NotZero(t, got.GetDetails().GetSequence())
			gotCD := got.GetDetails().GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEmpty(t, got.GetDetails().GetResourceOwner())
		})
	}
}

func TestServer_GetOrganizationByID(t *testing.T) {
	orgs, orgsName, err := createOrgs(1)
	if err != nil {
		assert.Fail(t, "unable to create org")
	}
	orgId := orgs[0].OrganizationId
	orgName := orgsName[0]

	tests := []struct {
		name    string
		ctx     context.Context
		req     *v2beta_org.GetOrganizationByIDRequest
		want    *v2beta_org.GetOrganizationByIDResponse
		wantErr bool
	}{
		{
			name: "get organization happy path",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.GetOrganizationByIDRequest{
				Id: orgId,
			},
		},
		{
			name: "get organization that doesn't exist",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.GetOrganizationByIDRequest{
				Id: "non existing organization",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GetOrganizationByID(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, orgId, got.Organization.Id)
			require.Equal(t, orgName, got.Organization.Name)
		})
	}
}

// TODO: finish off qyery testing in ListOrganizations
func TestServer_ListOrganization(t *testing.T) {
	noOfOrgs := 3
	orgs, orgsName, err := createOrgs(noOfOrgs)
	if err != nil {
		assert.Fail(t, "unable to create orgs")
	}

	tests := []struct {
		name    string
		ctx     context.Context
		req     *v2beta_org.ListOrganizationsRequest
		want    []*v2beta_org.Organization
		wantErr bool
	}{
		{
			name: "list organizations happy path",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req:  &v2beta_org.ListOrganizationsRequest{},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[0].OrganizationId,
					Name: orgsName[0],
				},
				{
					Id:   orgs[1].OrganizationId,
					Name: orgsName[1],
				},
				{
					Id:   orgs[2].OrganizationId,
					Name: orgsName[2],
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(context.Background(), 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListOrganizations(tt.ctx, tt.req)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				// require.Equal(t, len(tt.want), len(got.Result))

				// check details
				// assert.NotZero(t, got.GetDetails().GetSequence())
				// gotCD := got.GetDetails().GetChangeDate().AsTime()
				// now := time.Now()
				// assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
				// assert.NotEmpty(t, got.GetDetails().GetResourceOwner())

				foundOrgs := 0
				for _, got := range got.Result {
					for _, org := range tt.want {
						if org.Name == got.Name &&
							org.Id == got.Id {
							foundOrgs += 1
						}
						// require.Equal(t, org.do, got.Result[i].Name)
					}
				}
				require.Equal(t, len(tt.want), foundOrgs)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_DeleteOrganization(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		createOrgFunc func() string
		req           *v2beta_org.DeleteOrganizationRequest
		want          *v2beta_org.DeleteOrganizationResponse
		err           error
	}{
		{
			name: "delete org happy path",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			createOrgFunc: func() string {
				orgs, _, err := createOrgs(1)
				if err != nil {
					assert.Fail(t, "unable to create org")
				}
				return orgs[0].OrganizationId
			},
			req: &v2beta_org.DeleteOrganizationRequest{},
		},
		{
			name: "delete non existent org",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner),
			req: &v2beta_org.DeleteOrganizationRequest{
				Id: "non existent org id",
			},
			err: fmt.Errorf("Organisation not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createOrgFunc != nil {
				tt.req.Id = tt.createOrgFunc()
			}

			got, err := Client.DeleteOrganization(tt.ctx, tt.req)
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			// check details
			assert.NotZero(t, got.GetDetails().GetSequence())
			gotCD := got.GetDetails().GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEmpty(t, got.GetDetails().GetResourceOwner())

			_, err = Client.GetOrganizationByID(tt.ctx, &v2beta_org.GetOrganizationByIDRequest{
				Id: tt.req.Id,
			})
			require.Contains(t, err.Error(), "Organisation not found")
		})
	}
}

func TestServer_DeactivateReactivateNonExistentOrganization(t *testing.T) {
	ctx := Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)

	// deactivate non existent organization
	_, err := Client.DeactivateOrganization(ctx, &v2beta_org.DeactivateOrganizationRequest{
		Id: "non existent organization",
	})
	require.Contains(t, err.Error(), "Organisation not found")

	// reactivate non existent organization
	_, err = Client.ReactivateOrganization(ctx, &v2beta_org.ReactivateOrganizationRequest{
		Id: "non existent organization",
	})
	require.Contains(t, err.Error(), "Organisation not found")
}

func TestServer_DeactivateReactivateOrganization(t *testing.T) {
	// 1. create organization
	orgs, _, err := createOrgs(1)
	if err != nil {
		assert.Fail(t, "unable to create orgs")
	}
	orgId := orgs[0].OrganizationId
	ctx := Instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)

	// 2. check inital state of organization
	res, err := Client.GetOrganizationByID(ctx, &org.GetOrganizationByIDRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	require.Equal(t, v2beta_org.OrganizationState_ORGANIZATION_STATE_ACTIVE, res.Organization.State)

	// 3. deactivate organization once
	deactivate_res, err := Client.DeactivateOrganization(ctx, &v2beta_org.DeactivateOrganizationRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	assert.NotZero(t, deactivate_res.GetDetails().GetSequence())
	gotCD := deactivate_res.GetDetails().GetChangeDate().AsTime()
	now := time.Now()
	assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
	assert.NotEmpty(t, deactivate_res.GetDetails().GetResourceOwner())

	// 4. check organization state is deactivated
	res, err = Client.GetOrganizationByID(ctx, &v2beta_org.GetOrganizationByIDRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	require.Equal(t, v2beta_org.OrganizationState_ORGANIZATION_STATE_INACTIVE, res.Organization.State)

	// 5. repeat deactivate organization once
	deactivate_res, err = Client.DeactivateOrganization(ctx, &v2beta_org.DeactivateOrganizationRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	assert.NotZero(t, deactivate_res.GetDetails().GetSequence())
	gotCD = deactivate_res.GetDetails().GetChangeDate().AsTime()
	now = time.Now()
	assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
	assert.NotEmpty(t, deactivate_res.GetDetails().GetResourceOwner())

	// 6. repeat check organization state is still deactivated
	res, err = Client.GetOrganizationByID(ctx, &v2beta_org.GetOrganizationByIDRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	require.Equal(t, v2beta_org.OrganizationState_ORGANIZATION_STATE_INACTIVE, res.Organization.State)

	// 7. reactivate organization
	reactivate_res, err := Client.ReactivateOrganization(ctx, &v2beta_org.ReactivateOrganizationRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	assert.NotZero(t, reactivate_res.GetDetails().GetSequence())
	gotCD = reactivate_res.GetDetails().GetChangeDate().AsTime()
	now = time.Now()
	assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
	assert.NotEmpty(t, reactivate_res.GetDetails().GetResourceOwner())

	// 8. check organization state is active
	res, err = Client.GetOrganizationByID(ctx, &v2beta_org.GetOrganizationByIDRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	require.Equal(t, v2beta_org.OrganizationState_ORGANIZATION_STATE_ACTIVE, res.Organization.State)

	// 9. repeat reactivate organization
	reactivate_res, err = Client.ReactivateOrganization(ctx, &v2beta_org.ReactivateOrganizationRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	assert.NotZero(t, reactivate_res.GetDetails().GetSequence())
	gotCD = reactivate_res.GetDetails().GetChangeDate().AsTime()
	now = time.Now()
	assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
	assert.NotEmpty(t, reactivate_res.GetDetails().GetResourceOwner())

	// 10. repeat check organization state is still active
	res, err = Client.GetOrganizationByID(ctx, &v2beta_org.GetOrganizationByIDRequest{
		Id: orgId,
	})
	require.NoError(t, err)
	require.Equal(t, v2beta_org.OrganizationState_ORGANIZATION_STATE_ACTIVE, res.Organization.State)
}

func createOrgs(noOfOrgs int) ([]*v2beta_org.CreateOrganizationResponse, []string, error) {
	var err error
	orgs := make([]*v2beta_org.CreateOrganizationResponse, noOfOrgs)
	orgsName := make([]string, noOfOrgs)

	for i := range noOfOrgs {
		orgName := gofakeit.Name()
		orgsName[i] = orgName
		orgs[i], err = Client.CreateOrganization(CTX,
			&v2beta_org.CreateOrganizationRequest{
				Name: orgName,
			},
		)
		if err != nil {
			return nil, nil, err
		}
	}

	return orgs, orgsName, nil
}

func assertCreatedAdmin(t *testing.T, expected, got *v2beta_org.CreateOrganizationResponse_CreatedAdmin) {
	if expected.GetUserId() != "" {
		assert.NotEmpty(t, got.GetUserId())
	} else {
		assert.Empty(t, got.GetUserId())
	}
	if expected.GetEmailCode() != "" {
		assert.NotEmpty(t, got.GetEmailCode())
	} else {
		assert.Empty(t, got.GetEmailCode())
	}
	if expected.GetPhoneCode() != "" {
		assert.NotEmpty(t, got.GetPhoneCode())
	} else {
		assert.Empty(t, got.GetPhoneCode())
	}
}
