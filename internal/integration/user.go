package integration

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	internal_permission_v2 "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func (i *Instance) CreateMachineUserPATWithMembership(ctx context.Context, roles ...string) (id, pat string, err error) {
	user := i.CreateMachineUser(ctx)

	patResp, err := i.Client.Mgmt.AddPersonalAccessToken(ctx, &management.AddPersonalAccessTokenRequest{
		UserId:         user.GetUserId(),
		ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour)),
	})
	if err != nil {
		return "", "", err
	}

	orgRoles := make([]string, 0, len(roles))
	iamRoles := make([]string, 0, len(roles))

	for _, role := range roles {
		if strings.HasPrefix(role, "ORG_") {
			orgRoles = append(orgRoles, role)
		}
		if strings.HasPrefix(role, "IAM_") {
			iamRoles = append(iamRoles, role)
		}
	}

	if len(orgRoles) > 0 {
		_, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
			Resource: &internal_permission_v2.ResourceType{
				Resource: &internal_permission_v2.ResourceType_OrganizationId{OrganizationId: i.DefaultOrg.GetId()},
			},
			UserId: user.GetUserId(),
			Roles:  orgRoles,
		})
		if err != nil {
			return "", "", err
		}
	}
	if len(iamRoles) > 0 {
		_, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
			Resource: &internal_permission_v2.ResourceType{
				Resource: &internal_permission_v2.ResourceType_Instance{Instance: true},
			},
			UserId: user.GetUserId(),
			Roles:  iamRoles,
		})
		if err != nil {
			return "", "", err
		}
	}

	return user.GetUserId(), patResp.GetToken(), nil
}
