package view

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
)

const (
	applicationTable = "auth.applications"
)

func (v *View) ApplicationByClientID(ctx context.Context, clientID string) (*model.ApplicationView, error) {
	req := &proj_model.ApplicationSearchRequest{
		Limit: 1,
		Queries: []*proj_model.ApplicationSearchQuery{
			{
				Key:    proj_model.APPLICATIONSEARCHKEY_OIDC_CLIENT_ID,
				Method: global_model.SEARCHMETHOD_EQUALS,
				Value:  clientID,
			},
		},
	}
	apps, count, err := view.SearchApplications(v.Db, applicationTable, req)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "VIEW-sd6JQ", "cannot find client")
	}
	if count != 1 {
		return nil, errors.ThrowPreconditionFailed(nil, "VIEW-dfw3as", "cannot find client")
	}
	return apps[0], nil
}
