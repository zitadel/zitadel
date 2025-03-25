package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"strings"
	"sync"
	"text/template"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// oidcUserInfoTriggerHandlers slice can only be created after zitadel
// is fully initialized, otherwise the handlers are nil.
// OnceValue takes care of creating the slice on the first request
// and than will always return the same slice on subsequent requests.
var oidcUserInfoTriggerHandlers = sync.OnceValue(func() []*handler.Handler {
	return []*handler.Handler{
		projection.UserProjection,
		projection.UserMetadataProjection,
		projection.UserGrantProjection,
		projection.OrgProjection,
		projection.ProjectProjection,
	}
})

// TriggerOIDCUserInfoProjections triggers all projections
// relevant to userinfo queries concurrently.
func TriggerOIDCUserInfoProjections(ctx context.Context) {
	triggerBatch(ctx, oidcUserInfoTriggerHandlers()...)
}

var (
	//go:embed userinfo_by_id.sql
	oidcUserInfoQueryTmpl           string
	oidcUserInfoQuery               string
	oidcUserInfoWithRoleOrgIDsQuery string

	//go:embed userinfo_by_group_id.sql
	oidcGroupInfoQueryTmpl           string
	oidcGroupInfoQuery               string
	oidcGroupInfoWithRoleOrgIDsQuery string
)

// build the two variants of the userInfo query
// and the two variants of the groupInfo query
func init() {
	tmpl := template.Must(template.New("oidcUserInfoQuery").Parse(oidcUserInfoQueryTmpl))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, false); err != nil {
		panic(err)
	}
	oidcUserInfoQuery = buf.String()
	buf.Reset()

	if err := tmpl.Execute(&buf, true); err != nil {
		panic(err)
	}
	oidcUserInfoWithRoleOrgIDsQuery = buf.String()
	buf.Reset()

	grp_tmpl := template.Must(template.New("oidcGroupInfoQuery").Parse(oidcGroupInfoQueryTmpl))
	if err := grp_tmpl.Execute(&buf, false); err != nil {
		panic(err)
	}
	oidcGroupInfoQuery = buf.String()
	buf.Reset()

	if err := grp_tmpl.Execute(&buf, true); err != nil {
		panic(err)
	}
	oidcGroupInfoWithRoleOrgIDsQuery = buf.String()
	buf.Reset()
}

func (q *Queries) GetOIDCUserInfo(ctx context.Context, userID string, roleAudience []string, roleOrgIDs ...string) (userInfo *OIDCUserInfo, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if len(roleOrgIDs) > 0 {
		userInfo, err = database.QueryJSONObject[OIDCUserInfo](ctx, q.client, oidcUserInfoWithRoleOrgIDsQuery,
			userID, authz.GetInstance(ctx).InstanceID(), database.TextArray[string](roleAudience), database.TextArray[string](roleOrgIDs),
		)
	} else {
		userInfo, err = database.QueryJSONObject[OIDCUserInfo](ctx, q.client, oidcUserInfoQuery,
			userID, authz.GetInstance(ctx).InstanceID(), database.TextArray[string](roleAudience),
		)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Eey2a", "Errors.User.NotFound")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Oath6", "Errors.Internal")
	}
	if userInfo.User == nil {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.User.NotFound")
	}

	return userInfo, nil
}

type OIDCUserInfo struct {
	User       *User          `json:"user,omitempty"`
	Metadata   []UserMetadata `json:"metadata,omitempty"`
	Org        *UserInfoOrg   `json:"org,omitempty"`
	UserGrants []UserGrant    `json:"user_grants,omitempty"`
	// GroupGrants []GroupGrant   `json:"group_grants,omitempty"`
}

type UserInfoOrg struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	PrimaryDomain string `json:"primary_domain,omitempty"`
}

//go:embed userinfo_client_by_id.sql
var oidcUserinfoClientQuery string

func (q *Queries) GetOIDCUserinfoClientByID(ctx context.Context, clientID string) (projectID string, projectRoleAssertion bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	scan := func(row *sql.Row) error {
		err := row.Scan(&projectID, &projectRoleAssertion)
		return err
	}

	err = q.client.QueryRowContext(ctx, scan, oidcUserinfoClientQuery, authz.GetInstance(ctx).InstanceID(), clientID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, zerrors.ThrowNotFound(err, "QUERY-beeW8", "Errors.App.NotFound")
	}
	if err != nil {
		return "", false, zerrors.ThrowInternal(err, "QUERY-Ais4r", "Errors.Internal")
	}
	return projectID, projectRoleAssertion, nil
}

type OIDCGroupInfo struct {
	Group       *Group       `json:"group,omitempty"`
	GroupGrants []GroupGrant `json:"group_grants,omitempty"`
}

func (q *Queries) getOIDCGroupInfo(ctx context.Context, groupID string, roleAudience []string, roleOrgIDs ...string) (groupInfo *OIDCGroupInfo, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if len(roleOrgIDs) > 0 {
		groupInfo, err = database.QueryJSONObject[OIDCGroupInfo](ctx, q.client, oidcGroupInfoWithRoleOrgIDsQuery,
			groupID, authz.GetInstance(ctx).InstanceID(), database.TextArray[string](roleAudience), database.TextArray[string](roleOrgIDs),
		)
	} else {
		groupInfo, err = database.QueryJSONObject[OIDCGroupInfo](ctx, q.client, oidcGroupInfoQuery,
			groupID, authz.GetInstance(ctx).InstanceID(), database.TextArray[string](roleAudience),
		)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Eey2a", "Errors.Group.NotFound")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Oath6", "Errors.Internal")
	}
	if groupInfo.Group == nil {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.Group.NotFound")
	}

	return groupInfo, nil
}

type OIDCGroupInfos struct {
	Group []OIDCGroupInfo `json:"group,omitempty"`
}

func (q *Queries) GetOIDCGroupInfos(ctx context.Context, groupIDs []string, roleAudience []string, roleOrgIDs ...string) (groupInfos *OIDCGroupInfos, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupInfos = &OIDCGroupInfos{
		Group: make([]OIDCGroupInfo, 0, len(groupIDs)),
	}
	for _, groupID := range groupIDs {
		groupInfo, err := q.getOIDCGroupInfo(ctx, groupID, roleAudience, roleOrgIDs...)
		if err != nil {
			return nil, err
		}
		groupInfos.Group = append(groupInfos.Group, *groupInfo)
	}

	return groupInfos, nil
}

func (q *Queries) GetOIDCGroupInfosV2(ctx context.Context, groups *Groups, roleAudience []string, roleOrgIDs ...string) (groupInfos *OIDCGroupInfos, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupInfos = &OIDCGroupInfos{
		Group: make([]OIDCGroupInfo, 0, len(groups.Groups)),
	}
	for _, group := range groups.Groups {
		groupInfo, err := q.getOIDCGroupInfo(ctx, group.ID, roleAudience, roleOrgIDs...)
		if err != nil {
			return nil, err
		}
		groupInfos.Group = append(groupInfos.Group, *groupInfo)
	}

	return groupInfos, nil
}
