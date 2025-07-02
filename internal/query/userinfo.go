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

var (
	//go:embed userinfo_by_id.sql
	oidcUserInfoQueryTmpl           string
	oidcUserInfoQuery               string
	oidcUserInfoWithRoleOrgIDsQuery string
)

// build the two variants of the userInfo query
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
