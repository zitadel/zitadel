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

type PermissionCheck func(ctx context.Context, permission, resourceOwnerID, aggregateID string) (err error)

const (
	PermissionUserWrite                = "user.write"
	PermissionUserRead                 = "user.read"
	PermissionUserDelete               = "user.delete"
	PermissionUserCredentialWrite      = "user.credential.write"
	PermissionSessionWrite             = "session.write"
	PermissionSessionRead              = "session.read"
	PermissionSessionLink              = "session.link"
	PermissionSessionDelete            = "session.delete"
	PermissionOrgRead                  = "org.read"
	PermissionIDPRead                  = "iam.idp.read"
	PermissionOrgIDPRead               = "org.idp.read"
	PermissionProjectCreate            = "project.create"
	PermissionProjectWrite             = "project.write"
	PermissionProjectRead              = "project.read"
	PermissionProjectDelete            = "project.delete"
	PermissionProjectGrantWrite        = "project.grant.write"
	PermissionProjectGrantRead         = "project.grant.read"
	PermissionProjectGrantDelete       = "project.grant.delete"
	PermissionProjectRoleWrite         = "project.role.write"
	PermissionProjectRoleRead          = "project.role.read"
	PermissionProjectRoleDelete        = "project.role.delete"
	PermissionProjectAppWrite          = "project.app.write"
	PermissionProjectAppDelete         = "project.app.delete"
	PermissionProjectAppRead           = "project.app.read"
	PermissionInstanceMemberWrite      = "iam.member.write"
	PermissionInstanceMemberDelete     = "iam.member.delete"
	PermissionInstanceMemberRead       = "iam.member.read"
	PermissionOrgMemberWrite           = "org.member.write"
	PermissionOrgMemberDelete          = "org.member.delete"
	PermissionOrgMemberRead            = "org.member.read"
	PermissionProjectMemberWrite       = "project.member.write"
	PermissionProjectMemberDelete      = "project.member.delete"
	PermissionProjectMemberRead        = "project.member.read"
	PermissionProjectGrantMemberWrite  = "project.grant.member.write"
	PermissionProjectGrantMemberDelete = "project.grant.member.delete"
	PermissionProjectGrantMemberRead   = "project.grant.member.read"
	PermissionUserGrantWrite           = "user.grant.write"
	PermissionUserGrantRead            = "user.grant.read"
	PermissionUserGrantDelete          = "user.grant.delete"
)

// ProjectPermissionCheck is used as a check for preconditions dependent on application, project, user resourceowner and usergrants.
// Configurable on the project the application belongs to through the flags related to authentication.
type ProjectPermissionCheck func(ctx context.Context, clientID, userID string) (err error)
