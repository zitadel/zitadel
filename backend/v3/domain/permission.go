package domain

import "context"

// PermissionRepository is the interface that manages and checks Zitadel permissions.
//
// TODO(muhlemmer): This just defines the checker methods, rest to be done later.
type PermissionRepository interface {
	PermissionChecker
}

// PermissionChecker defines the methods needed to check permissions.
type PermissionChecker interface {
	// Check if the authenticated user has the given permission on instance level.
	CheckInstancePermission(ctx context.Context, permission string) error

	// Check if the authenticated user has the given permission on the given organization.
	// A permission may be inherited from the instance level.
	CheckOrganizationPermission(ctx context.Context, permission, orgID string) error

	// Check if the authenticated user has the given permission on the given project.
	// A permission may be inherited from the instance or organization level.
	CheckProjectPermission(ctx context.Context, permission, projectID string) error

	// Check if the authenticated user has the given permission on the given project grant.
	// A permission may be inherited from the instance or granted organization level.
	CheckProjectGrantPermission(ctx context.Context, permission, projectGrantID string) error
}

type noopPermissionChecker struct{}

var _ PermissionChecker = (*noopPermissionChecker)(nil)

func (*noopPermissionChecker) CheckInstancePermission(context.Context, string) error {
	return nil
}

func (*noopPermissionChecker) CheckOrganizationPermission(context.Context, string, string) error {
	return nil
}

func (*noopPermissionChecker) CheckProjectPermission(context.Context, string, string) error {
	return nil
}

func (*noopPermissionChecker) CheckProjectGrantPermission(context.Context, string, string) error {
	return nil
}
