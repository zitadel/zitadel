package query

import (
	"context"
	"database/sql"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Administrators struct {
	SearchResponse
	Administrators []*Administrator
}

type Administrator struct {
	Roles         database.TextArray[string]
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string

	User         *UserAdministrator
	Org          *OrgAdministrator
	Instance     *InstanceAdministrator
	Project      *ProjectAdministrator
	ProjectGrant *ProjectGrantAdministrator
}

type UserAdministrator struct {
	UserID        string
	LoginName     string
	DisplayName   string
	ResourceOwner string
}
type OrgAdministrator struct {
	OrgID string
	Name  string
}

type InstanceAdministrator struct {
	InstanceID string
	Name       string
}

type ProjectAdministrator struct {
	ProjectID     string
	Name          string
	ResourceOwner string
}

type ProjectGrantAdministrator struct {
	ProjectID     string
	ProjectName   string
	GrantID       string
	GrantedOrgID  string
	ResourceOwner string
}

func NewAdministratorUserResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(UserResourceOwnerCol, value, TextEquals)
}

func NewAdministratorUserLoginNameSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(LoginNameNameCol, value, TextEquals)
}

func NewAdministratorUserDisplayNameSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(HumanDisplayNameCol, value, TextEquals)
}

func administratorInstancePermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		InstanceMemberResourceOwner,
		domain.PermissionInstanceMemberRead,
		OwnedRowsPermissionOption(InstanceMemberUserID),
	)
	return query.JoinClause(join, args...)
}

func administratorOrgPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		OrgMemberResourceOwner,
		domain.PermissionOrgMemberRead,
		OwnedRowsPermissionOption(OrgMemberUserID),
	)
	return query.JoinClause(join, args...)
}

func administratorProjectPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		ProjectMemberResourceOwner,
		domain.PermissionProjectMemberRead,
		WithProjectsPermissionOption(ProjectMemberProjectID),
		OwnedRowsPermissionOption(ProjectMemberUserID),
	)
	return query.JoinClause(join, args...)
}

func administratorProjectGrantPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		ProjectGrantMemberResourceOwner,
		domain.PermissionProjectGrantMemberRead,
		WithProjectsPermissionOption(ProjectMemberProjectID),
		OwnedRowsPermissionOption(ProjectGrantMemberUserID),
	)
	return query.JoinClause(join, args...)
}

func administratorsCheckPermission(ctx context.Context, administrators *Administrators, permissionCheck domain.PermissionCheck) {
	selfUserID := authz.GetCtxData(ctx).UserID
	administrators.Administrators = slices.DeleteFunc(administrators.Administrators,
		func(administrator *Administrator) bool {
			if administrator.User != nil && administrator.User.UserID == selfUserID {
				return false
			}
			if administrator.ProjectGrant != nil {
				return administratorProjectGrantCheckPermission(ctx, administrator.ProjectGrant.ResourceOwner, administrator.ProjectGrant.ProjectID, administrator.ProjectGrant.GrantID, administrator.ProjectGrant.GrantedOrgID, permissionCheck) != nil
			}
			if administrator.Project != nil {
				return permissionCheck(ctx, domain.PermissionProjectMemberRead, administrator.Project.ResourceOwner, administrator.Project.ProjectID) != nil
			}
			if administrator.Org != nil {
				return permissionCheck(ctx, domain.PermissionOrgMemberRead, administrator.Org.OrgID, administrator.Org.OrgID) != nil
			}
			if administrator.Instance != nil {
				return permissionCheck(ctx, domain.PermissionInstanceMemberRead, administrator.Instance.InstanceID, administrator.Instance.InstanceID) != nil
			}
			return true
		},
	)
}

