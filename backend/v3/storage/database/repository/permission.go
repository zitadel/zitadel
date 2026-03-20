package repository

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type permission struct{}

// CheckInstancePermission implements [domain.PermissionChecker].
func (p permission) CheckInstancePermission(ctx context.Context, opts *domain.InvokeOpts, permission string) error {
	return p.executeCheck(ctx, opts.DB(),
		PermissionCondition(
			authz.GetInstance(ctx).InstanceID(),
			authz.GetCtxData(ctx).UserID,
			permission,
		),
	)
}

// CheckOrganizationPermission implements [domain.PermissionChecker].
func (p permission) CheckOrganizationPermission(ctx context.Context, opts *domain.InvokeOpts, permission string, orgID string) error {
	return p.executeCheck(ctx, opts.DB(),
		PermissionCondition(
			authz.GetInstance(ctx).InstanceID(),
			authz.GetCtxData(ctx).UserID,
			permission,
			WithOrganizationID(orgID),
		),
	)
}

// CheckProjectGrantPermission implements [domain.PermissionChecker].
func (p permission) CheckProjectGrantPermission(ctx context.Context, opts *domain.InvokeOpts, permission string, projectGrantID string) error {
	return p.executeCheck(ctx, opts.DB(),
		PermissionCondition(
			authz.GetInstance(ctx).InstanceID(),
			authz.GetCtxData(ctx).UserID,
			permission,
			WithProjectGrantID(projectGrantID),
		),
	)
}

// CheckProjectPermission implements [domain.PermissionChecker].
func (p permission) CheckProjectPermission(ctx context.Context, opts *domain.InvokeOpts, permission string, projectID string) error {
	return p.executeCheck(ctx, opts.DB(),
		PermissionCondition(
			authz.GetInstance(ctx).InstanceID(),
			authz.GetCtxData(ctx).UserID,
			permission,
			WithProjectID(projectID),
		),
	)
}

// CheckSessionPermission implements [domain.PermissionChecker].
func (p permission) CheckSessionPermission(ctx context.Context, opts *domain.InvokeOpts, permission string, sessionID string) error {
	panic("unimplemented")
}

func (p permission) executeCheck(ctx context.Context, client database.QueryExecutor, check database.Condition) error {
	var hasPermission bool
	builder := database.NewStatementBuilder("SELECT ")
	check.Write(builder)

	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&hasPermission)
	if err != nil {
		return err
	}

	if !hasPermission {
		return zerrors.ThrowPermissionDenied(nil, "PERM-CheckInstancePermission", "Errors.PermissionDenied")
	}

	return nil
}

func Permission() *permission {
	return new(permission)
}

var _ domain.PermissionChecker = (*permission)(nil)

type CheckPermissionOpt func(*permissionCondition)

func WithOrganizationID(organizationID string) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.organizationID = func(builder *database.StatementBuilder) {
			builder.WriteArgs(organizationID)
		}
	}
}

func WithOrganizationIDColumn(col database.Column) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.organizationID = col.WriteQualified
	}
}

func WithProjectID(projectID string) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.projectID = func(builder *database.StatementBuilder) {
			builder.WriteArgs(projectID)
		}
	}
}

func WithProjectIDColumn(col database.Column) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.projectID = col.WriteQualified
	}
}

func WithProjectGrantID(projectGrantID string) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.projectGrantID = func(builder *database.StatementBuilder) {
			builder.WriteArgs(projectGrantID)
		}
	}
}

func WithProjectGrantIDColumn(col database.Column) CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.projectGrantID = col.WriteQualified
	}
}

func WithRaiseIfDenied() CheckPermissionOpt {
	return func(p *permissionCondition) {
		p.raiseIfDenied = true
	}
}

func PermissionCondition(instanceID, userID, permission string, opts ...CheckPermissionOpt) database.Condition {
	cond := &permissionCondition{
		instanceID: instanceID,
		userID:     userID,
		permission: permission,
	}
	for _, opt := range opts {
		opt(cond)
	}
	return cond
}

type permissionCondition struct {
	instanceID     string
	userID         string
	permission     string
	raiseIfDenied  bool
	organizationID func(builder *database.StatementBuilder)
	projectID      func(builder *database.StatementBuilder)
	projectGrantID func(builder *database.StatementBuilder)
}

// IsRestrictingColumn implements [database.Condition].
func (p *permissionCondition) IsRestrictingColumn(col database.Column) bool {
	return false
}

// Matches implements [database.Condition].
func (p *permissionCondition) Matches(x any) bool {
	toMatch, ok := x.(*permissionCondition)
	if !ok {
		return false
	}
	var builder, toMatchBuilder database.StatementBuilder
	p.Write(&builder)
	toMatch.Write(&toMatchBuilder)
	return builder.String() == toMatchBuilder.String() && slices.Equal(builder.Args(), toMatchBuilder.Args())
}

// String implements [database.Condition].
func (p *permissionCondition) String() string {
	return "permissionCondition"
}

// Write implements [database.Condition].
func (p *permissionCondition) Write(builder *database.StatementBuilder) {
	builder.WriteString("zitadel.check_permission(")
	builder.WriteArgs(p.instanceID, p.userID, p.permission)
	if p.organizationID != nil {
		builder.WriteString(", p_organization_id => ")
		p.organizationID(builder)
	}
	if p.projectID != nil {
		builder.WriteString(", p_project_id => ")
		p.projectID(builder)
	}
	if p.projectGrantID != nil {
		builder.WriteString(", p_project_grant_id => ")
		p.projectGrantID(builder)
	}
	if p.raiseIfDenied {
		builder.WriteString(", p_raise_if_denied => ")
		builder.WriteArgs(true)
	}
	builder.WriteString(")")
}

var _ database.Condition = (*permissionCondition)(nil)
