package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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

func TriggerOIDCUserInfoProjections(ctx context.Context) {
	triggerBatch(ctx, oidcUserInfoTriggerHandlers()...)
}

//go:embed embed/userinfo_by_id.sql
var oidcUserInfoQuery string

func (q *Queries) GetOIDCUserInfo(ctx context.Context, userID string, roleAudience []string) (_ *OIDCUserInfo, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var data []byte
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&data)
	},
		oidcUserInfoQuery,
		userID, authz.GetInstance(ctx).InstanceID(), database.TextArray[string](roleAudience),
	)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Oath6", "Errors.Internal")
	}

	userInfo := new(OIDCUserInfo)
	if err = json.Unmarshal(data, userInfo); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Vohs6", "Errors.Internal")
	}
	if userInfo.User == nil {
		return nil, errors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.User.NotFound")
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
