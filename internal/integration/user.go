package integration

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Tester) CreateMachineUserPATWithMembership(ctx context.Context, roles ...string) (id, pat string, err error) {
	user := s.CreateMachineUser(ctx)

	patResp, err := s.Client.Mgmt.AddPersonalAccessToken(ctx, &management.AddPersonalAccessTokenRequest{
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
		_, err := s.Client.Mgmt.AddOrgMember(ctx, &management.AddOrgMemberRequest{
			UserId: user.GetUserId(),
			Roles:  orgRoles,
		})
		if err != nil {
			return "", "", err
		}
	}
	if len(iamRoles) > 0 {
		_, err := s.Client.Admin.AddIAMMember(ctx, &admin.AddIAMMemberRequest{
			UserId: user.GetUserId(),
			Roles:  iamRoles,
		})
		if err != nil {
			return "", "", err
		}
	}

	return user.GetUserId(), patResp.GetToken(), nil
}
