package domain

import "context"

type Permissions struct {
	Permissions []string
}

func (p *Permissions) AppendPermissions(ctxID string, permissions ...string) {
	for _, permission := range permissions {
		p.appendPermission(ctxID, permission)
	}
}

func (p *Permissions) appendPermission(ctxID, permission string) {
	if ctxID != "" {
		permission = permission + ":" + ctxID
	}
	for _, existingPermission := range p.Permissions {
		if existingPermission == permission {
			return
		}
	}
	p.Permissions = append(p.Permissions, permission)
}

type PermissionCheck func(ctx context.Context, permission, orgID, resourceID string) (err error)

const (
	PermissionUserWrite           = "user.write"
	PermissionUserRead            = "user.read"
	PermissionUserDelete          = "user.delete"
	PermissionUserCredentialWrite = "user.credential.write"
	PermissionSessionWrite        = "session.write"
	PermissionSessionRead         = "session.read"
	PermissionSessionLink         = "session.link"
	PermissionSessionDelete       = "session.delete"
	PermissionOrgRead             = "org.read"
	PermissionIDPRead             = "iam.idp.read"
	PermissionOrgIDPRead          = "org.idp.read"
	PermissionProjectWrite        = "project.write"
	PermissionProjectRead         = "project.read"
	PermissionProjectDelete       = "project.delete"
	PermissionProjectGrantWrite   = "project.grant.write"
	PermissionProjectGrantRead    = "project.grant.read"
	PermissionProjectGrantDelete  = "project.grant.delete"
	PermissionProjectRoleWrite    = "project.role.write"
	PermissionProjectRoleRead     = "project.role.read"
	PermissionProjectRoleDelete   = "project.role.delete"
)

// ProjectPermissionCheck is used as a check for preconditions dependent on application, project, user resourceowner and usergrants.
// Configurable on the project the application belongs to through the flags related to authentication.
type ProjectPermissionCheck func(ctx context.Context, clientID, userID string) (err error)
