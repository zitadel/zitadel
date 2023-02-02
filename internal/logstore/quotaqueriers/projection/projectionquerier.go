package projection

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
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
	quota, err := query.GetQuota(ctx, p.dbClient, instanceID, unit)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return quota, err
}

func (p *projectionQuerier) GetDueQuotaNotifications(ctx context.Context, q *query.Quota, used uint64) ([]*quota.NotifiedEvent, error) {
	return query.GetDueQuotaNotifications(ctx, p.dbClient, q, used)
}
