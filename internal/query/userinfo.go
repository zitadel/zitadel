package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
)

//go:embed embed/userinfo_by_id.sql
var oidcUserInfoQuery string

func (q *Queries) GetOIDCUserInfo(ctx context.Context, userID string) (_ *OIDCUserInfo, err error) {
	userInfo := new(OIDCUserInfo)
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		var data []byte
		if err := row.Scan(&data); err != nil {
			return err
		}
		return json.Unmarshal(data, userInfo)
	}, oidcUserInfoQuery, userID, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, fmt.Errorf("get oidc user info: %w", err)
	}
	if userInfo.User == nil {
		return nil, errors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.User.NotFound")
	}

	return userInfo, nil
}

type OIDCUserInfo struct {
	User     *User          `json:"user,omitempty"`
	Metadata []UserMetadata `json:"metadata,omitempty"`
	Org      *struct {
		ID            string `json:"id,omitempty"`
		Name          string `json:"name,omitempty"`
		PrimaryDomain string `json:"primary_domain,omitempty"`
	} `json:"org,omitempty"`
}
