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
	quota     *query.CurrentQuotaPeriod
	lastUsage uint64
}

func NewNoopQuerier(quota *query.CurrentQuotaPeriod) *inmemReporter {
	return &inmemReporter{quota: quota}
}

func (i *inmemReporter) GetQuota(context.Context, string, quota.Unit) (*query.CurrentQuotaPeriod, error) {
	return i.quota, nil
}

func (*inmemReporter) GetDueQuotaNotifications(context.Context, *query.CurrentQuotaPeriod, uint64) ([]*instance.QuotaNotifiedEvent, error) {
	return nil, nil
}