func administratorProjectGrantCheckPermission(ctx context.Context, resourceOwner, projectID, grantID, grantedOrgID string, permissionCheck domain.PermissionCheck) error {
	if err := permissionCheck(ctx, domain.PermissionProjectGrantMemberRead, resourceOwner, grantID); err != nil {
		if err := permissionCheck(ctx, domain.PermissionProjectGrantMemberRead, grantedOrgID, grantID); err != nil {
			if err := permissionCheck(ctx, domain.PermissionProjectGrantMemberRead, resourceOwner, projectID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (q *Queries) SearchAdministrators(ctx context.Context, queries *MembershipSearchQuery, permissionCheck domain.PermissionCheck) (*Administrators, error) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	admins, err := q.searchAdministrators(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil { // && !authz.GetFeatures(ctx).PermissionCheckV2 {
		administratorsCheckPermission(ctx, admins, permissionCheck)
	}
	return admins, nil
}

func (q *Queries) searchAdministrators(ctx context.Context, queries *MembershipSearchQuery, permissionCheckV2 bool) (administrators *Administrators, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, queryArgs, scan := prepareAdministratorsQuery(ctx, queries, permissionCheckV2)
	eq := sq.Eq{membershipInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-TODO", "Errors.Query.InvalidRequest")
	}
	latestState, err := q.latestState(ctx, orgMemberTable, instanceMemberTable, projectMemberTable, projectGrantMemberTable)
	if err != nil {
		return nil, err
	}
	queryArgs = append(queryArgs, args...)

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		administrators, err = scan(rows)
		return err
	}, stmt, queryArgs...)
	if err != nil {
		return nil, err
	}
	administrators.State = latestState
	return administrators, nil
}

func prepareAdministratorsQuery(ctx context.Context, queries *MembershipSearchQuery, permissionV2 bool) (sq.SelectBuilder, []interface{}, func(*sql.Rows) (*Administrators, error)) {
	query, args := getMembershipFromQuery(ctx, queries, permissionV2)
	return sq.Select(
			MembershipUserID.identifier(),
			membershipRoles.identifier(),
			MembershipCreationDate.identifier(),
			MembershipChangeDate.identifier(),
			membershipResourceOwner.identifier(),
			membershipOrgID.identifier(),
			membershipIAMID.identifier(),
			membershipProjectID.identifier(),
			membershipGrantID.identifier(),
			ProjectGrantColumnGrantedOrgID.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnName.identifier(),
			OrgColumnName.identifier(),
			InstanceColumnName.identifier(),
			LoginNameNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			UserTypeCol.identifier(),
			UserResourceOwnerCol.identifier(),
			countColumn.identifier(),
		).From(query).
			LeftJoin(join(ProjectColumnID, membershipProjectID)).
			LeftJoin(join(ProjectGrantColumnGrantID, membershipGrantID)).
			LeftJoin(join(OrgColumnID, membershipOrgID)).
			LeftJoin(join(InstanceColumnID, membershipInstanceID)).
			LeftJoin(join(HumanUserIDCol, OrgMemberUserID)).
			LeftJoin(join(MachineUserIDCol, OrgMemberUserID)).
			LeftJoin(join(UserIDCol, OrgMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, OrgMemberUserID)).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		args,
		func(rows *sql.Rows) (*Administrators, error) {
			administrators := make([]*Administrator, 0)
			var count uint64
			for rows.Next() {

				var (
					administrator        = new(Administrator)
					userID               = sql.NullString{}
					orgID                = sql.NullString{}
					instanceID           = sql.NullString{}
					projectID            = sql.NullString{}
					grantID              = sql.NullString{}
					grantedOrgID         = sql.NullString{}
					projectName          = sql.NullString{}
					orgName              = sql.NullString{}
					instanceName         = sql.NullString{}
					projectResourceOwner = sql.NullString{}
					loginName            = sql.NullString{}
					displayName          = sql.NullString{}
					machineName          = sql.NullString{}
					avatarURL            = sql.NullString{}
					userType             = sql.NullInt32{}
					userResourceOwner    = sql.NullString{}
				)

				err := rows.Scan(
					&userID,
					&administrator.Roles,
					&administrator.CreationDate,
					&administrator.ChangeDate,
					&administrator.ResourceOwner,
					&orgID,
					&instanceID,
					&projectID,
					&grantID,
					&grantedOrgID,
					&projectResourceOwner,
					&projectName,
					&orgName,
					&instanceName,
					&loginName,
					&displayName,
					&machineName,
					&avatarURL,
					&userType,
					&userResourceOwner,
					&count,
				)

				if err != nil {
					return nil, err
				}

				if userID.Valid {
					administrator.User = &UserAdministrator{
						UserID:        userID.String,
						LoginName:     loginName.String,
						DisplayName:   displayName.String,
						ResourceOwner: userResourceOwner.String,
					}
				}

				if orgID.Valid {
					administrator.Org = &OrgAdministrator{
						OrgID: orgID.String,
						Name:  orgName.String,
					}
				}
				if instanceID.Valid {
					administrator.Instance = &InstanceAdministrator{
						InstanceID: instanceID.String,
						Name:       instanceName.String,
					}
				}
				if projectID.Valid && grantID.Valid && grantedOrgID.Valid {
					administrator.ProjectGrant = &ProjectGrantAdministrator{
						ProjectID:     projectID.String,
						ProjectName:   projectName.String,
						GrantID:       grantID.String,
						GrantedOrgID:  grantedOrgID.String,
						ResourceOwner: projectResourceOwner.String,
					}
				} else if projectID.Valid {
					administrator.Project = &ProjectAdministrator{
						ProjectID:     projectID.String,
						Name:          projectName.String,
						ResourceOwner: projectResourceOwner.String,
					}
				}

				administrators = append(administrators, administrator)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-TODO", "Errors.Query.CloseRows")
			}

			return &Administrators{
				Administrators: administrators,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
