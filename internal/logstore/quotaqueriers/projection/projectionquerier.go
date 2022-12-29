package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/repository/instance"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.QuotaQuerier = (*projectionQuerier)(nil)

type projectionQuerier struct {
	dbClient *sql.DB
}

func NewQuerier(dbClient *sql.DB) *projectionQuerier {
	return &projectionQuerier{dbClient: dbClient}
}

func (p *projectionQuerier) GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (*query.Quota, error) {
	return query.GetInstanceQuota(ctx, p.dbClient, instanceID, unit)
}

func (p *projectionQuerier) GetDueQuotaNotifications(ctx context.Context, q *query.Quota, used uint64) ([]*instance.QuotaNotifiedEvent, error) {
	return query.GetDueInstanceQuotaNotifications(ctx, p.dbClient, q, used)
}
