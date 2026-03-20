package domain

import "context"

const (
	InstanceReadPermission      = "instance.read"
	InstanceWritePermission     = "instance.write"
	DomainReadPermission        = "domain.read"
	DomainWritePermission       = "domain.write"
	OrganizationReadPermission  = "organization.read"
	OrganizationWritePermission = "organization.write"
	SessionReadPermission       = "session.read"
	SessionWritePermission      = "session.write"
	SessionDeletePermission     = "session.delete"
)

// PermissionRepository is the interface that manages and checks Zitadel permissions.
//
// TODO(muhlemmer): This just defines the checker methods, rest to be done later.
type PermissionRepository interface {
	PermissionChecker
}

// PermissionChecker defines the methods needed to check permissions.
//
//go:generate mockgen -typed -package domainmock -destination ./mock/permission.mock.go . PermissionChecker
type PermissionChecker interface {
	// Check if the authenticated user has the given permission on instance level.
	CheckInstancePermission(ctx context.Context, opts *InvokeOpts, permission string) error

	// Check if the authenticated user has the given permission on the given organization.
	// A permission may be inherited from the instance level.
	CheckOrganizationPermission(ctx context.Context, opts *InvokeOpts, permission, orgID string) error

	// Check if the authenticated user has the given permission on the given project.
	// A permission may be inherited from the instance or organization level.
	CheckProjectPermission(ctx context.Context, opts *InvokeOpts, permission, projectID string) error

	// Check if the authenticated user has the given permission on the given project grant.
	// A permission may be inherited from the instance or granted organization level.
	CheckProjectGrantPermission(ctx context.Context, opts *InvokeOpts, permission, projectGrantID string) error

	// Check if the user in context has the given permission on the given session.
	// A permission may be inherited from the instance or granted organization level.
	CheckSessionPermission(ctx context.Context, opts *InvokeOpts, permission, sessionID string) error
}

type noopPermissionChecker struct{}

// CheckSessionPermission implements [PermissionChecker].
func (n *noopPermissionChecker) CheckSessionPermission(ctx context.Context, opts *InvokeOpts, permission, sessionID string) error {
	return nil
}

var _ PermissionChecker = (*noopPermissionChecker)(nil)

func (*noopPermissionChecker) CheckInstancePermission(context.Context, *InvokeOpts, string) error {
	return nil
}

func (*noopPermissionChecker) CheckOrganizationPermission(context.Context, *InvokeOpts, string, string) error {
	return nil
}

func (*noopPermissionChecker) CheckProjectPermission(context.Context, *InvokeOpts, string, string) error {
	return nil
}

func (*noopPermissionChecker) CheckProjectGrantPermission(context.Context, *InvokeOpts, string, string) error {
	return nil
}
