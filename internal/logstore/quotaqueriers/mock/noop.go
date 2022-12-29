package mock

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/instance"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.QuotaQuerier = (*inmemReporter)(nil)

type inmemReporter struct {
	quota     *query.Quota
	lastUsage uint64
}

func NewNoopQuerier(quota *query.Quota) *inmemReporter {
	return &inmemReporter{quota: quota}
}

func (i *inmemReporter) GetQuota(context.Context, string, quota.Unit) (*query.Quota, error) {
	return i.quota, nil
}

func (*inmemReporter) GetDueQuotaNotifications(context.Context, *query.Quota, uint64) ([]*instance.QuotaNotifiedEvent, error) {
	return nil, nil
}
