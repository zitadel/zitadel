//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_CreateProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())

	alreadyExistingProjectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), gofakeit.AppName(), false, false)
	alreadyExistingGrantedOrgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())
	instance.CreateProjectGrant(iamOwnerCtx, alreadyExistingProjectResp.GetId(), alreadyExistingGrantedOrgResp.GetOrganizationId())

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *project.CreateProjectGrantRequest)
		req     *project.CreateProjectGrantRequest
		want
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), gofakeit.AppName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name:    "empty projectID",
			ctx:     iamOwnerCtx,
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "empty granted organization",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectGrantRequest{
				ProjectId:             "something",
				GrantedOrganizationId: "",
			},
			wantErr: true,
		},
		{
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectGrantRequest{
				ProjectId:             alreadyExistingProjectResp.GetId(),
				GrantedOrganizationId: alreadyExistingGrantedOrgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "empty, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), gofakeit.AppName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.CreateProjectGrant(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateProjectGrantResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertCreateProjectGrantResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate bool, actualResp *project.CreateProjectGrantResponse) {
	if expectedCreationDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.CreationDate)
	}
}
