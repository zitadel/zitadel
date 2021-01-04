package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

type UserGrantRepo struct {
	View         *view.View
	IamID        string
	IamProjectID string
	Auth         authz.Config
	IamEvents    *iam_event.IAMEventstore
}

func (repo *UserGrantRepo) Health() error {
	return repo.View.Health()
}

func (repo *UserGrantRepo) ResolveGrants(ctx context.Context) (*authz.Grant, error) {
	err := repo.FillIamProjectID(ctx)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)

	orgGrant, err := repo.View.UserGrantByIDs(ctxData.OrgID, repo.IamProjectID, ctxData.UserID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	iamAdminGrant, err := repo.View.UserGrantByIDs(repo.IamID, repo.IamProjectID, ctxData.UserID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}

	return mergeOrgAndAdminGrant(ctxData, orgGrant, iamAdminGrant), nil
}

func (repo *UserGrantRepo) SearchMyZitadelPermissions(ctx context.Context) ([]string, error) {
	grant, err := repo.ResolveGrants(ctx)
	if err != nil {
		return nil, err
	}

	if grant == nil {
		return []string{}, nil
	}
	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, role := range grant.Roles {
		roleName, ctxID := authz.SplitPermission(role)
		for _, mapping := range repo.Auth.RolePermissionMappings {
			if mapping.Role == roleName {
				permissions.AppendPermissions(ctxID, mapping.Permissions...)
			}
		}
	}
	return permissions.Permissions, nil
}

func (repo *UserGrantRepo) FillIamProjectID(ctx context.Context) error {
	if repo.IamProjectID != "" {
		return nil
	}
	iam, err := repo.IamEvents.IAMByID(ctx, repo.IamID)
	if err != nil {
		return err
	}
	if iam.SetUpDone < domain.StepCount-1 {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-skiwS", "Setup not done")
	}
	repo.IamProjectID = iam.IAMProjectID
	return nil
}

func mergeOrgAndAdminGrant(ctxData authz.CtxData, orgGrant, iamAdminGrant *model.UserGrantView) (grant *authz.Grant) {
	if orgGrant != nil {
		roles := orgGrant.RoleKeys
		if iamAdminGrant != nil {
			roles = addIamAdminRoles(roles, iamAdminGrant.RoleKeys)
		}
		grant = &authz.Grant{OrgID: orgGrant.ResourceOwner, Roles: roles}
	} else if iamAdminGrant != nil {
		grant = &authz.Grant{
			OrgID: ctxData.OrgID,
			Roles: iamAdminGrant.RoleKeys,
		}
	}
	return grant
}

func addIamAdminRoles(orgRoles, iamAdminRoles []string) []string {
	result := make([]string, 0)
	result = append(result, iamAdminRoles...)
	for _, role := range orgRoles {
		if !authz.ExistsPerm(result, role) {
			result = append(result, role)
		}
	}
	return result
}
