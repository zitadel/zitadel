package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrantRepo struct {
	View         *view.View
	IamID        string
	IamProjectID string
	Auth         auth.Config
	IamEvents    *iam_event.IamEventstore
}

func (repo *UserGrantRepo) Health() error {
	return repo.View.Health()
}

func (repo *UserGrantRepo) ResolveGrants(ctx context.Context) (*auth.Grant, error) {
	err := repo.FillIamProjectID(ctx)
	if err != nil {
		return nil, err
	}
	ctxData := auth.GetCtxData(ctx)

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

	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, role := range grant.Roles {
		roleName, ctxID := auth.SplitPermission(role)
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
	iam, err := repo.IamEvents.IamByID(ctx, repo.IamID)
	if err != nil {
		return err
	}
	if !iam.SetUpDone {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-skiwS", "Setup not done")
	}
	repo.IamProjectID = iam.IamProjectID
	return nil
}

func mergeOrgAndAdminGrant(ctxData auth.CtxData, orgGrant, iamAdminGrant *model.UserGrantView) (grant *auth.Grant) {
	if orgGrant != nil {
		roles := orgGrant.RoleKeys
		if iamAdminGrant != nil {
			roles = addIamAdminRoles(roles, iamAdminGrant.RoleKeys)
		}
		grant = &auth.Grant{OrgID: orgGrant.ResourceOwner, Roles: roles}
	} else if iamAdminGrant != nil {
		grant = &auth.Grant{
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
		if !auth.ExistsPerm(result, role) {
			result = append(result, role)
		}
	}
	return result
}
